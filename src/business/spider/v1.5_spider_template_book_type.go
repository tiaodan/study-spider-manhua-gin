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
	"strings"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util"
	"study-spider-manhua-gin/src/util/langutil"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
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
		log.Debug("-------------- 匹配 .chapter_list a = ", e.Text)
		// -- 创建对象comic
		var chapterT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		log.Info("------------ 通过mapping规则,爬取结果 result = ", result)
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
	log.Info("------ update comicSpiderStats = ", comicSpiderStats)
	log.Info("------ update comicSpiderStats.ComicId = ", comicSpiderStats.ComicId)
	log.Info("------ update comicSpiderStats.LatestChapterId = ", *comicSpiderStats.LatestChapterId)
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

// 获取多个book所有chapter, 用colly, 通过mapping
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
主表数组
作用简单说：
*/
func GetManyBookAllChapterByCollyMappingV1_5[T any](mapping map[string]models.ModelHtmlMapping, websiteName string, bookIdArr []int) ([][]T, error) {
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
	var oneBookChapterArr []T    // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider{}
	var manyBookChapterArr [][]T //所有book 的chapte数组
	// 遍历一个book, 每个chapter
	// c.OnHTML(".chapter_list a", func(e *colly.HTMLElement) {
	StagesCfg := config.CfgSpiderYaml.Websites[websiteName].Stages["one_book_all_chapter"]
	everyChapterSelectStr := StagesCfg.Crawl.Selectors["item"].(string) // 每个chapter 选择器
	c.OnHTML(everyChapterSelectStr, func(e *colly.HTMLElement) {
		// 0. 处理异常内容
		// -- 处理 ”休刊公告“
		oneChapterStr, _ := langutil.TraditionalToSimplified(e.Text)
		if strings.Contains(oneChapterStr, "休刊") {
			return // ✅ 这个 return 只从匿名函数返回，不会影响 GetOneBookObjByCollyMapping 函数
		}

		// 1. 获取能获取到的
		log.Debug("-------------- 匹配 .chapter_list a = ", e.Text)
		// -- 创建对象comic
		var chapterT T

		// -- 通过mapping 爬内容
		result := GetOneObjByCollyMapping(e, mapping)
		log.Info("------------ 通过mapping规则,爬取结果 result = ", result)
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
		// 把 oneBookAllChapter 加到大数组中去
		manyBookChapterArr = append(manyBookChapterArr, oneBookChapterArr)
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
		// fullUrl := GetSpiderFullUrl(website.IsHttps, website.Domain, book.ComicUrlApiPath, nil) // 完整爬取 url
		apiUrlPath := fmt.Sprintf("/test/kxmanhua/spiderChapter/%d.html", bookId)
		fullUrl := GetSpiderFullUrl(false, "localhost:8080", apiUrlPath, nil) // 完整爬取 url，本地测试
		log.Info("生成的book 爬取 fullURl = ", fullUrl)
		// 打算使用 GET 请求校验 URL 可达性，通过后才加入抓取队列。爬取的一般都是get请求， 就用get请求下。但实际不用 用c.OnError() 就能有类似效果

		// 再添加到队列 --
		q.AddURL(fullUrl) // 只能添加1个

		// 测试用 - 添加任务到队列
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/社团学姐.html") // 章节url
		// q.AddURL("http://localhost:8080/test/kxmanhua/spiderChapter/1.html") // 章节url
	}

	// 启动对垒
	q.Run(c)
	return manyBookChapterArr, nil
}

// 把爬取 manyBookAllChapter 分成2部分。爬取部分 + 插入部分
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

返回 插入成功总数
 error
*/
func SpiderManyBookAllChapter_UpsertPart(websiteName string, bookIdArr []int, manyBookChapterArr [][]models.ChapterSpider) (int, error) {
	// 初始化
	okTotal := 0 // 插入成功总数

	for index, oneBookChapterArr := range manyBookChapterArr {
		log.Warnf("delete index=%v len=%v------------ oneBookChapterArr=%v, ", index, len(oneBookChapterArr), oneBookChapterArr)
		// 步骤4: 数据清洗/ 未爬到的字段赋值
		// -- 赋值上下文参数 + 数据清洗。（赋值上下文参数：是吧方法传参，给对象赋值。数据清洗：设置-爬取字段，或者默认数据）
		for i := range oneBookChapterArr {
			// -赋值 上下文传参。如parentId (非数据清洗业务，放在这里)
			oneBookChapterArr[i].ParentId = int(bookIdArr[index]) // 父id
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
		// 找到最后一章，从chapter里获取需要内容。找 chapterNUm=9999就行
		lastChapter, err := db.DBFindOneByField[models.ChapterSpider]("chapter_num", 9999)
		if err != nil {
			log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 找到最后一章失败, err: ", err)
		}

		// -- 创建 comic_spider_stats对象
		var comicSpiderStats models.ComicSpiderStats
		comicSpiderStats.ComicId = int(bookIdArr[index])
		comicSpiderStats.LatestChapterId = &lastChapter.Id    // 最后章节id
		comicSpiderStats.LatestChapterName = lastChapter.Name // 最后章节名称

		totalChapterDbRealUpsert, err := db.DBCountByField[models.ChapterSpider](db.DBComic, "parent_id", bookIdArr[index]) // db里真实插入 章节个数
		errorutil.ErrorPrint(err, "爬取oneBookAllChapter, 插入chapter_spider表后, 查询总插入数出错, err = ")
		comicSpiderStats.TotalChapter = totalChapterDbRealUpsert // 总章节数，从数据库查的
		okTotal += totalChapterDbRealUpsert

		// -- 更新 comic_spider_stats
		log.Info("------ update comicSpiderStats = ", comicSpiderStats)
		log.Info("------ update comicSpiderStats.ComicId = ", comicSpiderStats.ComicId)
		log.Info("------ update comicSpiderStats.LatestChapterId = ", *comicSpiderStats.LatestChapterId)
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
