/**
功能：封装 通用的爬虫 模板 (所有跟书类型相关的)
	什么是书类型？
	- 有书名、章节、章节里具体的内容(图片、视频、音频、文字等)
适用此模板有哪些？
	- 漫画网站
	- 有声书网站
	- 小说网站
	- 视频网站
*/

package spider

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util/langutil"
	"study-spider-manhua-gin/src/util/stringutil"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/tidwall/gjson"
)

// -- 初始化 ------------------------------------------------------------------------------

// -- 初始化 ------------------------------------------- end -----------------------------------

// -- 方法 ------------------------------------------------------------------------------

// 爬取类别。如“漫画”“小说”）。跟type相比，是更上层
func BookTemSpiderCategory(context *gin.Context) {

}

// 爬取分类，如“热血”“恋爱”“悬疑” -》 用type
/*
作用简单说：
	- 爬取网站上某一种分类。如有声书：悬疑、有声书：科幻、有声书：历史等

作用详细说:

核心思路:
	1. 准备爬取需要的参数
	2. 爬取
	3. 插入db

参考通用思路：
	1. 校验传参
		- 前端参数转成对象
		- 是否需要简单清洗？
		- 校验
		- 分析前端参数含义
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑 - 爬取 - 插入db
		-- 拼接第一页 完整url
		-- new 爬虫对象
		-- 建一个爬虫对象
		-- 设置并发数，和爬取限制
		-- 注册 HTML 解析逻辑
		-- 添加多个爬虫 到到队列中
	5. 返回结果

参数：
	1. context *gin.Context 类型 // 前端传参
		至少要包含的参数：
		1. websiteId int 外键-对应 website表 name_id // 因为：漫画、有声书、小说、视频都涉及
		2. pornTypeId int 外键-对应 porn_type表 name_id // 因为：漫画、有声书、小说、视频都涉及.都要区分是否是色情内容
		3. countryId int 外键-对应 country表 name_id // 因为：漫画、有声书、小说、视频都涉及
		4. typeId int 外键-对应 type表 name_id // 因为：漫画、有声书、小说、视频都涉及。都涉及类型：如 爱情、悬疑、冒险等
		例如：
		{
			"countryId": 1,                        // 国家id
			"pornTypeId": 1,                       // 色情类型id
			"websiteId": 3,                        // 网站id
			"targetSiteTypeId": 1,                 // 目标网站类型id,即：这个网站自己是怎么分类型的，如：悬疑、冒险等
			"typeId": 2                            // 类型id。 我的网站自己是怎么分类型的，如：悬疑、冒险等
			//
			"websitePrefix": "www.manhuagui.com",  // 网站前缀，现在想的是最后不带/
			"url": "list/c1-p",                    // 排除前缀后的，url路径，需要带/
			"needTcp": 1,   					   // 完整请求，是否需要带 http / https 因为有的爬取的 book的链接，有的带http，有的不带
			"needHttps": 1, 					   // 完整请求，是否需要带  https
			"endNum": 5, 						   // 尾页号码
			"endJudgeMethod": 2,                   // 完结判断方式 0：全部写成 未完结false 1：全部写成完结true 2：程序自动判断
		}

返回：
	[]T 泛型数组 // 各种表的对象 数组。如comic表对象，website表对象
	- 无
	- 不通过return返回数据，通过 context.JSON(响应码, "信息") 返回
	直接返回给前端json数据

注意：
*/
func BookTemSpiderTypeByHtml(context *gin.Context) {

	// 1. 校验传参
	// -- 打印参数
	// 参数头
	argsHeader := []string{"websiteId", "pornTypeId", "countryId", "typeId", "websitePrefix", "url", "needTcp", "needHttps", "endNum", "endJudgeMethod"}
	// log.Infof("book类, 通用爬取模板, 传参 body = %v", context.Request.Body) // 打印的顺序，和前端传参不一致，因此看这个日志也没有意义

	// -- 前端参数转成对象
	var requestBody models.SpiderRequestBody2
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		log.Error("解析请求体失败, 建议:1) 先检查传参类型 2) 检查传参数值。err= ", err)
		context.JSON(400, gin.H{"解析请求体失败。建议:1) 先检查传参类型 2) 检查传参数值。 error= ": err.Error()})
		return // 必须保留 return，确保绑定失败时提前退出
	}
	// 2. 数据清洗
	stringutil.TrimSpaceObj(&requestBody) // 去除空格

	log.Infof("book类, 通用爬取模板, 前端传参-> 对象, body-key=%v, body = %v", argsHeader, &requestBody) // 打印的顺序，和前端传参不一致，因此看这个日志也没有意义
	log.Info("book类, 通用爬取模板, 爬取整个分类. 原网站分类: ", requestBody.TargetSiteTypeId)

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 判断参数是否符合要求，不符合返回错误 return
	if requestBody.EndNum <= 0 {
		context.JSON(400, gin.H{"error": "参数错误, endNum 应>0"})
		return
	}

	// 3. 执行核心逻辑 - 爬取 - 插入db
	// -- 拼接第一页 完整url
	firstPageLink := GetOneTypeFirstPageFullURL(requestBody.NeedTcp, requestBody.NeedHttps, requestBody.WebsitePrefix,
		requestBody.ApiPath, requestBody.ParamArr)

	log.Infof("分类: %v, 第一页Link= %v", requestBody.TargetSiteTypeId, firstPageLink)
	// -- 爬取

	// -- 建一个爬虫对象
	c := colly.NewCollector()

	// -- 设置并发数，和爬取限制
	// 设置请求限制（每秒最多3个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: 3, // 和queue队列同时存在时，用queue控制并发就行。加这个有用，但没必要。默认是0，表示没限制
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})
	// -- 注册 HTML 解析逻辑
	// 线程安全的 去重map, 用于爬某类所有page数据 --
	var comicNamePool sync.Map

	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的 --
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
			switch requestBody.EndNum {
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
			"star", "need_tcp", "cover_need_tcp", "spider_end_status", "download_end_status", "upload_aws_end_status", "upload_baidu_end_status",
			"updated_at"} // 要传updated_at ，upsert必须传, UPDATE()方法不用传，会自动改
		db.DBUpsertBatch(db.DBComic, comicArr, uniqueIndexArr, updateColArr)
		// 7. 重置变量
		comicArr = comicArr[:0]
		moveRepeatComics = make(map[string]string)
		comicNamePool = sync.Map{}
	})

	// -- 添加多个页面到队列中
	// 使用队列控制任务调度（最多并发3个Url）
	q, _ := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})
	// 添加任务到队列
	// for i := 1; i <= requestBody.EndNum; i++ {
	// 	q.AddURL(fullUrl + strconv.Itoa(i))
	// }

	// 测试用
	q.AddURL("http://localhost:8080/test/index.html")

	// 启动对垒
	q.Run(c)

	// -- 插入db
	// 4. 返回结果
	context.JSON(500, "OK")
}

// 爬取分类 By Json,通过人工F12 查看的JSON返回数据，如“热血”“恋爱”“悬疑” -》
/*
参数：
	1. WholeJsonByteData []byte 类型 // 整个json的2进制。 传参，前端传来的json字节数据. 这个值，是某个方法读取 gin.Context后，把读取结果传过来的
	2. map[string]models.ModelMapping 类型 // 传参，表映射关系

返回：
	1. gjsonResult map[string]any 类型 // gjson 提取的字段值。比如comic表，一个对象
	2. error 类型 // 错误信息

作用简单说：
	- 爬取网站上某一种分类,通过Json。如有声书：悬疑、有声书：科幻、有声书：历史等

作用详细说:

核心思路:
	1. 从参数, 拿某个表的映射关系
	2. 通过映射关系, 提取某个表的 所有字段值
	3. 插入db

参考通用思路：
	1. 校验传参
		- 前端参数转成对象
		- 是否需要简单清洗？
		- 校验
		- 分析前端参数含义
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑 - 爬取 - 插入db
		-- 拼接第一页 完整url
		-- new 爬虫对象
		-- 建一个爬虫对象
		-- 设置并发数，和爬取限制
		-- 注册 HTML 解析逻辑
		-- 添加多个爬虫 到到队列中
	5. 返回结果


注意：

使用方式：
1. 先new 一个 映射关系
var comcicFieldMapping = map[string]models.ComicSpiderFieldMapping{
	"name":       {GetFieldPath: "adult.100.meta.title", FiledType: "string"}, // adult.100.meta.title 这样能获取第100个 的内容
	"websiteId":  {GetFieldPath: "websiteId", FiledType: "int"},
	"pornTypeId": {GetFieldPath: "pornTypeId", FiledType: "int"},
	"countryId":  {GetFieldPath: "countryId", FiledType: "int"},
	"typeId":     {GetFieldPath: "typeId", FiledType: "int"},
	"latestChapter":     {GetFieldPath: "adult.100.lastUpdated.episodeTitle", FiledType: "string"},
	"hits":       {GetFieldPath: "adult.100.meta.viewCount", FiledType: "int"},
	"comicUrlApiPath": {GetFieldPath: "adult.100.id", FiledType: "string",
		Transform: func(v any) any {
			id := v.(string)
			return "https://www.toptoon.net/comic/epList/" + id
		}}, // Template 表示模板：能实现拼接"https://www.toptoon.net/comic/epList/" + id
	"coverUrlApiPath":     {GetFieldPath: "adult.100.thumbnail.standard", FiledType: "string"},
	"briefLong":    {GetFieldPath: "adult.100.meta.description", FiledType: "string"},
	"star":         {GetFieldPath: "adult.100.meta.rating", FiledType: "float"},
	"needTcp":      {GetFieldPath: "needTcp", FiledType: "bool"},
	"coverNeedTcp": {GetFieldPath: "coverNeedTcp", FiledType: "bool"},
}

然后调用 func BookTemSpiderTypeByJson(c *gin.Context, map[string]models.ModelMapping)，传入参数
*/
func BookTemSpiderTypeByJson(WholeJsonByteData []byte, modelMapping map[string]models.ModelMapping) (map[string]any, error) {
	// v0.2 方式实现：用通用处理，适合所有表
	// 1. 校验传参
	// -- 判断 jsonByteData 是否空
	if len(WholeJsonByteData) == 0 {
		log.Error("func=BookTemSpiderTypeByJson(爬取JSON).参数 jsonByteData 不能为空. 建议排查步骤: 1. 判断上级传参是否读取出 gin.Context 2.确认前端传参json格式是否对 3. 前端传参json内容是否为空 4. 前端传参json是否少东西")
		return nil, errors.New("func=BookTemSpiderTypeByJson(爬取JSON).参数 jsonByteData 不能为空. 建议排查步骤: 1. 判断上级传参是否读取出 gin.Context 2.确认前端传参json格式是否对 3. 前端传参json内容是否为空 4. 前端传参json是否少东西")
	}
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// 添加这行来打印原始JSON数据
	log.Debug("func=BookTemSpiderTypeByJson(爬取JSON). 前端传参, 原始JSON = ", string(WholeJsonByteData))

	// 使用 gjson 获取字段
	gjsonResult := GetTableFieldValueBySpiderMapping(WholeJsonByteData, modelMapping)

	// 4. 执行核心逻辑

	// 5. 返回结果
	return gjsonResult, nil // 成功

	// v0,1 方式实现：没有用通用处理，仅仅适合comic表
	/*
		// 1. 校验传参
		// 2. 数据清洗

		// 3. 业务逻辑 需要的数据校验 +清洗
		// -- 读取 JSON Body
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
			return
		}

		// 使用 gjson 获取字段
		// -- 外键id，人为提供
		websiteId := gjson.GetBytes(data, "websiteId").Int()   // 网站id
		pornTypeId := gjson.GetBytes(data, "pornTypeId").Int() // 色情类型id
		countryId := gjson.GetBytes(data, "countryId").Int()   // 国家id
		typeId := gjson.GetBytes(data, "typeId").Int()         // 类型id

		// -- 爬取需要的参数
		adultResultArr := gjson.GetBytes(data, "adult").Array()
		adultResultArr100 := adultResultArr[100]

		bookName := adultResultArr100.Get("meta.title").String()
		udpate := adultResultArr100.Get("lastUpdated.episodeTitle").String()                       // 更新到
		hits := adultResultArr100.Get("meta.viewCount").Int()                                      // 点击量
		comicUrlApiPath := "https://www.toptoon.net/comic/epList/" + adultResultArr100.Get("id").String() // 人工给一个链接  https://www.toptoon.net/comic/epList + "/" + id
		coverUrlApiPath := adultResultArr100.Get("thumbnail.standard").String()                           // 封面链接
		breifLong := adultResultArr100.Get("meta.description").String()                            // 简介-long
		star := adultResultArr100.Get("meta.rating").Float()                                       // 评分
		needTcp := 1
		coverNeedTcp := 1

		// -- 循环清洗，空格+繁体 ！！！！！

		// 打印调试
		log.Info("adultResultArr[100] = ", adultResultArr[100])
		log.Info("adultResultArr[100].websiteId = ", websiteId)
		log.Info("adultResultArr[100].pornTypeId = ", pornTypeId)
		log.Info("adultResultArr[100].countryId = ", countryId)
		log.Info("adultResultArr[100].typeId = ", typeId)
		log.Info("adultResultArr[100].bookname = ", bookName)
		log.Info("adultResultArr[100].udpate = ", udpate)
		log.Info("adultResultArr[100].hits = ", hits)
		log.Info("adultResultArr[100].comicUrlApiPath = ", comicUrlApiPath)
		log.Info("adultResultArr[100].coverUrlApiPath = ", coverUrlApiPath)
		log.Info("adultResultArr[100].breifLong = ", breifLong)
		log.Info("adultResultArr[100].breifLstarong = ", star)
		log.Info("adultResultArr[100].needTcp = ", needTcp)
		log.Info("adultResultArr[100].coverNeedTcp = ", coverNeedTcp)

		// 4. 执行核心逻辑
		// 5. 返回结果
		c.JSON(200, gin.H{
			"bookName": bookName,
			"udpate":   udpate,
			"hits":     hits,
		})
	*/
}

// 通过要爬取 表的爬取映射关系, 获取某个表的所有字段值. 现在只能实现comic表，不能通用！！！！
// 最理想状态：只提取一个对象的所有信息
/*
参数：
	1 jsonByteData []byte json数据
	2 spiderMapping map[string]models.ComicSpiderFieldMapping 爬取映射关系

返回:
	1. map[string]any gjsonResult gjson解析JSON后的结果
*/
func GetTableFieldValueBySpiderMapping(jsonByteData []byte, spiderMapping map[string]models.ModelMapping) map[string]any {
	result := make(map[string]any) // 存放结果

	// 提取字段
	for key, fieldMapping := range spiderMapping {
		v := gjson.GetBytes(jsonByteData, fieldMapping.GetFieldPath)

		var value any

		// --- 简单处理 Transform ---
		if fieldMapping.Transform != nil {
			// Transform 优先执行
			value = fieldMapping.Transform(v.Value())
		} else {
			// 根据 FiledType 取值
			switch fieldMapping.FiledType {
			case "string":
				value = v.String()
			case "int":
				value = v.Int()
			case "float":
				value = v.Float()
			case "bool":
				value = v.Bool()
			case "array":
				value = v.Array()
			case "time":
				// value = v.Time() // 转成日期+时间 time.Time 格式 YYYY-MM-DD HH:MM:SS， 不好用，弃用
				// 使用time.Parse替代v.Time()，因为v.Time()无法解析"2025-11-18 22:00:00"这种格式
				value, _ = time.Parse("2006-01-02 15:04:05", v.String())

			default:
				value = v.Value() // fallback ?啥意思
			}
		}

		result[key] = value
	}

	return result
}

// 获取某一个类型，第一页的完整 url.如有声书：悬疑、言情分类
/*
参数：
	1 needTcp bool 是否需要带http
	2 needHttps bool 是否需要带https
	3 websitePrefix string 网站前缀，现在想的是最后不带/，如：www.manhuagui.com
	4 apiPath string 接口url/，如：list/c1-p, 要包括第1页
	5 paramArr []map[string]string 参数， // ?后面的,参数值。如：[{"type": "1"}, {"complete": "1"}]

作用简单说：
	- 拼接 类型第1页 完整url

作用详细说:

核心思路:

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑
		-- new url对象
		-- 设置域名
		-- 设置接口路径
		-- 设置参数
	5. 返回结果

完整url 示例：
全部: https://kxmanhua.com/manga/library?type=0&complete=1&page=1&orderby=1
3D:  https://kxmanhua.com/manga/library?type=1&complete=1&page=1&orderby=1
韩漫: https://kxmanhua.com/manga/library?type=2&complete=1&page=1&orderby=1
日漫: https://kxmanhua.com/manga/library?type=3&complete=1&page=1&orderby=1
*/
func GetOneTypeFirstPageFullURL(needTcp, needHttps bool, websitePrefix, apiPath string, paramArr []map[string]string) string {
	// v0.1 的写法。自己拼接，没有用到url这个库
	/*
		// 1. 校验传参
		// 2. 数据清洗
		// 3. 业务逻辑 需要的数据校验 +清洗
		// 4. 执行核心逻辑
		protocol := "" // 协议
		// -- 判断协议头
		if needTcp { // 需要带http 或 https
			if needHttps { // 需要https
				protocol = "https://"
			}
			protocol = "http://"
		}

		// 5. 返回结果
		return protocol + prefix + apiUrl
	*/

	// v0.2 的写法。用url这个库
	// 1. 校验传参
	// 2. 数据清洗
	apiPath = strings.TrimSpace(apiPath) // 去除前后空格

	// 3. 业务逻辑 需要的数据校验 +清洗
	// 4. 执行核心逻辑
	// -- new url对象
	u := url.URL{}

	// -- 设置协议头 https:// 或 http://
	if needTcp { // 需要带http 或 https
		if needHttps { // 需要https
			u.Scheme = "https"
		}
		u.Scheme = "http"
	}

	// -- 设置域名
	u.Host = websitePrefix

	// -- 设置接口路径
	u.Path = apiPath

	// -- 设置参数，要确保，按传参顺序赋值。下面方法可以保证
	paramsObj := url.Values{}
	for _, param := range paramArr {
		for k, v := range param {
			paramsObj.Set(k, v)
		}
	}
	u.RawQuery = paramsObj.Encode()

	// 5. 返回结果
	return u.String()
}

// 通用函数 根据 model 模型里的tag, 爬取到的json结果 result -> 转成成 模型对象. AI给的真的好用，仔细了解下这个方法！！
func MapByTag(result map[string]any, out any) {
	// v0.1 的写法。能适配大部分场景，comic里加了AuthorArr之后就不能用了
	// /*
	t := reflect.TypeOf(out).Elem()
	v := reflect.ValueOf(out).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		spiderKey := f.Tag.Get("spider") // 获取 tag

		if spiderKey == "" {
			continue
		}

		if val, ok := result[spiderKey]; ok {
			field := v.Field(i)
			if field.CanSet() {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		}
	}
	// */

	// v0.2 的写法。支持切片类型和结构体类型的转换
	/*
		t := reflect.TypeOf(out).Elem()
		v := reflect.ValueOf(out).Elem()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			spiderKey := f.Tag.Get("spider") // 获取 tag

			if spiderKey == "" {
				continue
			}

			if val, ok := result[spiderKey]; ok {
				field := v.Field(i)
				if field.CanSet() {
					// 处理不同类型的转换
					err := convertAndSet(val, field)
					if err != nil {
						// 如果转换失败，尝试使用原来的方法
						if reflect.ValueOf(val).Type().ConvertibleTo(field.Type()) {
							field.Set(reflect.ValueOf(val).Convert(field.Type()))
						} else {
							// 记录错误但继续处理其他字段
							fmt.Printf("字段转换失败: %v\n", err)
						}
					}
				}
			}
		}
	*/
}

// 通用转换函数，处理各种类型转换
func convertAndSet(source any, targetField reflect.Value) error {
	sourceValue := reflect.ValueOf(source)
	targetType := targetField.Type()

	// 处理指针类型
	if targetType.Kind() == reflect.Ptr {
		// 创建新对象
		newValue := reflect.New(targetType.Elem())
		// 递归处理
		err := convertAndSet(source, newValue.Elem())
		if err != nil {
			return err
		}
		targetField.Set(newValue)
		return nil
	}

	// 处理切片类型
	if targetType.Kind() == reflect.Slice {
		return handleSliceConversion(sourceValue, targetField)
	}

	// 处理结构体类型
	if targetType.Kind() == reflect.Struct {
		return handleStructConversion(sourceValue, targetField)
	}

	// 处理基本类型
	return handleBasicTypeConversion(sourceValue, targetField)
}

// 处理切片转换
func handleSliceConversion(sourceValue reflect.Value, targetField reflect.Value) error {
	// 检查源是否是切片/数组
	if sourceValue.Kind() != reflect.Slice && sourceValue.Kind() != reflect.Array {
		return fmt.Errorf("源不是切片/数组类型，无法转换为切片")
	}

	// 获取目标元素类型
	elemType := targetField.Type().Elem()

	// 创建新切片
	newSlice := reflect.MakeSlice(targetField.Type(), 0, sourceValue.Len())

	// 遍历源切片
	for i := 0; i < sourceValue.Len(); i++ {
		sourceElem := sourceValue.Index(i)

		// 创建新元素
		var newElem reflect.Value
		if elemType.Kind() == reflect.Ptr {
			newElem = reflect.New(elemType.Elem())
			err := convertAndSet(sourceElem.Interface(), newElem.Elem())
			if err != nil {
				return fmt.Errorf("切片元素转换失败: %v", err)
			}
			newSlice = reflect.Append(newSlice, newElem)
		} else {
			// 对于非指针类型，直接创建元素
			newElem = reflect.New(elemType).Elem()
			err := convertAndSet(sourceElem.Interface(), newElem)
			if err != nil {
				return fmt.Errorf("切片元素转换失败: %v", err)
			}
			newSlice = reflect.Append(newSlice, newElem)
		}
	}

	// 设置字段值
	targetField.Set(newSlice)
	return nil
}

// 处理结构体转换
func handleStructConversion(sourceValue reflect.Value, targetField reflect.Value) error {
	// 如果源是map[string]any，可以递归调用MapByTag
	if sourceMap, ok := sourceValue.Interface().(map[string]any); ok {
		MapByTag(sourceMap, targetField.Addr().Interface())
		return nil
	}

	// 如果源是gjson.Result，需要先转换为map[string]any
	if gjsonResult, ok := sourceValue.Interface().(gjson.Result); ok {
		// 尝试将gjson.Result转换为map
		if gjsonResult.IsObject() {
			resultMap := make(map[string]any)
			gjsonResult.ForEach(func(key, value gjson.Result) bool {
				resultMap[key.String()] = value.Value()
				return true
			})
			MapByTag(resultMap, targetField.Addr().Interface())
			return nil
		}
	}

	// 其他情况尝试直接转换
	if sourceValue.Type().ConvertibleTo(targetField.Type()) {
		targetField.Set(sourceValue.Convert(targetField.Type()))
		return nil
	}

	return fmt.Errorf("无法将 %v 转换为 %v", sourceValue.Type(), targetField.Type())
}

// 处理基本类型转换
func handleBasicTypeConversion(sourceValue reflect.Value, targetField reflect.Value) error {
	// 尝试直接转换
	if sourceValue.Type().ConvertibleTo(targetField.Type()) {
		targetField.Set(sourceValue.Convert(targetField.Type()))
		return nil
	}

	// 特殊处理：gjson.Result到基本类型的转换
	if gjsonResult, ok := sourceValue.Interface().(gjson.Result); ok {
		switch targetField.Kind() {
		case reflect.String:
			targetField.SetString(gjsonResult.String())
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			targetField.SetInt(gjsonResult.Int())
			return nil
		case reflect.Float32, reflect.Float64:
			targetField.SetFloat(gjsonResult.Float())
			return nil
		case reflect.Bool:
			targetField.SetBool(gjsonResult.Bool())
			return nil
		}
	}

	return fmt.Errorf("无法将 %v 转换为 %v", sourceValue.Type(), targetField.Type())
}

// 插入爬取后的表数据，方法小写，不导出。非通用方法，插入的核心实现逻辑
/*
参数：
	1 tableName string 表名
	2 gjsonResultArr []map[string]any 爬出来的对象，数组
*/
func upsertSpiderTableData(tableName string, gjsonResultArr []map[string]any) error {
	log.Debug("func=upsertSpiderTableData(插入爬取表数据) tableName=", tableName)

	// 1. 校验传参
	// 2. 数据清洗
	// 3. 业务逻辑 需要的数据校验 +清洗
	// 4. 执行核心逻辑

	// 插入数据库
	switch tableName {
	case "comic":
		// -- 准备初始化
		var comicArr []*models.ComicSpider

		// -- 把爬到的 gjsonResultArr 转成 表对象 数组
		for _, gjsonResult := range gjsonResultArr {
			// 准备插入参数, 循环清洗，空格+繁体 --
			comic := &models.ComicSpider{}
			MapByTag(gjsonResult, comic) // 爬取json内容，赋值给 comic对象

			// 数据清洗 (空格，转简体) --
			comic.TrimSpaces()  // 调下自己的方法，去空格
			comic.Trad2Simple() // 调用自己实现的接口方法，转简体

			// 添加到 comicArr数组 --
			comicArr = append(comicArr, comic)
		}

		// -- 数据校验，看有没有 不好用/错误的数据
		for i, comic := range comicArr {
			// 如果 comic.name是空，那这批数据不能用
			if comic.Name == "" {
				return errors.New("这批数据不能用, comic.name是空, 需修改前端传参json, index= " + strconv.Itoa(i))
			}
			log.Infof("爬取后将插入的comic, index =%v, comic = %v ", i, comic)
		}

		// -- 批量插入
		err := db.DBUpsertBatch(db.DBComic, comicArr, tableComicUniqueIndexArr, tableComicUpdateColArr)
		if err != nil {
			return err
		}
	case "author":
		// -- 准备初始化
		var authorArr []*models.Author

		// -- 把爬到的 gjsonResultArr 转成 表对象 数组
		log.Info("------- gjsonResultArr = ", gjsonResultArr)
		for _, gjsonResult := range gjsonResultArr {

			// 直接从gjsonResult获取name字段
			// 准备插入参数, 循环清洗，空格+繁体 --
			author := &models.Author{}
			MapByTag(gjsonResult, author) // 爬取json内容，赋值给 author 对象
			// 从gjsonResult中获取name字段
			if nameValue, ok := gjsonResult["name"]; ok {
				author.Name = fmt.Sprintf("%v", nameValue) // 获取name字段
				log.Debug("-- author = ", author)

				// 数据清洗 (空格，转简体) --
				author.TrimSpaces()  // 调下自己的方法，去空格
				author.Trad2Simple() // 调用自己实现的接口方法，转简体

				// 添加到 authorArr数组 --
				authorArr = append(authorArr, author)
			}

		}

		// -- 数据校验，看有没有 不好用/错误的数据
		for i, author := range authorArr {
			// 如果 comic.name是空，那这批数据不能用
			if author.Name == "" {
				return errors.New("这批数据不能用, author.name是空, 需修改前端传参json, index= " + strconv.Itoa(i))
			}
			log.Infof("爬取后将插入的 author, index =%v, author = %v ", i, author)
		}

		// -- 批量插入
		err := db.DBUpsertBatch(db.DBComic, authorArr, tableAuthorUniqueIndexArr, tableAuthorUpdateColArr)
		if err != nil {
			return err
		}
	default:
		return errors.New("未知的表名")
	}

	// 5. 返回结果
	return nil // 成功

}

// -- 方法 ------------------------------------------- end -----------------------------------
