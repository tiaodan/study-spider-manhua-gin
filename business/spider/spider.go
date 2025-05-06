package spider

import (
	"fmt"
	"regexp"
	"strconv"
	"study-spider-manhua-gin/db"
	"study-spider-manhua-gin/errorutil"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"
	"study-spider-manhua-gin/util/langutil"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func Spider(context *gin.Context) {
	log.Debug("爬虫开始-------------------------------")
	c := colly.NewCollector()

	// 设置请求限制（每秒最多2个请求）
	// spider.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*",
	// 	Parallelism: 2,
	// 	RandomDelay: 5 * time.Second,
	// })

	// 修改为解析所有元素并检查文本内容
	// spider.OnHTML("*", func(e *colly.HTMLElement) { // 修改选择器为"*"匹配所有元素
	// 	text := e.Text
	// 	if strings.Contains(text, "还没有看过的漫画") { // 新增文本判断逻辑
	// 		log.Debug("找到目标文本: %s", text)
	// 	}
	// 	link := e.Attr("href")
	// 	log.Debug("发现链接：%s", link)
	// })

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
		// 2. 转简体
		// 3. 数据清洗
		// 4. 把参数赋值给 comic对象,把每个对象存起来
		// 4.1 统一打印
		// 5. 去重
		// 6. 插入数据库
		// 7. 重置变量

		comicArr := []*models.Comic{}
		moveRepeatComics := make(map[string]string) // 用map做去重,保存漫画名称
		e.ForEach(".common-comic-item", func(i int, element *colly.HTMLElement) {
			// 创建对象comic
			comic := &models.Comic{}

			// 1. 爬数据, 自动去重前后空格
			// 1.1 爬名字,唯一索引,如果为空, return
			comicNameTradition := element.ChildText(".comic__title")
			if comicNameTradition == "" {
				log.Debug("漫画名称为空, 跳过")
				return
			}
			// 1.1.1 通过名字去重
			if _, exists := moveRepeatComics[comicNameTradition]; exists {
				log.Info("存在重复项: ", comicNameTradition)
				return
			}
			// 1.1.2 把不重复的加入到map里
			moveRepeatComics[comicNameTradition] = comicNameTradition

			// 1.2 爬更新到 ?集
			updateStrTrad := element.ChildText(".comic-update a")

			// 1.3 爬人气
			hitsStrTrad := element.ChildText(".comic-count")

			// 1.4 爬封面链接
			coverUrl := element.ChildAttr(".cover img", "data-original")

			// 1.5 爬漫画链接
			comicUrl := element.ChildAttr(".cover", "href")

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

			// 3. 数据清洗
			// 封面书籍链接 如果包含http或者https, needTcp = 0
			// isNeedTcp := 0   // 是否需要http 或https ,说明链接不带http头，默认不需要。不用赋值，因为new comic就有默认值0了
			// isNeedHttps := 0 // 默认,不需要https。不用赋值，因为new comic就有默认值0了
			if !(langutil.IsHTTPOrHTTPS(comicUrl)) {
				comic.NeedTcp = 1
			}

			// 判断封面链接
			if !(langutil.IsHTTPOrHTTPS(coverUrl)) {
				comic.CoverNeedTcp = 1
			}

			// 判断是否完结, 通过"更新至" 是否== "休刊公告"
			if updateStr == "休刊公告" || updateStr == "后记" {
				comic.End = 1 // 不用设置默认值0, 因为new comic 时会有默认值0
			}

			// 清洗 “人气”,提取字符串中数字
			re := regexp.MustCompile(`(\d+\.?\d*)\s*([^\d\s]+)`) // 定义正则表达式，匹配数字和单位
			matches := re.FindStringSubmatch(HitsStr)
			log.Info("--------------- matches = ", matches)
			if len(matches) >= 3 {
				hitsNumStr := matches[1] // 匹配全部字符串 如 95.2 万
				hitsUnit := matches[2]   // 人气数字 如：95.2
				numUnit := 1             // 单位 如：万
				if hitsUnit == "亿" {
					numUnit = 100000000
				} else if hitsUnit == "万" {
					numUnit = 10000
				} else if hitsUnit == "千" {
					numUnit = 1000
				}

				// 计算具体数字 HitsNum * hitsUnit
				hitsFloat, err := strconv.ParseFloat(hitsNumStr, 64)
				if err != nil || hitsFloat < 0 {
					comic.Hits = 0 // 错误或负值设为0
				} else {
					comic.Hits = uint(hitsFloat * float64(numUnit))
				}
			}

			// 4. 把参数赋值给 comic对象
			comic.Name = comicName
			comic.Update = updateStr
			comic.ComicUrl = comicUrl
			comic.CoverUrl = coverUrl

			// comic对象加入到数组中,把每个对象存起来
			comicArr = append(comicArr, comic)

			// 4.1 统一打印
			log.Debug("更新到: ", updateStr)
			log.Debug("人气: ", HitsStr)
			log.Debug("计算后人气: ", comic.Hits)
			log.Debug("封面链接: ", coverUrl)
			log.Debug("漫画链接:  ", comicUrl)
			log.Debugf("当前%d, 漫画名称转简体= %s -> %s", i+1, comicNameTradition, comicName)
			log.Infof("序号= %d, comic对象: id name 更新至 点击量 封面 书籍url 是否完结 needTcp  coverNeedTcp : %v", i+1, comic)
		})

		// 5. 插入数据库
		db.ComicBatchAdd(comicArr)
		// 7. 重置变量
		comicArr = comicArr[:0]
		moveRepeatComics = make(map[string]string)
	})

	// 修改本地文件访问路径为file协议格式
	err := c.Visit("http://localhost:8080/test/index.html") // 使用file协议访问本地文件

	errorutil.ErrorPrint(err, "访问失败")
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
