/**
V1.5 版本：爬取通用方法
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
	"strconv"
	"strings"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/errs"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util"
	"study-spider-manhua-gin/src/util/langutil"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"gorm.io/gorm"
)

// ------------------------------------------- 方法 -------------------------------------------

// 获取1个book所有chapter, 用colly, 通过mapping
/*

5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
	步骤1: 找到目标网站
	步骤2: 爬取
	步骤3: 提取数据 <- 往上是本方法
	步骤4: 数据清洗/ 未爬到的字段赋值
	步骤5: 验证爬取数据 准确性
	步骤6: 数据库插入

参数:
	2. mapping map[string]models.ModelMapping 爬取映射关系
	3. bookID

返回:
主表数组
作用简单说：
*/
func GetOneBookAllChapterByCollyMappingV1_5[T any](mapping map[string]models.ModelHtmlMapping, bookId int) ([]T, error) {
	// 步骤2: 爬取
	// 1. gjson 读取 前端 JSON 里 有用内容

	// 2. 爬虫相关
	// -- 建一个爬虫对象
	c := colly.NewCollector()

	// -- 设置并发数，和爬取限制
	// 设置请求限制（每秒最多3个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: 3, // 和queue队列同时存在时，用queue控制并发就行。加这个有用，但没必要。默认是0，表示没限制
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// 步骤3: 提取数据
	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的 --
	var chapterArr []T // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider{}
	// 遍历一个book, 每个chapter
	c.OnHTML(".chapter_list a", func(e *colly.HTMLElement) {
		// 0. 处理异常内容
		// -- 处理 ”休刊公告“
		oneChapterStr, _ := langutil.TraditionalToSimplified(e.Text)
		if strings.Contains(oneChapterStr, "休刊") {
			return // ✅ 这个 return 只从匿名函数返回，不会影响 GetOneBookObjByCollyMapping 函数
		}

		// 1. 获取能获取到的
		log.Debug("匹配 .chapter_list a = ", e.Text)
		// -- 创建对象comic
		var chapterT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		log.Info("通过mapping规则,爬取结果 result = ", result)
		if result != nil {
			// 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象
			MapByTag(result, &chapterT)
			log.Debugf("映射后的 chapter 对象, 还未清洗: %+v", chapterT)
		}
		// 2. 放到chapterArr里
		chapterArr = append(chapterArr, any(chapterT).(T))
	})

	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		if r == nil {
			// 网络层错误（DNS / timeout / TLS）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 网络层错误（DNS / timeout / TLS）, err: ", err)
			return
		}

		switch {
		case r.StatusCode >= 400 && r.StatusCode < 500:
			// 4xx：客户端错误（参数错误、被封、资源不存在）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 客户端错误（参数错误、被封、资源不存在）, err: ", err)

		case r.StatusCode >= 500 && r.StatusCode < 600:
			// 5xx：服务端错误（可重试）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 服务端错误（可重试）, err: ", err)

		default:
			// 其他非常规状态码
			// 可选重试
		}
	})

	// -- 添加多个页面到队列中
	// 使用队列控制任务调度（最多并发3个Url）
	q, _ := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})

	// -- 添加任务到队列
	// 通过bookId 查询，并拼接出 book的 url --
	book, err := db.DBFindOneByFieldV1_5[models.ComicSpider](db.DBComic, "id", bookId) // 通过bookId 查website
	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取book对象失败, err: ", err)
		return nil, err
	}
	websiteId := book.WebsiteId
	website, err := db.DBFindOneByFieldV1_5[models.Website](db.DBComic, "id", websiteId)
	if err != nil || website == nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取website对象失败, err: ", err)
		return nil, errors.New("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取website对象失败")
	}

	// 步骤1: 找到目标网站
	// fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, book.ComicUrlApiPath, nil) // 完整爬取 url
	fullUrl := GetSpiderFullUrl(false, "localhost:8080", "/test/kxmanhua/spiderChapter/社团学姐.html", nil) // 完整爬取 url，本地测试
	log.Info("生成的book 爬取 fullURl = ", fullUrl)
	// 打算使用 GET 请求校验 URL 可达性，通过后才加入抓取队列。爬取的一般都是get请求， 就用get请求下。但实际不用 用c.OnError() 就能有类似效果

	// 再添加到队列 --

	q.AddURL(fullUrl) // 只能添加1个

	// 测试用 - 添加任务到队列
	// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/社团学姐.html") // 章节url
	// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/1.html") // 章节url

	// 启动对垒
	q.Run(c)
	return chapterArr, nil
}

// 把爬取 oneBookAllChapter 分成2部分。爬取部分 + 插入部分
// 插入部分 - 插入单个书的 chapter
/*
5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值 <- 本方法
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
			- 6.1 插入 章节
			- 6.2 更新 book 表stats字段
参数:

返回 插入成功总数
 error
*/
func SpiderOneBookAllChapter_UpsertPart(websiteName string, bookId int, chapterArr []models.ChapterSpider) (int, error) {
	// 初始化
	// 获取请求种用的到的参数 website + bookId
	okTotal := 0 // 插入成功总数

	// 步骤4: 数据清洗/ 未爬到的字段赋值
	// -- 赋值上下文参数 + 数据清洗。（赋值上下文参数：是吧方法传参，给对象赋值。数据清洗：设置-爬取字段，或者默认数据）
	for i := range chapterArr {
		// -赋值 上下文传参。如parentId (非数据清洗业务，放在这里)
		chapterArr[i].ParentId = int(bookId) // 父id
		// -数据清洗
		chapterArr[i].DataClean() // 数据清洗
		log.Debug("清洗完数据 chapter = ", chapterArr[i])
	}

	// -- 检测下爬到的数据，有没有重复数据，需要注意下。只要判单章节号码 chapter_num 就可以了
	var spiderChapterNumArr []int // 爬到的章节号 arr
	for _, chapter := range chapterArr {
		spiderChapterNumArr = append(spiderChapterNumArr, chapter.ChapterNum)
	}
	if util.HasDuplicate(spiderChapterNumArr) { // 判断有重复
		log.Warn("爬取1本书 AllChapter, 爬到的章节有重复, 要注意下")
	}

	// 步骤6: 数据库插入
	// 6.1 插入 章节. upsert chapter
	// 获取配置 --
	webCfg := config.CfgSpiderYaml.Websites[websiteName]
	if webCfg == nil {
		// c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 %s 的配置", website)}) // 返回错误
		return 0, nil
	}
	log.Debug("------- webCfg = ", webCfg)

	// 获取 one_book_all_chapter 阶段配置
	stageCfg := webCfg.Stages["one_book_all_chapter"]
	if stageCfg == nil {
		// c.JSON(400, gin.H{"error": "func=爬 爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_book_all_chapter 阶段的配置"}) // 返回错误
		return 0, errors.New("func=爬 爬oneTypeAllBookArr V1.5, stageCfg 为空")
	}
	log.Debug("------- stageCfg = ", stageCfg)

	// 批量插入db chapter
	if len(chapterArr) == 0 {
		return 0, errors.New("func=爬 爬oneTypeAllBookArr V1.5, chapterArr 为空")
	}

	err := db.DBUpsertBatch(db.DBComic, chapterArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)

	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 批量插入db chapter 失败, err: ", err)
		// c.JSON(500, gin.H{"error": "批量插入db chapter 失败"}) // 返回错误
	}

	// 6.2 更新 book 表stats字段
	// 4. 考虑爬一次comic相关的，比如最后一章更新时间、评分等。现在没实现。真用到再实现

	// 5. 更新comic.因为有的需要有chapter数据，才可以。比如 最后一章id, 章节总数
	// -- 需要更新：  \ 最后章节名称 \ 总章节数 （假如需要爬 comic想的，那就得爬完插入chapter之后，再爬一次comic相关的）
	// 找到最后一章，从chapter里获取需要内容。找 chapterNUm=9999就行
	lastChapter, err := db.DBFindOneByField[models.ChapterSpider]("chapter_num", 9999)
	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 找到最后一章失败, err: ", err)
	}

	// -- 创建 comic_spider_stats对象
	var comicSpiderStats models.ComicSpiderStats
	comicSpiderStats.ComicId = int(bookId)
	comicSpiderStats.LatestChapterId = &lastChapter.Id    // 最后章节id
	comicSpiderStats.LatestChapterName = lastChapter.Name // 最后章节名称

	totalChapterDbRealUpsert, err := db.DBCountByField[models.ChapterSpider](db.DBComic, "parent_id", bookId) // db里真实插入 章节个数
	errorutil.ErrorPrint(err, "爬取oneBookAllChapter, 插入chapter_spider表后, 查询总插入数出错, err = ")
	comicSpiderStats.TotalChapter = totalChapterDbRealUpsert // 总章节数，从数据库查的

	// -- 更新 comic_spider_stats
	log.Info("update comicSpiderStats = ", comicSpiderStats)
	log.Info("update comicSpiderStats.ComicId = ", comicSpiderStats.ComicId)
	log.Info("update comicSpiderStats.LatestChapterId = ", *comicSpiderStats.LatestChapterId)
	err = db.DBUpdate(db.DBComic, &comicSpiderStats, stageCfg.UpdateParentStats.UniqueKeys, stageCfg.UpdateParentStats.UpdateKeys)
	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml, 更新comic_spider_stats失败, err: ", err)
		// c.JSON(500, gin.H{"error": "更新comic_spider_stats失败"}) // 返回错误
		return 0, err
	}

	okTotal = totalChapterDbRealUpsert
	log.Infof("插入成功 %v 条", okTotal)
	return okTotal, nil // 一切正常
}

// 根据条件生成 爬取完整url
/*
参数
	isHttps: 是否需要 https协议头

示例：
fullUrl := GetSpiderFullUrl(false, "localhost:8080", "/test/kxmanhua/spiderChapter/社团学姐.html", nil) // 完整爬取 url，本地测试
fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, book.ComicUrlApiPath, nil) // 完整爬取 url
*/
func GetSpiderFullUrl(isHttps bool, websitePrefix, apiPath string, paramArr []map[string]string) string {
	u := url.URL{} // 初始化

	// -- 设置协议头 https:// 或 http://
	u.Scheme = "http"
	if isHttps {
		u.Scheme = "https"
	}

	// -- 设置域名. 没有Port参数，要想写port,在host里写。例如：localhost:8888
	u.Host = websitePrefix

	// -- 设置接口路径
	u.Path = apiPath

	// -- 设置参数，要确保，按传参顺序赋值。下面方法可以保证。按传参顺序赋值意思就是： 比如传的顺序 是 [b] [a] [c]， 生成的链接是 b=?a=?c=?
	if len(paramArr) > 0 { // 没传参数情况下
		return u.String()
	}

	paramsObj := url.Values{}
	for _, param := range paramArr {
		for k, v := range param {
			paramsObj.Set(k, v)
		}
	}
	u.RawQuery = paramsObj.Encode()

	// -- 返回
	return u.String()
}

// 获取多个book所有chapter, 用colly, 通过mapping V1,只适用于 kxmanhua , 不通用
/*

5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
	步骤1: 找到目标网站
	步骤2: 爬取
	步骤3: 提取数据 <- 往上是本方法
	步骤4: 数据清洗/ 未爬到的字段赋值
	步骤5: 验证爬取数据 准确性
	步骤6: 数据库插入

参数:
	2. mapping map[string]models.ModelMapping 爬取映射关系
	2. websiteName string 网站名称
	3. bookIdArr

返回:
map id -> 数组

主表数组
作用简单说：
*/
func GetManyBookAllChapterByCollyMappingV1_5_V1_OnlyForKxmanhua[T any](mapping map[string]models.ModelHtmlMapping, websiteName string, bookIdArr []int) (map[int][]T, error) {
	// 初始化
	funcName := "GetManyBookAllChapterByCollyMappingV1_5"
	// bookId 和 fullUrl 映射关系, key 是url
	bookIdFullUrlMapKeyUrl := make(map[string]int)

	// 步骤2: 爬取

	// 2. 爬虫相关
	// -- 建一个爬虫对象
	c := colly.NewCollector()

	// -- 设置并发数，和爬取限制
	// 设置请求限制（每秒最多3个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: 3, // 和queue队列同时存在时，用queue控制并发就行。加这个有用，但没必要。默认是0，表示没限制
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// 步骤3: 提取数据
	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的 --
	var oneBookChapterArr []T                     // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider{}
	var manyBookChapterArrMap = make(map[int][]T) //所有book 的chapte数组 map -》 bookId 对应 所有章节
	// 遍历一个book, 每个chapter
	// c.OnHTML(".chapter_list a", func(e *colly.HTMLElement) {
	StagesCfg := config.CfgSpiderYaml.Websites[websiteName].Stages["one_book_all_chapter"]
	everyChapterSelectStr := StagesCfg.Crawl.Selectors["item"].(string) // 每个chapter 选择器
	c.OnHTML(everyChapterSelectStr, func(e *colly.HTMLElement) {
		// 0 初始化
		currentUrl := e.Request.URL.String() // 获取当前正在爬取的 URL

		// 0. 处理异常内容
		// -- 处理 ”休刊公告“
		oneChapterStr, _ := langutil.TraditionalToSimplified(e.Text)
		if strings.Contains(oneChapterStr, "休刊") || strings.Contains(oneChapterStr, "停刊") {
			return // ✅ 这个 return 只从匿名函数返回，不会影响 GetOneBookObjByCollyMapping 函数
		}

		// 1. 获取能获取到的
		log.Debug("匹配 .chapter_list a = ", e.Text)
		// -- 创建对象comic
		var chapterT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		log.Infof(" bookIdArr=%v,当前爬取url=%v 通过mapping规则,爬取结果 result = %v", bookIdArr, currentUrl, result)
		if result != nil {
			// 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象
			MapByTag(result, &chapterT)
			log.Debugf("映射后的 chapter 对象, 还未清洗: %+v", chapterT)
		}
		// 2. 放到 oneBookChapterArr 里
		oneBookChapterArr = append(oneBookChapterArr, any(chapterT).(T))
	})

	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		if r == nil {
			// 网络层错误（DNS / timeout / TLS）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 网络层错误（DNS / timeout / TLS）, err: ", err)
			return
		}

		switch {
		case r.StatusCode >= 400 && r.StatusCode < 500:
			// 4xx：客户端错误（参数错误、被封、资源不存在）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 客户端错误（参数错误、被封、资源不存在）, err: ", err)

		case r.StatusCode >= 500 && r.StatusCode < 600:
			// 5xx：服务端错误（可重试）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 服务端错误（可重试）, err: ", err)

		default:
			// 其他非常规状态码
			// 可选重试
		}
	})

	// 成功爬完1页，回调
	c.OnScraped(func(r *colly.Response) {
		// 把 oneBookAllChapter 加到 要返回的map 中去
		currentUrl := r.Request.URL.String()         // 获取当前正在爬取的 URL
		bookId := bookIdFullUrlMapKeyUrl[currentUrl] // 通过url 获取到bookId
		manyBookChapterArrMap[bookId] = oneBookChapterArr

		oneBookChapterArr = nil // 或者 oneBookChapterArr = []T{} 重置切片
	})

	// -- 添加多个页面到队列中
	// 使用队列控制任务调度（最多并发3个Url）
	q, _ := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})

	// -- 添加任务到队列
	for _, bookId := range bookIdArr {
		// 通过bookId 查询，并拼接出 book的 url --
		book, err := db.DBFindOneByFieldV1_5[models.ComicSpider](db.DBComic, "id", bookId) // 通过bookId 查website
		if err != nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取book对象失败, err: ", err)
			return nil, err
		}
		websiteId := book.WebsiteId
		website, err := db.DBFindOneByFieldV1_5[models.Website](db.DBComic, "id", websiteId)
		if err != nil || website == nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取website对象失败, err: ", err)
			return nil, errors.New("func= DispatchApi_OneBookAllChapterByHtml, 通过bookId查询website信息失败,因为通过bookId获取website对象失败")
		}

		// 步骤1: 找到目标网站
		fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, book.ComicUrlApiPath, nil) // 完整爬取 url
		// apiUrlPath := fmt.Sprintf("/test/kxmanhua/spiderChapter/%d.html", bookId)
		// fullUrl := GetSpiderFullUrl(false, "localhost:8080", apiUrlPath, nil) // 完整爬取 url，本地测试
		log.Info("生成的book 爬取 fullURl = ", fullUrl)
		bookIdFullUrlMapKeyUrl[fullUrl] = bookId // 记录bookId 和 fullUrl 的映射关系
		// 打算使用 GET 请求校验 URL 可达性，通过后才加入抓取队列。爬取的一般都是get请求， 就用get请求下。但实际不用 用c.OnError() 就能有类似效果

		// 再添加到队列 --
		q.AddURL(fullUrl) // 只能添加1个

		// 测试用 - 添加任务到队列
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/社团学姐.html") // 章节url
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/1.html") // 章节url
	}
	log.Infof("func=%v, bookId 和 请求fullUrl映射关系 bookIdFullUrlMapKeyUrl = %v", funcName, bookIdFullUrlMapKeyUrl)

	// 启动对垒
	q.Run(c)
	return manyBookChapterArrMap, nil
}

// 把爬取 manyBookAllChapter 分成2部分。爬取部分 + 插入部分. V1 实现，只适用于kxmanhua，如果章节页面有其它表数据(comic,comic_stats, author)，不会同步更新
// 插入部分 - 插入多个书的 chapter
/*
5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值 <- 本方法
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
			- 6.1 插入 章节
			- 6.2 更新 book 表stats字段
参数:
manyBookChapterArrMap []

返回 插入成功总数
 error
*/
func SpiderManyBookAllChapter_UpsertPart_V1_OnlyForKxmanhua(websiteName string, manyBookChapterArrMap map[int][]models.ChapterSpider) (int, error) {
	// 初始化
	okTotal := 0 // 插入成功总数
	// funcName := "SpiderManyBookAllChapter_UpsertPart"

	// 异常处理
	if len(manyBookChapterArrMap) == 0 {
		return 0, errors.New("func=爬 爬oneTypeAllBookArr V1.5, manyBookChapterArrMap 为空")
	}

	for bookId, oneBookChapterArr := range manyBookChapterArrMap {
		// 步骤4: 数据清洗/ 未爬到的字段赋值
		// -- 赋值上下文参数 + 数据清洗。（赋值上下文参数：是吧方法传参，给对象赋值。数据清洗：设置-爬取字段，或者默认数据）
		for i := range oneBookChapterArr {
			// -赋值 上下文传参。如parentId (非数据清洗业务，放在这里)
			// oneBookChapterArr[i].ParentId = int(bookIdArr[index]) // delete - 弃用。会导致 comic和chapter内容对不上!!!. 父id，应该从manyBookChapterArr 里拿，这里是最准确的。因为bookIdArr 是从小到大，但爬出来的 manyBookChapterArr 是按id 随机的。容易导致：comic和chapter的 真实章节 对不上
			oneBookChapterArr[i].ParentId = bookId // 父id，应该从manyBookChapterArr 里拿，这里是最准确的。因为bookIdArr 是从小到大，但爬出来的 manyBookChapterArr 是按id 随机的。容易导致：comic和chapter的 真实章节 对不上
			// -数据清洗
			oneBookChapterArr[i].DataClean() // 数据清洗
			log.Debug("清洗完数据 chapter = ", oneBookChapterArr[i])
		}

		// -- 检测下爬到的数据，有没有重复数据，需要注意下。只要判单章节号码 chapter_num 就可以了
		var spiderChapterNumArr []int // 爬到的章节号 arr
		for _, chapter := range oneBookChapterArr {
			spiderChapterNumArr = append(spiderChapterNumArr, chapter.ChapterNum)
		}
		// log.Warn("---------- delete 判断去重数组 spiderChapterNumArr= ", spiderChapterNumArr)
		if util.HasDuplicate(spiderChapterNumArr) { // 判断有重复
			log.Warn("爬取1本书 AllChapter, 爬到的章节号码 有重复, 要注意下, bookId = ", bookId)
		}

		// 步骤6: 数据库插入
		// 6.1 插入 章节. upsert chapter
		// 获取配置 --
		webCfg := config.CfgSpiderYaml.Websites[websiteName]
		if webCfg == nil {
			// c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 %s 的配置", website)}) // 返回错误
			return 0, nil
		}
		log.Debug("------- webCfg = ", webCfg)

		// 获取 one_book_all_chapter 阶段配置
		stageCfg := webCfg.Stages["one_book_all_chapter"]
		if stageCfg == nil {
			// c.JSON(400, gin.H{"error": "func=爬 爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_book_all_chapter 阶段的配置"}) // 返回错误
			return 0, errors.New("func=爬 爬oneTypeAllBookArr V1.5, stageCfg 为空")
		}
		log.Debug("------- stageCfg = ", stageCfg)

		// 批量插入db chapter
		if len(oneBookChapterArr) == 0 {
			return 0, errors.New("func=爬 爬oneTypeAllBookArr V1.5, chapterArr 为空")
		}

		err := db.DBUpsertBatch(db.DBComic, oneBookChapterArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)

		if err != nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 批量插入db chapter 失败, err: ", err)
			// c.JSON(500, gin.H{"error": "批量插入db chapter 失败"}) // 返回错误
		}

		// 6.2 更新 book 表stats字段
		// 4. 考虑爬一次comic相关的，比如最后一章更新时间、评分等。现在没实现。真用到再实现

		// 5. 更新comic.因为有的需要有chapter数据，才可以。比如 最后一章id, 章节总数
		// -- 需要更新：  \ 最后章节名称 \ 总章节数 （假如需要爬 comic想的，那就得爬完插入chapter之后，再爬一次comic相关的）
		// 找到最后一章，从chapter里获取需要内容。方案1：通过bookID + 号码=9999 找到最后一章 x 因为有的没有叫"最终话"开头的 方案2: 找parent_id + chapter_real_sort_num 最大的那个数
		// 弃用，如果没有 "最终话开头"，会报错lastChapter, err := db.DBFindOneByMapCondition[models.ChapterSpider](map[string]any{"parent_id": bookIdArr[index], "chapter_real_sort_num": 9999})
		lastChapter, err := db.DBFindOneV2[models.ChapterSpider](db.DBComic,
			db.WithWhere("parent_id = ?", bookId),
			db.WithOrder("chapter_real_sort_num DESC"),
			db.WithLimit(1),
		)

		if err != nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 找到最后一章失败, err: ", err)
		}

		// -- 创建 comic_spider_stats对象
		var comicSpiderStats models.ComicSpiderStats
		comicSpiderStats.ComicId = bookId
		comicSpiderStats.LatestChapterId = &lastChapter.Id    // 最后章节id
		comicSpiderStats.LatestChapterName = lastChapter.Name // 最后章节名称

		// totalChapterDbRealUpsert, err := db.DBCountByField[models.ChapterSpider](db.DBComic, "parent_id", bookIdArr[index]) // 最新版弃用，db里真实插入 章节个数，之前写法，不判断 real_sort_num !=0的情况
		totalChapterDbRealUpsert, err := db.DBCountV2[models.ChapterSpider](db.DBComic, db.WithWhere("parent_id = ? AND chapter_real_sort_num != 0", bookId)) // db里真实插入 章节个数，找parent_id=x. and chapter_real_sort_num !=0
		errorutil.ErrorPrint(err, "爬取oneBookAllChapter, 插入chapter_spider表后, 查询总插入数出错, err = ")
		comicSpiderStats.TotalChapter = totalChapterDbRealUpsert // 总章节数，从数据库查的
		okTotal += totalChapterDbRealUpsert

		// -- 更新 comic_spider_stats
		log.Info("update comicSpiderStats = ", comicSpiderStats)
		log.Info("update comicSpiderStats.ComicId = ", comicSpiderStats.ComicId)
		log.Info("update comicSpiderStats.LatestChapterId = ", *comicSpiderStats.LatestChapterId)
		err = db.DBUpdate(db.DBComic, &comicSpiderStats, stageCfg.UpdateParentStats.UniqueKeys, stageCfg.UpdateParentStats.UpdateKeys)
		if err != nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml, 更新comic_spider_stats失败, err: ", err)
			// c.JSON(500, gin.H{"error": "更新comic_spider_stats失败"}) // 返回错误
			return 0, err
		}
	}

	log.Infof("插入成功 %v 条", okTotal)
	return okTotal, nil // 一切正常
}

// 爬取manybook allChapter V1实现。只适用与 kxmanhua，不能爬章节时，同时处理 能爬到的表数据。如：comic、comic_stats、authoer
// 把爬取+插入放到1个方法，且和 gin.context 解耦
/*
返回：
	okTotal
	error
*/
func SpiderManyBookAllChapter2DB_V1_OnlyForKxmanhua(websiteName string, bookIdArr []int) (int, error) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "SpiderManyBookAllChapter2DB_V1_OnlyForKxmanhua"
	var funcErr error

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	log.Infof("func=%v, 要爬的bookId = %v", funcName, bookIdArr)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterMappingForSpiderKxmanhuaByHTML,
		"rouman8":  ChapterMappingForSpiderRouman8ByHTML,
	}
	mapping := mappingFactory[websiteName]

	// 5.1. 爬取 chapter
	// -- 请求html页面
	manyBookChapterArrMap, err := GetManyBookAllChapterByCollyMappingV1_5_V1_OnlyForKxmanhua[models.ChapterSpider](mapping.(map[string]models.ModelHtmlMapping), websiteName, bookIdArr)
	// manyBookChapterArrMap, err := GetManyBookAllChapterByCollyMappingV1_5_V2_Common_ForAllWebsite[models.ChapterSpider](mapping.(map[string]models.ModelHtmlMapping), websiteName, bookIdArr)
	chapterNamePreviewCount = 0 // ！！！！重要,必有，重置计数器。chapter中 name包含"Preview"次数
	// -- 插入前数据校验
	if manyBookChapterArrMap == nil || err != nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		return 0, err // 直接结束
	}

	// 5.2. 执行核心逻辑 - 插入部分
	if okTotal, funcErr = SpiderManyBookAllChapter_UpsertPart_V1_OnlyForKxmanhua(websiteName, manyBookChapterArrMap); funcErr != nil {
		log.Errorf("爬取失败, reaason: 插入db失败. website=%v, bookIdArr=%v", websiteName, bookIdArr)
		return 0, funcErr
	}

	// 步骤5.3：更新book 字段：spider_sub_chapter_end_status
	funcErr = db.DBUpdateBatchByIdArr[models.ComicSpider](db.DBComic, bookIdArr, map[string]any{"spider_sub_chapter_end_status": 1})
	if funcErr != nil {
		log.Errorf("func= %v 失败, 更新db book spider_sub_chapter_end_status 状态失败, err: %v", funcName, funcErr)
		return 0, funcErr
	}

	// 6. 返回结果
	log.Info("爬取成功,插入" + strconv.Itoa(okTotal) + "条chapter数据")
	return okTotal, nil
}

// 获取需要爬取的任何id数组，比如 bookIds
func DBGetIdsNeedCrawlByFiled[T any](db *gorm.DB, ids []int, field string, value any) ([]int, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var need []int
	err := db.Model(new(T)).
		Where("id IN ?", ids).
		Where(field+" = ?", value). // Where("spider_end != ?", 0)
		Pluck("id", &need).Error
	if err != nil {
		return nil, err
	}

	return need, nil
}

// 获取需要爬取的任何id数组，比如 bookIds
/*
dbConn db连接对象 -> grom.DB
idsNoFilter 未过滤前 的数组
needCrawlValue 需要爬 对应的值 ，比如 where spider_end = 0, 0 就是 needCrawlValue

返回
[]int 过滤后，需要爬的数组
*/
func DBGetIdsNeedCrawl[T any](dbConn *gorm.DB, idsNoFilter []int, field string, needCrawlValue any) ([]int, error) {
	if len(idsNoFilter) == 0 {
		return nil, nil
	}

	// 获取
	var idsNeedCrawl []int // 需要爬取的数组
	err := db.DBPluckV2[models.ChapterSpider](dbConn, "id", &idsNeedCrawl, db.WithWhere("id IN ?", idsNoFilter), db.WithWhere(field+" = ?", needCrawlValue))
	if err != nil {
		return nil, err
	}

	return idsNeedCrawl, nil
}

// 爬取manychapter allContent V1实现。把爬取+插入放到1个方法，且和 gin.context 解耦
/*
返回：
	okTotal
	error
*/
func SpiderManyChapterAllContent2DB(websiteId int, websiteName string, chapterIdArr []int) (int, error) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "SpiderManyChapterAllContent2DB"
	var funcErr error

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	log.Infof("func=%v, 要爬的chapterIds = %v", funcName, chapterIdArr)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterContentMappingForSpiderKxmanhuaByHTML,
	}
	mapping := mappingFactory[websiteName]

	// 5.1. 爬取 chapter
	// -- 请求html页面
	// manyChapterContentArrMap, err := GetManyChapterAllContentByCollyMappingV1_5_V1[models.ChapterContentSpider](mapping.(map[string]models.ModelHtmlMapping), websiteId, websiteName, chapterIdArr)
	manyChapterContentArrMap, err := GetManyChapterAllContentByCollyMappingV1_5_V2_Async[models.ChapterContentSpider](mapping.(map[string]models.ModelHtmlMapping), websiteId, websiteName, chapterIdArr)
	chapterNamePreviewCount = 0 // ！！！！重要,必有，重置计数器。chapter中 name包含"Preview"次数
	// -- 插入前数据校验
	if err != nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		return 0, err // 直接结束
	}
	if len(manyChapterContentArrMap) == 0 {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空,没爬到, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		return 0, errors.New("manyChapterContentArrMap 为空, 没爬到, 拒绝进入下一步") // 直接结束
	}
	// 5.1.1 插入前必要数据清洗，比如parentId num subNUm要加上
	for chapterId, oneChapterContentArr := range manyChapterContentArrMap {
		for i := range oneChapterContentArr {
			oneChapterContentArr[i].ParentId = chapterId
			oneChapterContentArr[i].Num = i + 1
			oneChapterContentArr[i].SubNum = 0
			oneChapterContentArr[i].RealSortNum = oneChapterContentArr[i].Num
		}
	}

	// 5.2. 执行核心逻辑 - 插入部分
	if okTotal, funcErr = SpiderManyChapterAllContent_UpsertPart(websiteId, websiteName, manyChapterContentArrMap); funcErr != nil {
		log.Errorf("爬取失败, reaason: 插入db失败. website=%v, chapterIdArr=%v", websiteName, chapterIdArr)
		return 0, funcErr
	}

	// 步骤5.3：更新 chapter 字段：spider_end_status (只有爬取个数>0的，才需要更新)
	// 步骤5.3.1 找需要更新的数组
	needUpdateChapterIdArr := []int{}   // 爬取个数>0的
	noNeedUpdateChapterIdArr := []int{} // 不需要更新数组, 爬取个数=0的
	for chapterId, oneChapterArr := range manyChapterContentArrMap {
		if len(oneChapterArr) > 0 {
			needUpdateChapterIdArr = append(needUpdateChapterIdArr, chapterId)
		} else {
			noNeedUpdateChapterIdArr = append(noNeedUpdateChapterIdArr, chapterId)
		}
	}
	if len(noNeedUpdateChapterIdArr) > 0 {
		log.Warn("爬取结果为0的 chapterIdArr = ", noNeedUpdateChapterIdArr)
	}

	// 步骤5.3.2 更新父表
	funcErr = db.DBUpdateBatchByIdArr[models.ComicSpider](db.DBComic, needUpdateChapterIdArr, map[string]any{"spider_end_status": 1})
	if funcErr != nil {
		log.Errorf("func= %v 失败, 更新db book spider_end_status 状态失败, err: %v", funcName, funcErr)
		return 0, funcErr
	}

	// 6. 返回结果
	log.Info("爬取成功,插入" + strconv.Itoa(okTotal) + "条chapter数据")
	return okTotal, nil
}

// 获取多个chapter所有content, 用colly, 通过mapping。串行方式，实现
/*

5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
	步骤1: 找到目标网站
	步骤2: 爬取
	步骤3: 提取数据 <- 往上是本方法
	步骤4: 数据清洗/ 未爬到的字段赋值
	步骤5: 验证爬取数据 准确性
	步骤6: 数据库插入

参数:
	2. mapping map[string]models.ModelMapping 爬取映射关系
	2. websiteName string 网站名称
	3. bookIdArr

返回:
map id -> 数组

主表数组
作用简单说：
*/
func GetManyChapterAllContentByCollyMappingV1_5_V1[T any](mapping map[string]models.ModelHtmlMapping, websiteId int, websiteName string, chapterIdArr []int) (map[int][]T, error) {
	// 初始化
	funcName := "GetManyChapterAllContentByCollyMappingV1_5_V1"
	// chapterId 和 fullUrl 映射关系, key 是url
	chapterIdFullUrlMapKeyUrl := make(map[string]int)

	// 步骤2: 爬取

	// 2. 爬虫相关
	// -- 建一个爬虫对象
	c := colly.NewCollector(
	// colly.Async(true), // ← 这一行没加就一直是串行的
	)

	// -- 设置并发数，和爬取限制
	// 设置请求限制（每秒最多3个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: 3, // 和queue队列同时存在时，用queue控制并发就行。加这个有用，但没必要。默认是0，表示没限制
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// 步骤3: 提取数据
	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的 --
	var oneChapterContentArr []T                    // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider{}
	var manyChapterContenArrMap = make(map[int][]T) //所有 chapter 的 content 数组 map 。key=chapterId value=arr
	// 遍历一个book, 每个chapter
	StagesCfg := config.CfgSpiderYaml.Websites[websiteName].Stages["one_chapter_all_content"]
	everyChapterSelectStr := StagesCfg.Crawl.Selectors["item"].(string) // 每个chapter 选择器
	c.OnHTML(everyChapterSelectStr, func(e *colly.HTMLElement) {
		// 0 初始化
		currentUrl := e.Request.URL.String() // 获取当前正在爬取的 URL

		// 0. 处理异常内容

		// 1. 获取能获取到的
		log.Debug("匹配 oneChapterStr = ", e.Text)
		// -- 创建对象comic
		var chapterContentT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		// log.Infof(" chapterId=%v,当前爬取url=%v 通过mapping规则,爬取结果 result = %v", chapterIdArr, currentUrl, result)
		log.Infof("当前爬取url=%v 通过mapping规则,爬取结果 result = %v", currentUrl, result)
		if result != nil {
			// 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象
			MapByTag(result, &chapterContentT)
			log.Debugf("映射后的 chapter 对象, 还未清洗: %+v", chapterContentT)
		}

		// 数据清洗/校验，如果url 是空的，不处理 --
		// 方式1：用类型断言
		// 数据清洗/校验，如果urlApiPath 是空的，不处理
		/*
			if chapterContent, ok := any(chapterContentT).(models.ChapterContentSpider); ok {
				if chapterContent.UrlApiPath == "" {
					log.Warnf("当前爬取url=%v 的 urlApiPath 为空，跳过该记录", currentUrl)
					return
				}
			}*/
		// 方式2：直接从 result map 中判断
		if urlApiPath, ok := result["urlApiPath"].(string); ok && urlApiPath == "" {
			return
		}

		// 2. 放到 oneChapterContentArr 里
		oneChapterContentArr = append(oneChapterContentArr, any(chapterContentT).(T))
	})

	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		if r == nil {
			// 网络层错误（DNS / timeout / TLS）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 网络层错误（DNS / timeout / TLS）, err: ", err)
			return
		}

		switch {
		case r.StatusCode >= 400 && r.StatusCode < 500:
			// 4xx：客户端错误（参数错误、被封、资源不存在）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 客户端错误（参数错误、被封、资源不存在）, err: ", err)

		case r.StatusCode >= 500 && r.StatusCode < 600:
			// 5xx：服务端错误（可重试）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 服务端错误（可重试）, err: ", err)

		default:
			// 其他非常规状态码
			// 可选重试
		}
	})

	// 成功爬完1页，回调
	c.OnScraped(func(r *colly.Response) {
		// 把 oneBookAllChapter 加到 要返回的map 中去
		currentUrl := r.Request.URL.String()               // 获取当前正在爬取的 URL
		chapterId := chapterIdFullUrlMapKeyUrl[currentUrl] // 通过url 获取到chapterId
		manyChapterContenArrMap[chapterId] = oneChapterContentArr

		oneChapterContentArr = nil // 或者 oneChapterContentArr = []T{} 重置切片
	})

	// -- 添加多个页面到队列中
	// 使用队列控制任务调度（最多并发3个Url）
	q, _ := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})

	// -- 添加任务到队列
	for _, chapterId := range chapterIdArr {
		// 通过websiteId 查询，并拼接出 book的 url --
		// 步骤1: 找到目标网站
		website, err := db.DBFindOneByFieldV1_5[models.Website](db.DBComic, "id", websiteId)
		if err != nil || website == nil {
			log.Errorf("func= %v, 通过websiteId查询website信息失败,因为通过websiteId获取website对象失败, err: %v", funcName, err)
			return nil, errors.New("通过websiteId查询website信息失败")
		}

		chapter, err := db.DBFindOneByFieldV1_5[models.ChapterSpider](db.DBComic, "id", chapterId) // 通过 chapterId 获取 chapter信息
		if err != nil {
			log.Errorf("func= %v, 通过chapterId查询chapter信息失败,因为通过chapterId获取chapter对象失败, err: %v", funcName, err)
			return nil, err
		}
		fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, chapter.UrlApiPath, nil) // 完整爬取 url
		// apiUrlPath := fmt.Sprintf("/test/kxmanhua/spiderChapterContent/%d.html", chapterId)
		// fullUrl := GetSpiderFullUrl(false, "localhost:8080", apiUrlPath, nil) // 完整爬取 url，本地测试
		log.Info("生成的 chapter 爬取 fullURl = ", fullUrl)
		chapterIdFullUrlMapKeyUrl[fullUrl] = chapterId // 记录chapterId 和 fullUrl 的映射关系
		// 打算使用 GET 请求校验 URL 可达性，通过后才加入抓取队列。爬取的一般都是get请求， 就用get请求下。但实际不用 用c.OnError() 就能有类似效果

		// 再添加到队列 --
		q.AddURL(fullUrl) // 只能添加1个

		// 测试用 - 添加任务到队列
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/社团学姐.html") // 章节url
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/1.html") // 章节url
	}
	log.Infof("func=%v, chapterId 和 请求fullUrl映射关系 chapterIdFullUrlMapKeyUrl = %v", funcName, chapterIdFullUrlMapKeyUrl)

	// 启动对垒
	q.Run(c)
	return manyChapterContenArrMap, nil
}

// 获取多个chapter所有content, 用colly, 通过mapping。并行方式，实现
/*

5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
	步骤1: 找到目标网站
	步骤2: 爬取
	步骤3: 提取数据 <- 往上是本方法
	步骤4: 数据清洗/ 未爬到的字段赋值
	步骤5: 验证爬取数据 准确性
	步骤6: 数据库插入

参数:
	2. mapping map[string]models.ModelMapping 爬取映射关系
	2. websiteName string 网站名称
	3. bookIdArr

返回:
map id -> 数组

主表数组
作用简单说：
*/
func GetManyChapterAllContentByCollyMappingV1_5_V2_Async[T any](mapping map[string]models.ModelHtmlMapping, websiteId int, websiteName string, chapterIdArr []int) (map[int][]T, error) {
	// 初始化
	funcName := "GetManyChapterAllContentByCollyMappingV1_5_V2_Async"

	// 步骤2: 爬取

	// 2. 爬虫相关
	// -- 建一个爬虫对象
	c := colly.NewCollector(
		colly.Async(true), // ← 这一行没加就一直是串行的
	)

	// -- 设置并发数，和爬取限制
	// 设置请求限制（每秒最多3个请求, 5秒后发）
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,                         // 和queue队列同时存在时，这个必须有。和queue无关，它是真正控制并发的！！！！！！！
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// ↓↓↓ 新增：保护最终结果 map 的并发安全
	var mu sync.Mutex

	// 步骤3: 提取数据
	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的 --
	// var oneChapterContentArr []T                    // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider...
	// ↓↓↓ 注释掉全局 oneChapterContentArr，因为并发时会串数据
	// var oneChapterContentArr []T
	var manyChapterContenArrMap = make(map[int][]T) //所有 chapter 的 content 数组 map 。key=chapterId value=arr

	// ↓↓↓ 新增：在 OnRequest 阶段为每个请求绑定 chapterId 和 独立的 items slice 指针
	c.OnRequest(func(r *colly.Request) {
		chapterIdStr := r.Ctx.Get("chapter_id")
		if chapterIdStr == "" {
			return
		}

		// 不要 new 一个全新的 context
		// 直接在现有的 r.Ctx 上放
		items := &[]T{}
		r.Ctx.Put("items", items)

		/* 错误写法
		chapterIdStr := r.Ctx.Get("chapter_id")
		if chapterIdStr == "" {
			return
		}

		// 为这个请求创建独立的上下文和 slice
		newCtx := colly.NewContext()
		newCtx.Put("chapter_id", chapterIdStr)
		newCtx.Put("items", &[]T{}) // 每个请求有自己独立的 slice
		r.Ctx = newCtx              // 覆盖为新的 context
		*/
	})

	// 遍历一个book, 每个chapter
	StagesCfg := config.CfgSpiderYaml.Websites[websiteName].Stages["one_chapter_all_content"]
	everyChapterSelectStr := StagesCfg.Crawl.Selectors["item"].(string) // 每个chapter 选择器
	c.OnHTML(everyChapterSelectStr, func(e *colly.HTMLElement) {
		// 0 初始化
		// ↓↓↓ 新增：从上下文取 chapterId 和 items slice
		ctx := e.Request.Ctx
		chapterIdStr := ctx.Get("chapter_id")
		if chapterIdStr == "" {
			return
		}

		chapterId, err := strconv.Atoi(chapterIdStr)
		if err != nil {
			return
		}

		itemsPtrAny := ctx.GetAny("items")
		if itemsPtrAny == nil {
			return
		}
		itemsPtr, ok := itemsPtrAny.(*[]T)
		if !ok || itemsPtr == nil {
			return
		}

		// 1. 获取能获取到的
		log.Debug("匹配 oneChapterStr = ", e.Text)
		// -- 创建对象comic
		var chapterContentT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		log.Infof("当前爬取 chapterId = %v url=%v 通过mapping规则,爬取结果 result = %v", chapterId, e.Request.URL.String(), result)
		if result != nil {
			// 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象
			MapByTag(result, &chapterContentT)
			log.Debugf("映射后的 chapter 对象, 还未清洗: %+v", chapterContentT)
		}

		// 数据清洗/校验，如果url 是空的，不处理 --
		// 方式1：用类型断言
		// 数据清洗/校验，如果urlApiPath 是空的，不处理
		/*
			if chapterContent, ok := any(chapterContentT).(models.ChapterContentSpider); ok {
				if chapterContent.UrlApiPath == "" {
					log.Warnf("当前爬取url=%v 的 urlApiPath 为空，跳过该记录", currentUrl)
					return
				}
			}*/
		// 方式2：直接从 result map 中判断
		if urlApiPath, ok := result["urlApiPath"].(string); ok && urlApiPath == "" {
			return
		}

		// 2. 放到 oneChapterContentArr 里
		// oneChapterContentArr = append(oneChapterContentArr, any(chapterContentT).(T))
		// ↓↓↓ 修改为：追加到当前请求专属的 slice
		*itemsPtr = append(*itemsPtr, chapterContentT)
	})

	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		if r == nil {
			// 网络层错误（DNS / timeout / TLS）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 网络层错误（DNS / timeout / TLS）, err: ", err)
			return
		}
		// 下面好像触发不到！！！
		switch {
		case r.StatusCode >= 400 && r.StatusCode < 500:
			// 4xx：客户端错误（参数错误、被封、资源不存在）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 客户端错误（参数错误、被封、资源不存在）, err: ", err)

		case r.StatusCode >= 500 && r.StatusCode < 600:
			// 5xx：服务端错误（可重试）
			log.Error("func= GetOneBookAllChapterByCollyMappingV1_5, 服务端错误（可重试）, err: ", err)

		default:
			// 其他非常规状态码
			// 可选重试
		}
	})

	c.OnScraped(func(r *colly.Response) {
		// ↓↓↓ 修改为：从上下文取 chapterId
		ctx := r.Ctx
		chapterIdStr := ctx.Get("chapter_id")
		if chapterIdStr == "" {
			return
		}

		chapterId, err := strconv.Atoi(chapterIdStr)
		if err != nil {
			return
		}

		// ↓↓↓ 新增：从上下文取 items 并存入最终 map
		itemsPtrAny := ctx.GetAny("items")
		if itemsPtrAny == nil {
			return
		}
		itemsPtr, ok := itemsPtrAny.(*[]T)
		if !ok || itemsPtr == nil {
			return
		}

		mu.Lock()
		// manyChapterContenArrMap[chapterId] = oneChapterContentArr  // 原来写法
		// oneChapterContentArr = nil // 或者 oneChapterContentArr = []T{} 重置切片  // 原来写法
		manyChapterContenArrMap[chapterId] = *itemsPtr // 复制内容
		mu.Unlock()

	})
	// -- 添加多个页面到队列中
	// 使用队列控制任务调度（最多并发3个Url）
	q, err := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})
	if err != nil {
		return nil, err
	}

	// -- 添加任务到队列
	for _, chapterId := range chapterIdArr {
		// 通过websiteId 查询，并拼接出 book的 url --
		// 步骤1: 找到目标网站
		website, err := db.DBFindOneByFieldV1_5[models.Website](db.DBComic, "id", websiteId)
		if err != nil || website == nil {
			log.Errorf("func= %v, 通过websiteId查询website信息失败,因为通过websiteId获取website对象失败, err: %v", funcName, err)
			return nil, errors.New("通过websiteId查询website信息失败")
		}

		chapter, err := db.DBFindOneByFieldV1_5[models.ChapterSpider](db.DBComic, "id", chapterId) // 通过 chapterId 获取 chapter信息
		if err != nil {
			log.Errorf("func= %v, 通过chapterId查询chapter信息失败,因为通过chapterId获取chapter对象失败, err: %v", funcName, err)
			return nil, err
		}
		fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, chapter.UrlApiPath, nil) // 完整爬取 url
		// apiUrlPath := fmt.Sprintf("/test/kxmanhua/spiderChapterContent/%d.html", chapterId)
		// fullUrl := GetSpiderFullUrl(false, "localhost:8080", apiUrlPath, nil) // 完整爬取 url，本地测试
		log.Info("生成的 chapter 爬取 fullURl = ", fullUrl)

		// 打算使用 GET 请求校验 URL 可达性，通过后才加入抓取队列。爬取的一般都是get请求， 就用get请求下。但实际不用 用c.OnError() 就能有类似效果
		ctx := colly.NewContext()
		ctx.Put("chapter_id", strconv.Itoa(chapterId))

		// 关键：使用 c.Request() 创建 request（自动处理 URL 解析）
		// 正确调用 c.Request()，它只返回 *colly.Request 和 error
		err = c.Request("GET", fullUrl, nil, ctx, nil) // 用c.Request 就不用写q.AddRequest(req) 这种了？自动就进去queue了
		if err != nil {
			log.Warnf("创建 request 失败 %s: %v", fullUrl, err)
			continue
		}

		// q.AddRequest(req)  // 用不到.Request 就表示 自动就 入队列了
	}

	// 在 q.Run(c) 前加
	// 启动队列并等待完成
	q.Run(c)

	// ↓↓↓ 新增：等待所有异步请求完成（虽然 q.Run 通常会阻塞，但加 c.Wait() 更保险）
	c.Wait()

	return manyChapterContenArrMap, nil
}

// 把爬取 manyChapterAllCcontent 分成2部分。爬取部分 + 插入部分
// 插入部分 - 插入多个章节的 内容
/*
5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值 <- 本方法
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
			- 6.1 插入 章节
			- 6.2 更新 book 表stats字段
参数:
manyBookChapterArrMap []

返回 插入成功总数
 error
*/
func SpiderManyChapterAllContent_UpsertPart(websiteId int, websiteName string, manyChapterContentArrMap map[int][]models.ChapterContentSpider) (int, error) {
	// 初始化
	okTotal := 0 // 插入成功总数
	funcName := "SpiderManyChapterAllContent_UpsertPart"

	// 异常处理
	if len(manyChapterContentArrMap) == 0 {
		return 0, errors.New("func=爬chapter AllContent, manyChapterContentArrMap 为空")
	}

	for chapterId, oneChapterContentArr := range manyChapterContentArrMap {
		if len(oneChapterContentArr) == 0 { // 如果是空数组，跳过，防止后面返回错误err
			log.Warn("chapterId = ", chapterId, " 爬取结果为0，跳过该记录")
			continue
		}
		// 步骤4: 数据清洗/ 未爬到的字段赋值
		// -- 赋值上下文参数 + 数据清洗。（赋值上下文参数：是吧方法传参，给对象赋值。数据清洗：设置-爬取字段，或者默认数据）
		for i := range oneChapterContentArr {
			// -赋值 上下文传参。如parentId (非数据清洗业务，放在这里)
			// oneBookChapterArr[i].ParentId = int(bookIdArr[index]) // delete - 弃用。会导致 comic和chapter内容对不上!!!. 父id，应该从manyBookChapterArr 里拿，这里是最准确的。因为bookIdArr 是从小到大，但爬出来的 manyBookChapterArr 是按id 随机的。容易导致：comic和chapter的 真实章节 对不上
			oneChapterContentArr[i].ParentId = chapterId // 父id，应该从manyChapterContentArr 里拿，这里是最准确的。因为chapterId 是从小到大，但爬出来的 manyChapterContentArr 是按id 随机的。容易导致：comic和chapter的 真实章节 对不上
			// -数据清洗
			oneChapterContentArr[i].DataClean() // 数据清洗
			log.Debug("清洗完数据 chapter = ", oneChapterContentArr[i])
		}

		// -- 检测下爬到的数据，有没有重复数据，需要注意下。只要判单章节号码 num 就可以了
		var spiderChapterNumArr []int // 爬到的章节号 arr
		for _, chapter := range oneChapterContentArr {
			spiderChapterNumArr = append(spiderChapterNumArr, chapter.Num)
		}
		if util.HasDuplicate(spiderChapterNumArr) { // 判断有重复
			log.Warn("爬取1本书 AllChapter, 爬到的章节号码 有重复, 要注意下, chapterId = ", chapterId)
		}

		// 步骤6: 数据库插入
		// 6.1 插入 章节. upsert chapter
		// 获取配置 --
		webCfg := config.CfgSpiderYaml.Websites[websiteName]
		if webCfg == nil {
			return 0, nil
		}
		log.Debug("------- webCfg = ", webCfg)

		// 获取 one_chapter_all_content 阶段配置
		stageCfg := webCfg.Stages["one_chapter_all_content"]
		if stageCfg == nil {
			return 0, errors.New("func=爬chapter AllContent V1.5, stageCfg 为空")
		}
		log.Debug("------- stageCfg = ", stageCfg)

		// 批量插入db chapter
		if len(oneChapterContentArr) == 0 {
			return 0, errors.New("func=爬chapter AllContent V1.5, chapterArr 为空")
		}

		// err := db.DBUpsertBatch(db.DBComic, oneChapterContentArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys) // v1写法，保留，不能指定表名，比如content_005这个分表，因此引入v2写法
		tableName := fmt.Sprintf("chapter_content_spider_%04d", websiteId)
		err := db.DBUpsertBatchV2SpecifyTableName(db.DBComic, tableName, oneChapterContentArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys) // v2写法，可以指定表名，比如content_005这个分表

		if err != nil {
			log.Errorf("func= %v, 批量插入db chapter 失败, err: %v", funcName, err)
		}
		okTotal += len(oneChapterContentArr) // 更新总数

		// 6.2 更新 chapter 表stats字段
		// -- 创建 chapter对象
		var chapterSpider models.ChapterSpider
		chapterSpider.Id = chapterId
		chapterSpider.SpiderEndStatus = 1 // 爬取完成状态, 这是要修改的字段

		err = db.DBUpdate(db.DBComic, chapterSpider, stageCfg.UpdateParentStats.UniqueKeys, stageCfg.UpdateParentStats.UpdateKeys)
		if err != nil {
			log.Errorf("func= %v, 更新 chapter spider_end_status 失败, err: %v", funcName, err)
			return 0, err
		}
	}

	log.Infof("插入成功 %v 条", okTotal)
	return okTotal, nil // 一切正常
}

// 爬取 SpiderOneTypePageAllBook2DB V1实现。把爬取+插入放到1个方法，且和 gin.context 解耦
/*
返回：
	okTotal 成功总数
	error
*/
func SpiderOneTypeAllBook2DBV1(reqDTO SpiderOneTypeAllBookReqV15V1) (int, error) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "SpiderOneTypeAllBook2DBV1"

	// 1. 获取传参。实现方式: 从req中拿

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ComicMappingForSpiderKxmanhuaByHtml,
		"rouman8":  ComicMappingForSpiderRouman8ByHtml,
	}

	// 5.1 爬取
	manyBookArr2D, err := SpiderOneTypeAllBookUseCollyByMappingV2_Sync[models.ComicSpider](mappingFactory, reqDTO)
	if err != nil {
		log.Errorf("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		return 0, err // 直接结束
	}
	if len(manyBookArr2D) == 0 {
		log.Errorf("func=%v,失败, 爬取数据为空", funcName)
		return 0, errs.ErrNull
	}
	// 5.2 插入
	err = SpiderOneTypeAllBook2DBV1_Upsert_Part(reqDTO.WebsiteId, reqDTO.SpiderTag.Website, &manyBookArr2D)
	if err != nil {
		log.Errorf("func=%v, 批量插入db 失败, err: %v", funcName, err)
		return 0, err
	}

	// 6 返回
	// 计算总数
	for _, oneBookArr := range manyBookArr2D {
		okTotal += len(oneBookArr)
	}

	return okTotal, nil // 成功
}

// 爬取某类 所有book V2实现,根据mapping. 同步
/*
常用方法写法：
func SpiderOneTypeAllBookUseCollyByMappingV2_Sync[T any](websiteId int, websiteName string, spiderUrlNoSetValue string, startPageNum, endPageNum int) ([][]T, error) {
常用步骤：
	// 步骤0：初始化
	// 步骤0.5：从参数中再获取参数。有的参数，传的 前端请求reqDTO
	// 步骤1：参数异常判断
	// 步骤2: 爬取
	// 2.1 建一个爬虫对象
	// 2.2 设置请求限制（例如：每秒最3个请求, 每个请求发前随机延迟5秒）
	// 步骤3: 处理回调-colly请求前
	// 步骤4: 提取数据 c.OhHTML() ,之前需要从获取配置
	// 4.1 获取配置
	// 4.2 提取数据

	// 步骤5：处理回调-成功完成时
	// 步骤6：处理回调-发生错误时
	// 步骤7: 创建队列 queue 对象
	// 步骤8: 添加任务到队列
	// 8.1 拼接请求url
	// 8.2 添加请求url到队列
	// 8.2.1 传递自定义参数，判断当前是第N页，url 应该可以之间获取到
	// 8.2.2 添加到队列
	// 步骤9: 启动队列

	T 主表
	SubT 子表
*/
// func SpiderOneTypeAllBookUseCollyByMappingV2_Sync[T any, SubT any](mappingFactory map[string]models.ModelHtmlMapping, reqDTO SpiderOneTypeAllBookReqV15V1) ([][]T, error) {
func SpiderOneTypeAllBookUseCollyByMappingV2_Sync[T any](mappingFactory map[string]any, reqDTO SpiderOneTypeAllBookReqV15V1) ([][]T, error) {
	// 步骤0：初始化
	funcName := "SpiderOneTypeAllBookUseCollyByMappingV2_Sync"

	// 步骤0.5：从参数中再获取参数。有的参数，传的 前端请求reqDTO
	websiteName := reqDTO.SpiderTag.Website // 就是:kxmanhua rouman8 这种，一般不加.com,参考 配置文件：rouman8 kxmanhua 配置
	websiteId := reqDTO.WebsiteId
	pornTypeId := reqDTO.PornTypeId
	countryId := reqDTO.CountryId
	typeId := reqDTO.TypeId
	processId := reqDTO.ProcessId
	authorConcatType := reqDTO.AuthorConcatType
	spiderUrlNoSetValue := reqDTO.SpiderUrl // 未传值的rul，带%d的写法。如： "https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
	startPageNum := reqDTO.StartPageNum
	endPageNum := reqDTO.EndPageNum

	// 步骤1：参数异常判断
	// 1.1 endPageNum 必须 >= startPageNUm
	if endPageNum < startPageNum {
		return nil, errs.ErrInvalidPageNum
	}

	// 步骤2: 爬取
	// 2.1 建一个爬虫对象
	c := colly.NewCollector(
		colly.Async(false), // ← 这一行没加就一直是串行的，或者改成false
	)

	// 2.2 设置请求限制（例如：每秒最3个请求, 每个请求发前随机延迟5秒）
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,                         // 和queue队列同时存在时，这个必须有。和queue无关，它是真正控制并发的！！！！！！！ 这里非并发，所有注释了
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// 步骤3: 处理回调-colly请求前
	// 步骤4: 提取数据
	// 4.1 获取配置
	StagesCfg := config.CfgSpiderYaml.Websites[websiteName].Stages["one_type_all_book"]
	bookArrCssSelector := StagesCfg.Crawl.Selectors["arr"].(string)      // 某一页 all book选择器
	bookArrItemCssSelector := StagesCfg.Crawl.Selectors["item"].(string) // 每个book 选择器

	// 4.2 提取数据
	var allPageBookArr2D [][]T // 存放爬好的 obj, 二维数组
	var mu sync.Mutex          // 添加互斥锁
	mapping := mappingFactory[websiteName].(map[string]models.ModelHtmlMapping)
	// 遍历每一个 bookArr .c.OnHTML() 根据 CSS选择器, 就让触发1次
	c.OnHTML(bookArrCssSelector, func(eBookArr *colly.HTMLElement) {
		// 遍历每一个 bookArrItem, 用forEach. colly，用Html遍历
		var onePageBookArr []T
		eBookArr.ForEach(bookArrItemCssSelector, func(i int, e *colly.HTMLElement) {
			// 1. 获取能获取到的
			// 通过mapping -> 转成1个对象
			// 创建对象comic
			var comicT T
			// var comicTStats SubT //   comicSpiderStats := models.ComicSpiderStats{} // 子表，统计数据。不知道用不用得着

			// 通过mapping 爬内容
			result := GetOneObjByCollyMapping(e, mapping)
			if result != nil {
				MapByTag(result, &comicT) // 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象
				log.Infof("映射后的comic对象-只有爬取到的数据: %+v", comicT)
			}

			// 2. 设置对象值
			// -- T 类型 -》 具体struct 类型
			comic := any(comicT).(models.ComicSpider)

			// -- 进度id逻辑
			if processId == 1 {
				comic.ProcessId = comic.End // 如果用户传 1(待分类) - 》程序自己判断
			} else {
				comic.ProcessId = int(processId) // 如果是2/3, 就直接替换赋值，以前端传参为主
			}
			// -- 其它直接赋值
			comic.WebsiteId = int(websiteId)               // 网站id
			comic.PornTypeId = int(pornTypeId)             // 色情类型id
			comic.CountryId = int(countryId)               // 国家id
			comic.TypeId = int(typeId)                     // 类型id
			comic.AuthorConcatType = int(authorConcatType) // 作者拼接方式 id

			// 3 数据清洗
			comic.DataClean()

			// 3.5 处理子表
			// 如果子表爬到数据了，才处理 没实现--
			comicStats := comic.Stats
			// 修复指针类型比较错误：LatestChapterId是*int类型，需要先检查是否为nil再解引用比较
			latestChapterIdValid := comicStats.LatestChapterId != nil && *comicStats.LatestChapterId > 0
			if latestChapterIdValid || comicStats.Star > 0 || comicStats.LatestChapterName != "" || comicStats.TotalChapter > 0 || comicStats.Hits > 0 || comicStats.LastestChapterReleaseDate.After(time.Date(1001, 1, 1, 8, 0, 0, 0, time.UTC)) { // 程序爬不到发布时间就 1001-01-01 08:00:00
				log.Warn("还没实现, 子表爬到数据，要插入有用数据，不然全是0！！")
			}

			// 4 把爬好的单个数据，放到数组里，准备插入数据库
			onePageBookArr = append(onePageBookArr, any(comic).(T))
		})

		// 3. 遍历完之后，加到 allPageBookArr 里 - 使用互斥锁保护 (因为要操作多个线程 用的共享对象 -》 allPageBookArr)
		mu.Lock()
		allPageBookArr2D = append(allPageBookArr2D, onePageBookArr)
		mu.Unlock()
	})

	// 步骤5：处理回调-成功完成时
	c.OnScraped(func(r *colly.Response) {
		// 暂时啥也没写
	})

	// 步骤6：处理回调-发生错误时
	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		if r == nil {
			// 网络层错误（DNS / timeout / TLS）
			log.Errorf("func= %v, 网络层错误（DNS / timeout / TLS）, err= %v ", funcName, err)
			return
		}

		switch {
		case r.StatusCode >= 400 && r.StatusCode < 500:
			// 4xx：客户端错误（参数错误、被封、资源不存在）
			log.Errorf("func= %v, 客户端错误（参数错误、被封、资源不存在）, err= %v ", funcName, err)

		case r.StatusCode >= 500 && r.StatusCode < 600:
			// 5xx：服务端错误（可重试）
			log.Errorf("func= %v, 服务端错误（可重试）, err=%v ", funcName, err)

		default:
			// 其他非常规状态码
			// 可选重试
		}
	})

	// 步骤7: 创建队列 queue 对象
	q, err := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum, &queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum})
	if err != nil {
		return nil, err
	}

	// 步骤8: 添加任务到队列
	// 8.1 拼接请求url
	spiderUrlArr := make([]string, endPageNum-startPageNum+1)
	for i := range spiderUrlArr {
		spiderUrlArr[i] = fmt.Sprintf(spiderUrlNoSetValue, startPageNum+i)
		log.Infof("爬取请求url, 第%v个, url= =%v", i+1, spiderUrlArr[i])

		// 8.2 添加请求url到队列
		// 8.2.1 传递自定义参数，判断当前是第N页，url 应该可以之间获取到
		ctx := colly.NewContext()
		ctx.Put("currentPageNum", strconv.Itoa(i+1))

		// 8.2.2 添加到队列
		c.Request("GET", spiderUrlArr[i], nil, nil, nil)
	}

	// 步骤9: 启动队列
	q.Run(c)

	return allPageBookArr2D, nil
}

// SpiderOneTypeAllBook2DBV1_Upsert_Part  插入部分
// 插入部分 - 插入多个章节的 内容
/*
5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值 <- 本方法
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
			- 6.1 插入 章节
			- 6.2 更新 book 表stats字段

// 总结出插入方法，经典流程 comic - 简单
	// 步骤0：初始化
	// 步骤1：参数异常判断
	// 步骤2：参数异常业务逻辑判断。比如判断websiteId是否存在
	// 步骤3：真实业务，插入comic表
	// 步骤4：真实业务，插入stats 关联表
	// 步骤5：返回

// 总结出插入方法，经典流程 comic - 详细
	// 步骤0：初始化
		// 0.1 根据网站名字，获取配置文件中 相关配置
	// 步骤1：参数异常判断
		// 1.1 websiteId 必须>0
	// 步骤2：参数异常业务逻辑判断。比如判断websiteId是否存在
		// 2.1 websiteId必须 能从数据库 查到
	// 步骤3：真实业务，插入comic表 （要想插入的comic对象回填id, allPageComicArr 必须传指针，for allPageComicArr 里也必须传指针）
		// 3.1 从二维数组中，取出每一页,爬的数据，插入数据库
		// 3.2 插入主表
		// 3.3 判断是否回填ID成功
	// 步骤4：真实业务，插入stats 关联表
		// 4.1 准备插入数据
			// 4.1.1 新建插入数组
			// 4.1.2 从comic取有用数据
					// 取的数据 异常判断
							- 取的ID是否无效
					// 取有用数据
			// 4.1.3 通用数据清洗
		// 4.2 插入数据, 异常判断
			// 4.2.1 父表comic 和stats关联表插入数据不一致
		// 4.3 插入
	// 步骤5：返回

参数:
manyBookChapterArrMap []

返回 插入成功总数
 error
 不用返回 成功总数了，因为 len(allPageBookArr) 就是, -> 因为 stats 个数，一定= comic个数
*/
func SpiderOneTypeAllBook2DBV1_Upsert_Part(websiteId int, websiteName string, allPageBookArr2D *[][]models.ComicSpider) error {
	// 步骤0：初始化

	// 0.1 根据网站名字，获取配置文件中 相关配置
	webCfg := config.CfgSpiderYaml.Websites[websiteName]
	if webCfg == nil {
		log.Error("未根据websiteName, 找到配置文件中 配置webCfg")
		return errs.ErrNoGetConfig
	}
	log.Debug("------- webCfg = ", webCfg)

	// 获取 one_type_all_book 阶段配置
	stageCfg := webCfg.Stages["one_type_all_book"]
	if stageCfg == nil {
		log.Error("未根据websiteName, 找到配置文件中 阶段配置 stageCfg")
		return errs.ErrNoGetConfig
	}
	log.Debug("---------- 配置文件, stageCfg.Insert.UniqueKeys = ", stageCfg.Insert.UniqueKeys)
	log.Debug("---------- 配置文件，stageCfg.Insert.UpdateKeys = ", stageCfg.Insert.UpdateKeys)

	// 步骤1：参数异常判断
	// 1.1 websiteId 必须>0
	if websiteId <= 0 {
		log.Error("参数错误, websiteId必须大于0")
		return errs.GetErr("参数错误, websiteId必须大于0")
	}

	// 步骤2：参数异常业务逻辑判断。比如判断websiteId是否存在
	// 2.1 websiteId必须 能从数据库 查到
	_, err := db.DBFindOneByField[models.Website]("id", websiteId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)
			return errs.GetErr(fmt.Sprintf("websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId))
		}
		return errs.GetErr("查询website失败")
	}

	// 步骤3：真实业务，插入comic表
	// 3.1 从二维数组中，取出每一页,爬的数据，插入数据库 --
	for pageIdx, onePageBookArr := range *allPageBookArr2D {
		log.Infof("========== 开始处理第 %d 页，共 %d 条数据 ==========", pageIdx+1, len(onePageBookArr))

		// 3.2 插入主表
		err := db.DBUpsertBatch(db.DBComic, onePageBookArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)
		if err != nil {
			log.Errorf("第%d页: 批量插入db-comic失败, err = %v", pageIdx+1, err)
			return errs.GetErr(fmt.Sprintf("第%d页批量插入db-comic失败: %v", pageIdx+1, err))
		}
		log.Infof("第%d页: 主表 comic_spider 插入成功，共 %d 条", pageIdx+1, len(onePageBookArr))

		// 3.3 判断是否回填ID成功 (mysql永远无法成功，因为涉及update的操作，不会回填id。因此只能插入前，自己手动查询并回填下正确id)
		if onePageBookArr[0].Id <= 0 { // gorm mysql upsert 总会回填一个非0Id，只是update的时候，回填的错的。因此，这里永远不会触发
			log.Error("Id 回填失败, 跳过本次循环")
			continue
		}

		// -- 重要：由于GORM批量Upsert时不会更新对象的ID字段，需要重新查询获取正确的ID
		for comicIdx := range onePageBookArr {
			comic := &onePageBookArr[comicIdx]
			// 使用对象作为查询条件
			existingComic, err := db.DBFindOneByUniqueIndexMapCondition(comic, stageCfg.Insert.UniqueKeys)
			if err == nil {
				// 更新对象的ID为数据库中的实际ID
				oldId := comic.Id
				comic.Id = existingComic.Id
				log.Debugf("第%d页第%d条: 更新comic ID成功,  旧ID=%d -> 新ID=%d, 名称=%s", pageIdx+1, comicIdx+1, oldId, comic.Id, comic.Name)
			}
		}

		// 步骤4：真实业务，插入stats 关联表
		// 4.1 准备插入数据
		// 4.1.1 新建插入数组
		var comicStatsArr []models.ComicSpiderStats

		// 4.1.2 从comic取有用数据
		statsBuildFailedCount := 0
		for i, comic := range onePageBookArr {
			// 取的数据 异常判断 --
			// 取的ID是否无效
			if comic.Id <= 0 {
				statsBuildFailedCount++
				log.Errorf("第%d页第%d条: 构建stats失败，comic ID无效 (ID=%d), comic名称=%s", pageIdx+1, i+1, comic.Id, comic.Name)
				continue
			}

			// 取有用数据
			stats := models.ComicSpiderStats{
				ComicId:                   comic.Id, // 现在使用正确的ID
				Star:                      comic.Stats.Star,
				LatestChapterName:         comic.Stats.LatestChapterName, // 最新章节名字
				Hits:                      comic.Stats.Hits,
				TotalChapter:              comic.Stats.TotalChapter,
				LastestChapterReleaseDate: comic.Stats.LastestChapterReleaseDate,
			}

			// 4.1.3 通用数据清洗
			stats.DataClean()

			comicStatsArr = append(comicStatsArr, stats)
			log.Debugf("第%d页第%d条: 构建stats成功, ComicId=%d, Star=%.2f, Hits=%d", pageIdx+1, i+1, stats.ComicId, stats.Star, stats.Hits)
		}

		// 4.2 插入数据, 异常判断
		// 4.2.1 父表comic 和stats关联表插入数据不一致
		if len(comicStatsArr) != len(onePageBookArr) {
			log.Errorf("第%d页: 数据不一致！主表插入 %d 条，但stats只构建了 %d 条", pageIdx+1, len(onePageBookArr), len(comicStatsArr))
			return errs.GetErr(fmt.Sprintf("第%d页: 数据不一致！主表 %d 条, 但stats只构建了 %d 条", pageIdx+1, len(onePageBookArr), len(comicStatsArr)))
		}

		// 4.3 插入
		err = db.DBUpsertBatch(db.DBComic, comicStatsArr, stageCfg.RelatedTables["comic_stats"].Insert.UniqueKeys, stageCfg.RelatedTables["comic_stats"].Insert.UpdateKeys)
		if err != nil {
			log.Errorf("第%d页: 批量插入db-comic-stats表失败, err = %v", pageIdx+1, err)
			return errs.GetErr(fmt.Sprintf("第%d页, comic-stats表失败, err = %v", pageIdx+1, err))
		}
		log.Infof("第%d页: 关联表 comic_spider_stats 插入成功，共 %d 条", pageIdx+1, len(comicStatsArr))

		// 打印结果
		onePageOkTotal := len(onePageBookArr) // 每页成功条数
		log.Infof("========== 第%d页处理完成，成功插入 %d 条 book 数据（主表+stats表） ==========", pageIdx+1, onePageOkTotal)
	}

	// 步骤5：返回
	return nil // 一切正常
}

// 爬取manybook allChapter V2实现。适用所有网站，能爬章节时，同时处理 能爬到的表数据。如：comic、comic_stats、authoer
// 把爬取+插入放到1个方法，且和 gin.context 解耦
/*
返回：
	okTotal
	error
*/
func SpiderManyBookAllChapter2DB_V2_Common_CanUpdateOtherTable(websiteName string, bookIdArr []int) (int, error) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "SpiderManyBookAllChapter2DB_V2_Common_CanUpdateOtherTable"
	var funcErr error

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	log.Infof("func=%v, 要爬的bookId = %v", funcName, bookIdArr)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterMappingForSpiderKxmanhuaByHTML,
		// 爬章节依次要传 4个 mapping: 1. 爬章节mapping 2. 通过章节-爬book 3. 通过章节-爬book-stats 4. 通过章节-爬author, 爬不到就传nil
		"rouman8": []any{ChapterMappingForSpiderRouman8ByHTML, ChapterMappingForSpiderRouman8ByHTML_CanGetBook,
			ChapterMappingForSpiderRouman8ByHTML_CanGetBookStats, ChapterMappingForSpiderRouman8ByHTML_CanGetAuthor},
	}
	mappingArr := mappingFactory[websiteName]

	// 5.1. 爬取 chapter
	// -- 请求html页面
	manyBookChapterArrMap, err := GetManyBookAllChapterByCollyMappingV1_5_V1_OnlyForKxmanhua[models.ChapterSpider](mapping.(map[string]models.ModelHtmlMapping), websiteName, bookIdArr)
	manyBookChapterArrMap, err := GetManyBookAllChapterByCollyMappingV1_5_V2_Common_ForAllWebsite[models.ChapterSpider](mappingArr.(map[string]models.ModelHtmlMapping), websiteName, bookIdArr)
	chapterNamePreviewCount = 0 // ！！！！重要,必有，重置计数器。chapter中 name包含"Preview"次数
	// -- 插入前数据校验
	if manyBookChapterArrMap == nil || err != nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		return 0, err // 直接结束
	}

	// 5.2. 执行核心逻辑 - 插入部分
	if okTotal, funcErr = SpiderManyBookAllChapter_UpsertPart_V1_OnlyForKxmanhua(websiteName, manyBookChapterArrMap); funcErr != nil {
		log.Errorf("爬取失败, reaason: 插入db失败. website=%v, bookIdArr=%v", websiteName, bookIdArr)
		return 0, funcErr
	}

	// 步骤5.3：更新book 字段：spider_sub_chapter_end_status
	funcErr = db.DBUpdateBatchByIdArr[models.ComicSpider](db.DBComic, bookIdArr, map[string]any{"spider_sub_chapter_end_status": 1})
	if funcErr != nil {
		log.Errorf("func= %v 失败, 更新db book spider_sub_chapter_end_status 状态失败, err: %v", funcName, funcErr)
		return 0, funcErr
	}

	// 6. 返回结果
	log.Info("爬取成功,插入" + strconv.Itoa(okTotal) + "条chapter数据")
	return okTotal, nil
}

// 获取多个book所有chapter, 用colly, 通过mapping V2,适用于 各种网站
/*

5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
	步骤1: 找到目标网站
	步骤2: 爬取
	步骤3: 提取数据 <- 往上是本方法
	步骤4: 数据清洗/ 未爬到的字段赋值
	步骤5: 验证爬取数据 准确性
	步骤6: 数据库插入

参数:
	2. mapping map[string]models.ModelMapping 爬取映射关系
	2. websiteName string 网站名称
	3. bookIdArr

返回:
map id -> 数组

主表数组
作用简单说：
*/
func GetManyBookAllChapterByCollyMappingV1_5_V2_Common_ForAllWebsite[T any](mapping map[string]models.ModelHtmlMapping, websiteName string, bookIdArr []int) (map[int][]T, error) {
	还需实现
}
