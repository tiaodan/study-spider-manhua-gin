/*
*
参考代码，弃用 !!
*/
package spider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util/langutil"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

// 思路：
// 1. 提取请求数据,拼接完整请求url
// 1.1 判断参数是否符合要求，不符合返回错误 return
// 2. 爬取
// 2.1 new 爬取对象
// 2.2 建一个任务队列
// 2.3 设置并发数，和爬取限制
// 2.4 注册 HTML 解析逻辑
// 2.5 可选：注册请求前的回调
// 2.6 可选：注册错误处理
// 2.7 添加多个页面到队列中
// 2.8 启动爬虫并等待所有任务完成
func Spider(context *gin.Context) {
	log.Debug("爬虫开始-------------------------------")
	// 读取 Body
	log.Info("传参 body = ", context.Request.Body)

	// 1. 提取请求数据
	var requestBody models.SpiderRequestBody
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		log.Error("解析请求体失败, err: ", err)
		context.JSON(400, gin.H{"error": err.Error()})
		return // 必须保留 return，确保绑定失败时提前退出
	}
	var fullUrl string // 完整请求url
	if requestBody.NeedTcp == 1 {
		if requestBody.NeedHttps == 1 { // 如果需要https
			fullUrl += "https://"
		}
		fullUrl += "http://"
	}
	fullUrl += requestBody.WebsitePrefix + requestBody.Url // 完整请求url,差结尾数字

	log.Debug("请求数据: url: ", requestBody.Url)
	log.Debug("请求数据: websitePrefix: ", requestBody.WebsitePrefix)
	log.Debug("请求数据: needTcp: ", requestBody.NeedTcp)
	log.Debug("请求数据: needHttps: ", requestBody.NeedHttps)
	log.Debug("请求数据: endNum: ", requestBody.EndNum)
	log.Debug("请求数据: fullUrl: ", fullUrl)

	// 1.1 判断参数是否符合要求，不符合返回错误 return
	if requestBody.EndNum <= 0 {
		context.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	c := colly.NewCollector()

	// 创建一个并发为3的任务队列，使用内存存储

	// 设置请求限制（每秒最多2个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		RandomDelay: 5 * time.Second,
	})
	// 使用队列控制任务调度（最多并发3个Url）
	q, _ := queue.New(3, &queue.InMemoryQueueStorage{MaxSize: 100})

	// 线程安全的 去重map, 用于爬某类所有page数据
	var comicNamePool sync.Map

	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的
	c.OnHTML(".cate-comic-list", func(e *colly.HTMLElement) {
		// log.Debug("html 元素名称Name= ", e.Name) // 一般显示div a span 等标签名
		// log.Debug("html 元素DOM= ", e.DOM) // 就显示地址 &{[0xc0002fa540] 0xc00060e4e0 <nil>}
		// log.Debug("html 元素内容Text= ", e.Text)

		// 在选一个可以foreach匹配的子标签
		log.Debug("匹配到.cate-comic-list")
		// 思路：
		// 1. 爬数据, 自动去重前后空格
		// 1.1 爬名字
		// 1.2 爬更新到 ?集
		// 1.3 爬人气
		// 1.4 爬封面链接
		// 1.5 爬漫画链接
		// 1.6 爬简介-short
		// 1.7 标志位不用管，插入默认0值 (spider_end、download_end、upload_aws_end、upload_baidu_end)
		// 1.8 不用管 (评分)，插入默认0值，因为这个页面爬取不到
		// 2. 转简体
		// 3. 数据清洗
		// 4. 把参数赋值给 comic对象,把每个对象存起来
		// 4.1 统一打印
		// 5. 去重
		// 6. 插入数据库
		// 7. 重置变量

		comicArr := []*models.ComicSpider{}
		moveRepeatComics := make(map[string]string) // 用map做去重,保存漫画名称
		e.ForEach(".common-comic-item", func(i int, element *colly.HTMLElement) {
			// 创建对象comic
			comic := &models.ComicSpider{}

			// 1. 爬数据, 自动去重前后空格
			// 1.1 爬名字,唯一索引,如果为空, return
			comicNameTradition := strings.TrimSpace(element.ChildText(".comic__title"))
			if comicNameTradition == "" {
				log.Debug("漫画名称为空, 跳过")
				return
			}
			// 1.1.1 通过名字去重
			if _, exists := moveRepeatComics[comicNameTradition]; exists {
				log.Info("存在重复项: ", comicNameTradition)
				return
			}
			if _, exists := comicNamePool.Load(comicNameTradition); exists { // 线程安全方式提取
				log.Info("--------------------大comic池, 存在重复项: ", comicNameTradition)
				return
			}

			// 1.1.2 把不重复的加入到map里
			moveRepeatComics[comicNameTradition] = comicNameTradition
			comicNamePool.Store(comicNameTradition, comicNameTradition) // 把不重复的加到大comic池里

			// 1.2 爬更新到 ?集
			updateStrTrad := strings.TrimSpace(element.ChildText(".comic-update a"))

			// 1.3 爬人气
			hitsStrTrad := strings.TrimSpace(element.ChildText(".comic-count"))

			// 1.4 爬封面链接
			coverUrlApiPath := strings.TrimSpace(element.ChildAttr(".cover img", "data-original"))

			// 1.5 爬漫画链接
			comicUrlApiPath := strings.TrimSpace(element.ChildAttr(".cover", "href"))

			// 1.6 爬简介-short - 繁体
			comicBriefShortTrad := strings.TrimSpace(element.ChildText(".cover p"))

			// 2 转简体
			// 2.1 漫画名称, 转换为简体中文
			comicName, err := langutil.TraditionalToSimplified(comicNameTradition)
			if err != nil {
				log.Errorf("转换为简体中文失败: %v", err)
				comicName = comicNameTradition // 如果转换失败，使用原名称
			}

			// 2.2 转换更新到 ?集
			updateStr, err := langutil.TraditionalToSimplified(updateStrTrad)
			if err != nil {
				log.Errorf("转换为简体中文失败: %v", err)
				updateStr = updateStrTrad // 如果转换失败，使用原名称
			}

			// 2.3 人气
			HitsStr, err := langutil.TraditionalToSimplified(hitsStrTrad)
			if err != nil {
				log.Errorf("转换为简体中文失败: %v", err)
				HitsStr = hitsStrTrad // 如果转换失败，使用原名称
			}

			// 2.4 爬简介-short
			comicBriefShort, err := langutil.TraditionalToSimplified(comicBriefShortTrad)
			if err != nil {
				log.Errorf("转换为简体中文失败: %v", err)
				comicBriefShort = comicBriefShortTrad // 如果转换失败，使用原名称
			}

			// 3. 数据清洗

			// 判断是否完结, 传参如果带标志位，就不判断了。通过字段包含 "更新至" 是否== "休刊公告"
			switch requestBody.End {
			case 1:
				comic.End = true // 不用设置默认值0, 因为new comic 时会有默认值0
			case 2: // 传参2, 程序自行判断
				if strings.Contains(updateStr, "休刊公告") || strings.Contains(updateStr, "后记") {
					comic.End = true // 不用设置默认值0, 因为new comic 时会有默认值0
				}
			}

			// 清洗 “人气”,提取字符串中数字
			re := regexp.MustCompile(`(\d+\.?\d*)\s*([^\d\s]+)`) // 定义正则表达式，匹配数字和单位
			matches := re.FindStringSubmatch(HitsStr)
			log.Info("--------------- matches = ", matches)
			var hitsNumStr string // xx数字
			var hitsUnit string   // 单位 如；万、千
			numUnit := 1          // 单位 如：万, 默认1个
			if len(matches) >= 3 {
				hitsNumStr = matches[1] // 匹配全部字符串 如 95.2 万
				hitsUnit = matches[2]   // 人气数字 如：95.2
				switch hitsUnit {
				case "亿":
					numUnit = 100000000
				case "万":
					numUnit = 10000
				case "千":
					numUnit = 1000
				}
			} else { // 重新正则匹配
				re = regexp.MustCompile(`(\d+\.?\d*)\s*`) // 定义正则表达式，匹配数字和单位
				newMatches := re.FindStringSubmatch(HitsStr)
				log.Info("--------------- newMatches !=3 ", newMatches)
				log.Info("--------------- newMatches[1] ", newMatches[1])
				hitsNumStr = newMatches[1] // 匹配全部字符串 如 95.2 万
			}

			// 计算具体数字 HitsNum * hitsUnit
			hitsFloat, err := strconv.ParseFloat(hitsNumStr, 64)
			if err != nil || hitsFloat < 0 {
				comic.Hits = 0 // 错误或负值设为0
			} else {
				comic.Hits = int(hitsFloat * float64(numUnit))
			}

			// 4. 把参数赋值给 comic对象
			comic.Name = comicName
			comic.LatestChapter = updateStr
			comic.ComicUrlApiPath = comicUrlApiPath
			comic.CoverUrlApiPath = coverUrlApiPath
			comic.BriefShort = comicBriefShort
			comic.CountryId = requestBody.CountryId // 外键
			comic.WebsiteId = requestBody.WebsiteId
			comic.PornTypeId = requestBody.PornTypeId
			comic.TypeId = requestBody.TypeId

			// comic对象加入到数组中,把每个对象存起来
			comicArr = append(comicArr, comic)

			// 4.1 统一打印
			log.Debug("更新到: ", updateStr)
			log.Debug("人气: ", HitsStr)
			log.Debug("计算后人气: ", comic.Hits)
			log.Debug("封面链接: ", coverUrlApiPath)
			log.Debug("漫画链接:  ", comicUrlApiPath)
			log.Debugf("当前%d, 漫画名称转简体= %s -> %s", i+1, comicNameTradition, comicName)
			log.Infof("序号= %d, comic对象: id name 更新至 点击量 封面 书籍url 是否完结 needTcp  coverNeedTcp : %v", i+1, comic)
		})

		// 5. 插入数据库
		uniqueIndexArr := []string{"Name", "CountryId", "WebsiteId", "pornTypeId", "TypeId"}
		updateColArr := []string{"update", "hits", "comic_url_api_path", "cover_url_api_path", "brief_short", "brief_long", "end",
			"star", "need_tcp", "cover_need_tcp", "spider_end", "download_end", "upload_aws_end", "upload_baidu_end",
			"updated_at"} // 要传updated_at ，upsert必须传, UPDATE()方法不用传，会自动改
		db.DBUpsertBatch(db.DBComic, comicArr, uniqueIndexArr, updateColArr)
		// 7. 重置变量
		comicArr = comicArr[:0]
		moveRepeatComics = make(map[string]string)
		comicNamePool = sync.Map{}
	})

	// 修改本地文件访问路径为file协议格式 - 注释了
	// err := c.Visit("http://localhost:8080/test/index.html") // 使用file协议访问本地文件，用了队列后，就不用c.Visit了

	// 添加任务到队列
	// for i := 1; i <= requestBody.EndNum; i++ {
	// 	q.AddURL(fullUrl + strconv.Itoa(i))
	// }

	// 测试用
	q.AddURL("http://localhost:8080/test/index.html")

	// 启动对垒
	q.Run(c)
	// 返回信息
	context.JSON(200, "添加成功")
	log.Debug("爬虫结束-------------------------------")
}

// 爬百度主页, 熟悉练手
func SpiderBaiduTest(context *gin.Context) {
	c := colly.NewCollector()

	// 获取 "百度一下" 按钮文本
	c.OnHTML("#su", func(e *colly.HTMLElement) {
		baiduBtn := e.Attr("value")
		fmt.Println(baiduBtn)
	})

	// 开始访问
	err := c.Visit("http://www.baidu.com/") // 百度 测试
	errorutil.ErrorPrint(err, "访问失败")
}

// 根据第一页链接，自动提取出 请求链接、尾页号码、是否需要http请求。待后续封装
// 前期可以先人为提供信息
func GetFirstPageLink(url string) (string, int, int) {
	// 提取请求参数

	return "", 0, 0 // 暂时这么些，待修改
}
