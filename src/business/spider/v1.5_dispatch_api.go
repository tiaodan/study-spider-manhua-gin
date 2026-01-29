/*
功能: 处理爬取api, V1.5版本实现方式
*/

package spider

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// ------------------------------------------- 初始化 -------------------------------------------
// 请求结构体(带验证规则) - 爬oneBookAllChapter
type SpiderOneBookAllChapterReqV15 struct {
	SpiderTag struct {
		Website string `json:"website" binding:"required" ` // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	BookId int `json:"bookId" binding:"required,min=1" ` // 必填, 且大于0。 gt=0 这样写也可以
}

// 请求结构体(带验证规则) - 爬manyBookAllChapter
type SpiderManyBookAllChapterReqV15 struct {
	SpiderTag struct {
		Website string `json:"website" binding:"required" ` // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	BookIdArr []int `json:"bookIdArr" binding:"required,min=1,dive,gt=0" ` // required 必填, min=1 数组长度最小为1, dive 判断每个子元素, gt=0 个元素必须 > 0
}

// 请求结构体(带验证规则) - 爬manyBookAllChapter V2, 增加字段 BookEverytimeMax
type SpiderManyBookAllChapterReqV15V2 struct {
	SpiderTag struct {
		Website string `json:"website" binding:"required" ` // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	BookIdArr        []int `json:"bookIdArr" binding:"required,min=1,dive,gt=0" ` // required 必填, min=1 数组长度最小为1, dive 判断每个子元素, gt=0 个元素必须 > 0
	BookEverytimeMax int   `json:"bookEverytimeMax" binding:"required,min=1" `    // 可选, 数得 >0。解释: book每组最大请求个数
}

// 请求结构体(带验证规则) - 爬manyBookAllChapter V1, 增加字段 BookEverytimeMax
type SpiderManyChapterAllContentReqV15V1 struct {
	SpiderTag struct {
		Website string `json:"website" binding:"required" ` // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	ChapterIdArr        []int `json:"chapterIdArr" binding:"required,min=1,dive,gt=0" ` // required 必填, min=1 数组长度最小为1, dive 判断每个子元素, gt=0 个元素必须 > 0
	ChapterEverytimeMax int   `json:"chapterEverytimeMax" binding:"required,min=1" `    // 必填, 数得 >0。解释: chapter每组最大请求个数
}

// 请求结构体(带验证规则) - 爬manyBookAllChapter V1, 增加字段 spiderTag.websiteId 必填
type SpiderManyChapterAllContentReqV15V2 struct {
	SpiderTag struct {
		WebsiteId int    `json:"websiteId" binding:"required" ` // 必填。required 同时满足 "非空字符串"
		Website   string `json:"website" binding:"required" `   // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	ChapterIdArr        []int `json:"chapterIdArr" binding:"required,min=1,dive,gt=0" ` // required 必填, min=1 数组长度最小为1, dive 判断每个子元素, gt=0 个元素必须 > 0
	ChapterEverytimeMax int   `json:"chapterEverytimeMax" binding:"required,min=1" `    // 必填, 数得 >0。解释: chapter每组最大请求个数
}

// 请求结构体(带验证规则) - 爬manyBookAllChapter V3, 可以2种方式传数组：1. 传一个随机数组 2 传开始、结束号码，用这个的时候，ChapterIdArr 传空数组就行
type SpiderManyChapterAllContentReqV15V3 struct {
	SpiderTag struct {
		WebsiteId int    `json:"websiteId" binding:"required" ` // 必填。required 同时满足 "非空字符串"
		Website   string `json:"website" binding:"required" `   // 必填。required 同时满足 "非空字符串"
	} `json:"spiderTag" binding:"required" ` // 必填
	ChapterIdArr        []int `json:"chapterIdArr" binding:"required,min=0,dive,gt=0" ` // required 必填, V3 min=0 数组长度最小为0, dive 判断每个子元素, gt=0 个元素必须 > 0。与V2区别：V2是数组长度最小1
	ChapterEverytimeMax int   `json:"chapterEverytimeMax" binding:"required,min=1" `    // 必填, 数得 >0。解释: chapter每组最大请求个数
	ChapterIdStart      *int  `json:"chapterIdStart" binding:"omitempty,gt=0" `         // 可选, 大于0。解释: chapterId开始号码
	ChapterIdEnd        *int  `json:"chapterIdEnd" binding:"omitempty,gt=0" `           // 可选, 大于0。解释: chapterId结束号码
}

// ------------------------------------------- 各种处理API 方法 -------------------------------------------

// 爬某一类所有书籍 - V1.5版本实现方式: 从配置文件读参数 ,此方法的V2 实现
// 与V1主要差别：移除 switch case 区分网站，进行对应网站的爬取逻辑
/*
参考通用思路：
 1. 校验传参
 2. 数据清洗
 3. 业务逻辑 需要的数据校验 +清洗
 4. 执行核心逻辑
	- 读取html内容
	- 通过mapping 爬取字段，赋值给chapter_spider对象
	- 验证业务逻辑，保证稳定性(比如 websiteId是否存在, countryId是否存在等)
	- 插入前, 数据清洗
	- 批量插入db
 5. 返回结果
*/
func DispatchApi_SpiderOneTypeAllBookArr_V1_5_V2(c *gin.Context) {
	// V0.2 要排查 comic_spider + comci_spider_stats 插入数据不一致问题
	// 0. 初始化
	okTotal := 0        // 成功条数
	onePageOkTotal := 0 // 每页成功条数

	// 1. 校验传参
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 找到应该爬哪个网站
	// 读取 JSON Body --
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
		return
	}

	// 1. gjson 读取 前端 JSON 里 有用数据
	website := gjson.Get(string(data), "spiderTag.website").String() // websiteTag - website 字段
	spiderUrl := gjson.Get(string(data), "spiderUrl").String()       // spiderUrl 要爬取的url。这里传的是某个分类的 url (页码用%d替代)。如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
	startPageNum := gjson.Get(string(data), "startPageNum").Int()    // 爬取起始页码 。如： 1-10页，startPageNum=1
	endPageNum := gjson.Get(string(data), "endPageNum").Int()        // 爬取页码结束数。如： 1-10页，endNum=10
	websiteId := int(gjson.Get(string(data), "websiteId").Int())

	// 2. 生成 爬取的url 数组
	// 判断传参 --
	// endPageNum > startPageNUm
	totalPages := int(endPageNum - startPageNum + 1)
	if totalPages <= 0 {
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneTypeAllBookArr_V1_5, endPageNum必须大于等于startPageNum"})
		return
	}

	spiderUrlArr := make([]string, totalPages)
	for i := range spiderUrlArr {
		spiderUrlArr[i] = fmt.Sprintf(spiderUrl, startPageNum+int64(i)) // 如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
		log.Info("delete spiderUrl = ", spiderUrlArr[i])
	}

	// -- 根据网站字段，使用不同的爬虫 ModelMapping映射表
	webCfg := config.CfgSpiderYaml.Websites[website]
	if webCfg == nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 %s 的配置", website)}) // 返回错误
		return
	}
	log.Debug("------- webCfg = ", webCfg)

	// 获取 one_type_all_book 阶段配置
	stageCfg := webCfg.Stages["one_type_all_book"]
	if stageCfg == nil {
		c.JSON(400, gin.H{"error": "func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_type_all_book 阶段的配置"}) // 返回错误
		return
	}

	// 通过mapping 获取 book 对象
	// 插入booK
	// 测试-- mapping结果
	// -- 最终返回结果：二维数组 var AllPageBookArr []onePageBookArr
	allPageBookArr := GetOneTypeAllBookUseCollyByMappingV1[models.ComicSpider](data, ComicMappingForSpiderKxmanhuaByHtml, spiderUrlArr)
	if len(allPageBookArr) == 0 {
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneTypeAllBookArr_V1_5, 获取到的所有书籍为空。推荐排查: 1.爬取网站是不是挂了 2. 本地模拟爬取网站是不是挂了"}) // 获取所有书籍失败
		return
	}
	// log.Debug("---------- 返回 allPageBookArr = ", allPageBookArr)
	log.Debug("---------- stageCfg.Insert.UniqueKeys = ", stageCfg.Insert.UniqueKeys)
	log.Debug("---------- stageCfg.Insert.UpdateKeys = ", stageCfg.Insert.UpdateKeys)

	// 验证业务逻辑
	// websiteId必须 能从数据库 查到 --
	if websiteId <= 0 {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, websiteId必须大于0"})
		return
	}

	// 使用 DBFindOneByField 查询 website 表，验证 websiteId 是否存在 --
	websiteRecord, err := db.DBFindOneByField[models.Website]("id", websiteId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)})
		} else {
			log.Errorf("func=爬一本书chapter 失败, 查询website失败: %v", err)
			log.Error("func=爬一本书chapter 失败, 查询website = ", websiteRecord)
			c.JSON(500, gin.H{"error": "func=爬一本书chapter 失败, 查询website失败"})
		}
		return
	}

	// 2. 从二维数组中，取出每一页,爬的数据，插入数据库
	var debugNewComicIds []int
	for pageIdx, onePageBookArr := range allPageBookArr {
		log.Infof("========== 开始处理第 %d 页，共 %d 条数据 ==========", pageIdx+1, len(onePageBookArr))

		// 2. 插入数据库
		// -- 插入主表
		log.Debugf("第%d页: 开始插入主表 comic_spider，共 %d 条", pageIdx+1, len(onePageBookArr))
		err := db.DBUpsertBatch(db.DBComic, onePageBookArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)
		if err != nil {
			log.Errorf("第%d页: 批量插入db-comic失败, err = %v", pageIdx+1, err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeAllBook), 第%d页批量插入db-comic失败: %v", pageIdx+1, err)}) // 返回错误
			return
		}
		log.Infof("第%d页: 主表 comic_spider 插入成功，共 %d 条", pageIdx+1, len(onePageBookArr))

		// -- 重要：由于GORM批量Upsert时不会更新对象的ID字段，需要重新查询获取正确的ID
		// 构建查询条件：根据唯一索引字段查询
		comicIdQueryFailedCount := 0
		comicIdZeroCount := 0

		for comicIdx := range onePageBookArr {
			comic := &onePageBookArr[comicIdx]
			// 使用对象作为查询条件
			existingComic, err := db.DBFindOneByUniqueIndexMapCondition(comic, stageCfg.Insert.UniqueKeys)
			if err == nil {
				// 更新对象的ID为数据库中的实际ID
				oldId := comic.Id
				comic.Id = existingComic.Id
				if comic.Id <= 0 {
					comicIdZeroCount++
					log.Errorf("第%d页第%d条: 查询到的comic ID无效 (ID=%d), comic名称=%s, 唯一索引字段=%v",
						pageIdx+1, comicIdx+1, comic.Id, comic.Name, stageCfg.Insert.UniqueKeys)
				} else {
					log.Debugf("第%d页第%d条: 更新comic ID成功, 名称=%s, 旧ID=%d -> 新ID=%d",
						pageIdx+1, comicIdx+1, comic.Name, oldId, comic.Id)
					debugNewComicIds = append(debugNewComicIds, comic.Id) // 测试 delete
				}
			} else {
				comicIdQueryFailedCount++
				log.Errorf("第%d页第%d条: 查询comicId失败, comic名称=%s, 唯一索引字段=%v, err = %v",
					pageIdx+1, comicIdx+1, comic.Name, stageCfg.Insert.UniqueKeys, err)
				// 打印comic的详细信息，便于排查
				log.Errorf("第%d页第%d条: comic详细信息 - Name=%s, WebsiteId=%d, CountryId=%d, TypeId=%d, PornTypeId=%d, ProcessId=%d, AuthorConcat=%s",
					pageIdx+1, comicIdx+1, comic.Name, comic.WebsiteId, comic.CountryId, comic.TypeId, comic.PornTypeId, comic.ProcessId, comic.AuthorConcat)
				c.JSON(500, gin.H{
					"error": fmt.Sprintf("/spider/oneTypeAllBook失败,第%d页第%d条更新关联表前查询comicId失败, comic名称=%s, err=%v",
						pageIdx+1, comicIdx+1, comic.Name, err),
					"pageIndex":  pageIdx + 1,
					"comicIndex": comicIdx + 1,
					"comicName":  comic.Name,
				})
				return // 不进行下一步
			}
		}

		// 检查是否有ID为0的情况
		if comicIdZeroCount > 0 {
			log.Errorf("第%d页: 警告！有 %d 条comic的ID为0，这会导致stats插入失败", pageIdx+1, comicIdZeroCount)
			c.JSON(500, gin.H{
				"error":       fmt.Sprintf("第%d页: 有 %d 条comic的ID为0，无法插入stats表", pageIdx+1, comicIdZeroCount),
				"pageIndex":   pageIdx + 1,
				"zeroIdCount": comicIdZeroCount,
			})
			return
		}

		log.Infof("第%d页: 所有comic ID查询成功，共 %d 条", pageIdx+1, len(onePageBookArr))

		// -- 插入关联表
		var comicStatsArr []models.ComicSpiderStats
		statsBuildFailedCount := 0
		for comicIdx, comic := range onePageBookArr {
			if comic.Id <= 0 {
				statsBuildFailedCount++
				log.Errorf("第%d页第%d条: 构建stats失败，comic ID无效 (ID=%d), comic名称=%s",
					pageIdx+1, comicIdx+1, comic.Id, comic.Name)
				continue
			}

			stats := models.ComicSpiderStats{
				ComicId:                   comic.Id, // 现在使用正确的ID
				Star:                      comic.Stats.Star,
				LatestChapterName:         comic.Stats.LatestChapterName, // 最新章节名字
				Hits:                      comic.Stats.Hits,
				TotalChapter:              comic.Stats.TotalChapter,
				LastestChapterReleaseDate: comic.Stats.LastestChapterReleaseDate,
			}
			stats.DataClean() // 数据清洗下
			comicStatsArr = append(comicStatsArr, stats)
			log.Debugf("第%d页第%d条: 构建stats成功, ComicId=%d, Star=%.2f, Hits=%d",
				pageIdx+1, comicIdx+1, stats.ComicId, stats.Star, stats.Hits)
		}

		if statsBuildFailedCount > 0 {
			log.Errorf("第%d页: 构建stats失败 %d 条，实际构建成功 %d 条", pageIdx+1, statsBuildFailedCount, len(comicStatsArr))
			c.JSON(500, gin.H{
				"error":       fmt.Sprintf("第%d页: 构建stats失败 %d 条", pageIdx+1, statsBuildFailedCount),
				"pageIndex":   pageIdx + 1,
				"failedCount": statsBuildFailedCount,
			})
			return
		}

		if len(comicStatsArr) != len(onePageBookArr) {
			log.Errorf("第%d页: 数据不一致！主表插入 %d 条，但stats只构建了 %d 条",
				pageIdx+1, len(onePageBookArr), len(comicStatsArr))
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("第%d页: 数据不一致！主表 %d 条，stats %d 条",
					pageIdx+1, len(onePageBookArr), len(comicStatsArr)),
				"pageIndex":      pageIdx + 1,
				"mainTableCount": len(onePageBookArr),
				"statsCount":     len(comicStatsArr),
			})
			return
		}

		log.Debugf("第%d页: 开始插入关联表 comic_spider_stats，共 %d 条", pageIdx+1, len(comicStatsArr))
		log.Debugf("第%d页: stats唯一索引字段=%v, 更新字段=%v",
			pageIdx+1, stageCfg.RelatedTables["comic_stats"].Insert.UniqueKeys,
			stageCfg.RelatedTables["comic_stats"].Insert.UpdateKeys)

		err = db.DBUpsertBatch(db.DBComic, comicStatsArr, stageCfg.RelatedTables["comic_stats"].Insert.UniqueKeys, stageCfg.RelatedTables["comic_stats"].Insert.UpdateKeys)
		if err != nil {
			log.Errorf("第%d页: 批量插入db-comic-stats表失败, err = %v", pageIdx+1, err)
			log.Errorf("第%d页: 主表已插入 %d 条，但stats表插入失败，数据不一致！", pageIdx+1, len(onePageBookArr))
			// 打印前几条stats数据，便于排查
			for i := 0; i < len(comicStatsArr) && i < 3; i++ {
				log.Errorf("第%d页: stats[%d] = ComicId=%d, Star=%.2f, Hits=%d",
					pageIdx+1, i, comicStatsArr[i].ComicId, comicStatsArr[i].Star, comicStatsArr[i].Hits)
			}
			c.JSON(400, gin.H{
				"error":             fmt.Sprintf("/spider/oneTypeAllBook失败,第%d页批量插入db-comic-stats表失败: %v", pageIdx+1, err),
				"pageIndex":         pageIdx + 1,
				"mainTableInserted": len(onePageBookArr),
				"statsInsertFailed": true,
			}) // 返回错误
			return
		}
		log.Infof("第%d页: 关联表 comic_spider_stats 插入成功，共 %d 条", pageIdx+1, len(comicStatsArr))

		// 打印结果
		onePageOkTotal = len(onePageBookArr) // 每页成功条数
		okTotal += onePageOkTotal            // 总成功条数
		log.Infof("========== 第%d页处理完成，成功插入 %d 条 book 数据（主表+stats表） ==========", pageIdx+1, onePageOkTotal)
	}

	// 测试-一会删除，遍历内容，判断 debugNewComicIds 里是否有重复数据
	if util.HasDuplicate(debugNewComicIds) {
		log.Warn("!!! 检测到 debugNewComicIds 中有重复数据, 如果comic_spider 和 stats表数据个数不一致, 就需要检查！")
	}

	log.Infof("爬取某个分类allBook,, 爬取成功, 插入%d条 book 数据", okTotal)

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条数据")

	/* V0.1 要排查 comic_spider + comci_spider_stats 插入数据不一致问题前，备份
	// 0. 初始化
	okTotal := 0        // 成功条数
	onePageOkTotal := 0 // 每页成功条数

	// 1. 校验传参
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 找到应该爬哪个网站
	// 读取 JSON Body --
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
		return
	}

	// 1. gjson 读取 前端 JSON 里 有用数据
	website := gjson.Get(string(data), "spiderTag.website").String() // websiteTag - website 字段
	spiderUrl := gjson.Get(string(data), "spiderUrl").String()       // spiderUrl 要爬取的url。这里传的是某个分类的 url (页码用%d替代)。如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
	startPageNum := gjson.Get(string(data), "startPageNum").Int()    // 爬取起始页码 。如： 1-10页，startPageNum=1
	endPageNum := gjson.Get(string(data), "endPageNum").Int()        // 爬取页码结束数。如： 1-10页，endNum=10
	websiteId := int(gjson.Get(string(data), "websiteId").Int())

	// 2. 生成 爬取的url 数组
	// 判断传参 --
	// endPageNum > startPageNUm
	totalPages := int(endPageNum - startPageNum + 1)
	if totalPages <= 0 {
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneTypeAllBookArr_V1_5, endPageNum必须大于等于startPageNum"})
		return
	}

	spiderUrlArr := make([]string, totalPages)
	for i := range spiderUrlArr {
		spiderUrlArr[i] = fmt.Sprintf(spiderUrl, startPageNum+int64(i)) // 如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
		log.Info("----- delete spiderUrl = ", spiderUrlArr[i])
	}

	// -- 根据网站字段，使用不同的爬虫 ModelMapping映射表
	webCfg := config.CfgSpiderYaml.Websites[website]
	if webCfg == nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 %s 的配置", website)}) // 返回错误
		return
	}
	log.Debug("------- webCfg = ", webCfg)

	// 获取 one_type_all_book 阶段配置
	stageCfg := webCfg.Stages["one_type_all_book"]
	if stageCfg == nil {
		c.JSON(400, gin.H{"error": "func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_type_all_book 阶段的配置"}) // 返回错误
		return
	}

	// 通过mapping 获取 book 对象
	// 插入booK
	// 测试-- mapping结果
	// -- 最终返回结果：二维数组 var AllPageBookArr []onePageBookArr
	allPageBookArr := GetOneTypeAllBookUseCollyByMappingV1[models.ComicSpider](data, ComicMappingForSpiderKxmanhuaByHtml, spiderUrlArr)
	if len(allPageBookArr) == 0 {
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneTypeAllBookArr_V1_5, 获取到的所有书籍为空。推荐排查: 1.爬取网站是不是挂了 2. 本地模拟爬取网站是不是挂了"}) // 获取所有书籍失败
		return
	}
	// log.Debug("---------- 返回 allPageBookArr = ", allPageBookArr)
	log.Debug("---------- stageCfg.Insert.UniqueKeys = ", stageCfg.Insert.UniqueKeys)
	log.Debug("---------- stageCfg.Insert.UpdateKeys = ", stageCfg.Insert.UpdateKeys)

	// 验证业务逻辑
	// websiteId必须 能从数据库 查到 --
	if websiteId <= 0 {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, websiteId必须大于0"})
		return
	}

	// 使用 DBFindOneByField 查询 website 表，验证 websiteId 是否存在 --
	websiteRecord, err := db.DBFindOneByField[models.Website]("id", websiteId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)})
		} else {
			log.Errorf("func=爬一本书chapter 失败, 查询website失败: %v", err)
			log.Error("func=爬一本书chapter 失败, 查询website = ", websiteRecord)
			c.JSON(500, gin.H{"error": "func=爬一本书chapter 失败, 查询website失败"})
		}
		return
	}

	// 2. 从二维数组中，取出每一页,爬的数据，插入数据库
	for i, onePageBookArr := range allPageBookArr {
		// 2. 插入数据库
		// -- 插入主表
		err := db.DBUpsertBatch(db.DBComic, onePageBookArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)
		if err != nil {
			c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeAllBook), 批量插入db-comic 失败"}) // 返回错误
			return
		}

		// -- 重要：由于GORM批量Upsert时不会更新对象的ID字段，需要重新查询获取正确的ID
		// 构建查询条件：根据唯一索引字段查询
		for i := range onePageBookArr {
			// 使用对象作为查询条件
			existingComic, err := db.DBFindOneByUniqueIndexMapCondition(&onePageBookArr[i], stageCfg.Insert.UniqueKeys)
			if err == nil {
				// 更新对象的ID为数据库中的实际ID
				onePageBookArr[i].Id = existingComic.Id
				log.Debugf("更新comic ID: %s -> %d", onePageBookArr[i].Name, existingComic.Id)
			} else {
				log.Errorf("/spider/oneTypeAllBook失败,更新关联表前,查询comicId %v失败, err = %v", onePageBookArr[i].Name, err)
				c.JSON(500, gin.H{"error": "/spider/oneTypeAllBook失败,更新关联表前,查询comicId失败"})
				return // 不进行下一步
			}
		}

		// -- 插入关联表
		var comicStatsArr []models.ComicSpiderStats
		for _, comic := range onePageBookArr {
			stats := models.ComicSpiderStats{
				ComicId:                   comic.Id, // 现在使用正确的ID
				Star:                      comic.Stats.Star,
				LatestChapterName:         comic.Stats.LatestChapterName, // 最新章节名字
				Hits:                      comic.Stats.Hits,
				TotalChapter:              comic.Stats.TotalChapter,
				LastestChapterReleaseDate: comic.Stats.LastestChapterReleaseDate,
			}
			stats.DataClean() // 数据清洗下
			comicStatsArr = append(comicStatsArr, stats)
		}
		err = db.DBUpsertBatch(db.DBComic, comicStatsArr, stageCfg.RelatedTables["comic_stats"].Insert.UniqueKeys, stageCfg.RelatedTables["comic_stats"].Insert.UpdateKeys)
		if err != nil {
			c.JSON(400, gin.H{"error": "/spider/oneTypeAllBook), 批量插入db-comic-stats表 失败"}) // 返回错误
			return
		}

		// 打印结果
		onePageOkTotal = len(onePageBookArr) // 每页成功条数
		okTotal += onePageOkTotal            // 总成功条数
		log.Infof("爬取某个分类allBook, 第%d页, 爬取成功, 插入%d条 book 数据", i+1, onePageOkTotal)
	}
	log.Infof("爬取某个分类allBook,, 爬取成功, 插入%d条 book 数据", okTotal)

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条数据")
	*/
}

// 爬某一类所有书籍 - V1.5版本实现方式: 从配置文件读参数 ,此方法的V1 实现
// 与V2主要差别：有 switch case 区分网站，进行对应网站的爬取逻辑
/*
参考通用思路：
 1. 校验传参
 2. 数据清洗
 3. 业务逻辑 需要的数据校验 +清洗
 4. 执行核心逻辑
	- 读取html内容
	- 通过mapping 爬取字段，赋值给chapter_spider对象
	- 验证业务逻辑，保证稳定性(比如 websiteId是否存在, countryId是否存在等)
	- 插入前, 数据清洗
	- 批量插入db
 5. 返回结果
*/
func DispatchApi_SpiderOneTypeAllBookArr_V1_5_V1(c *gin.Context) {
	// 0. 初始化
	okTotal := 0        // 成功条数
	onePageOkTotal := 0 // 每页成功条数

	// 1. 校验传参
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 找到应该爬哪个网站
	// 读取 JSON Body --
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
		return
	}

	// 1. gjson 读取 前端 JSON 里 有用数据
	website := gjson.Get(string(data), "spiderTag.website").String()   // websiteTag - website 字段
	spiderUrl := gjson.Get(string(data), "spiderUrl").String()         // spiderUrl 要爬取的url。这里传的是某个分类的 url (页码用%d替代)。如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
	startPageNum := int(gjson.Get(string(data), "startPageNum").Int()) // 爬取起始页码 。如： 1-10页，startPageNum=1
	endPageNum := int(gjson.Get(string(data), "endPageNum").Int())     // 爬取页码结束数。如： 1-10页，endNum=10
	websiteId := int(gjson.Get(string(data), "websiteId").Int())

	// 2. 生成 爬取的url 数组
	// 判断传参 --
	// endPageNum 必须 >= startPageNUm
	if endPageNum < startPageNum {
		c.JSON(400, gin.H{"error": "func: 爬取oneType V1.5 失败。endPageNum 必须 >= startPageNUm"})
		return
	}
	spiderUrlArr := make([]string, endPageNum)
	for i := range spiderUrlArr {
		// spiderUrlArr[i] = fmt.Sprintf(spiderUrl, i+1) // 原来写法：如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
		spiderUrlArr[i] = fmt.Sprintf(spiderUrl, startPageNum+i) // 如："https://kxmanhua.com/manga/library?type=2&complete=1&page=%d&orderby=1"
		log.Info("delete spiderUrl = ", spiderUrlArr[i])
	}

	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// 0. 使用 kxmanhua 相关配置
	webCfg := config.CfgSpiderYaml.Websites[website]
	if webCfg == nil {
		c.JSON(400, gin.H{"error": "func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 kxmanhua 的配置"}) // 返回错误
		return
	}
	log.Debug("webCfg = ", webCfg)

	// 获取 one_type_all_book 阶段配置
	stageCfg := webCfg.Stages["one_type_all_book"]
	if stageCfg == nil {
		c.JSON(400, gin.H{"error": "func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_type_all_book 阶段的配置"}) // 返回错误
		return
	}

	// 通过mappping 获取 book 对象
	// 插入booK
	// 测试-- mapping结果
	// -- 最终返回结果：二维数组 var AllPageBookArr []onePageBookArr
	allPageBookArr := GetOneTypeAllBookUseCollyByMappingV1[models.ComicSpider](data, ComicMappingForSpiderKxmanhuaByHtml, spiderUrlArr)
	if len(allPageBookArr) == 0 {
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneTypeAllBookArr_V1_5, 获取到的所有书籍为空。推荐排查: 1.爬取网站是不是挂了 2. 本地模拟爬取网站是不是挂了"}) // 获取所有书籍失败
		return
	}
	// log.Debug("---------- 返回 allPageBookArr = ", allPageBookArr)
	log.Debug("---------- stageCfg.Insert.UniqueKeys = ", stageCfg.Insert.UniqueKeys)
	log.Debug("---------- stageCfg.Insert.UpdateKeys = ", stageCfg.Insert.UpdateKeys)

	// 验证业务逻辑
	// websiteId必须 能从数据库 查到 --
	if websiteId <= 0 {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, websiteId必须大于0"})
		return
	}

	// 使用 DBFindOneByField 查询 website 表，验证 websiteId 是否存在 --
	websiteRecord, err := db.DBFindOneByField[models.Website]("id", websiteId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)})
		} else {
			log.Errorf("func=爬一本书chapter 失败, 查询website失败: %v", err)
			log.Error("func=爬一本书chapter 失败, 查询website = ", websiteRecord)
			c.JSON(500, gin.H{"error": "func=爬一本书chapter 失败, 查询website失败"})
		}
		return
	}

	// 2. 从二维数组中，取出每一页,爬的数据，插入数据库
	for i, onePageBookArr := range allPageBookArr {
		// 2. 插入数据库
		// -- 插入主表
		err := db.DBUpsertBatch(db.DBComic, onePageBookArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)
		if err != nil {
			c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeAllBook), 批量插入db-comic 失败"}) // 返回错误
			return
		}

		// -- 重要：由于GORM批量Upsert时不会更新对象的ID字段，需要重新查询获取正确的ID
		// 构建查询条件：根据唯一索引字段查询
		for i := range onePageBookArr {
			/* 之前代码
			var existingComic models.ComicSpider
			condition := map[string]interface{}{
				"name":          onePageBookArr[i].Name,
				"country_id":    onePageBookArr[i].CountryId,
				"website_id":    onePageBookArr[i].WebsiteId,
				"porn_type_id":  onePageBookArr[i].PornTypeId,
				"type_id":       onePageBookArr[i].TypeId,
				"author_concat": onePageBookArr[i].AuthorConcat,
			}
			result := db.DBComic.Where(condition).First(&existingComic)
			if result.Error == nil {
				// 更新对象的ID为数据库中的实际ID
				onePageBookArr[i].Id = existingComic.Id
				log.Debugf("更新comic ID: %s -> %d", onePageBookArr[i].Name, existingComic.Id)
			} else {
				log.Errorf("查询comic失败: %s, err: %v", onePageBookArr[i].Name, result.Error)
			}
			*/
			// 使用对象作为查询条件
			existingComic, err := db.DBFindOneByUniqueIndexMapCondition(&onePageBookArr[i], stageCfg.Insert.UniqueKeys)
			if err == nil {
				// 更新对象的ID为数据库中的实际ID
				onePageBookArr[i].Id = existingComic.Id
				log.Debugf("更新comic ID: %s -> %d", onePageBookArr[i].Name, existingComic.Id)
			} else {
				log.Errorf("/spider/oneTypeAllBook失败,更新关联表前,查询comicId %v失败, err = %v", onePageBookArr[i].Name, err)
				c.JSON(500, gin.H{"error": "/spider/oneTypeAllBook失败,更新关联表前,查询comicId失败"})
				return // 不进行下一步
			}
		}

		// -- 插入关联表
		var comicStatsArr []models.ComicSpiderStats
		for _, comic := range onePageBookArr {
			stats := models.ComicSpiderStats{
				ComicId:                   comic.Id, // 现在使用正确的ID
				Star:                      comic.Stats.Star,
				LatestChapterName:         comic.Stats.LatestChapterName, // 最新章节名字
				Hits:                      comic.Stats.Hits,
				TotalChapter:              comic.Stats.TotalChapter,
				LastestChapterReleaseDate: comic.Stats.LastestChapterReleaseDate,
			}
			stats.DataClean() // 数据清洗下
			comicStatsArr = append(comicStatsArr, stats)
		}
		err = db.DBUpsertBatch(db.DBComic, comicStatsArr, stageCfg.RelatedTables["comic_stats"].Insert.UniqueKeys, stageCfg.RelatedTables["comic_stats"].Insert.UpdateKeys)
		if err != nil {
			c.JSON(400, gin.H{"error": "/spider/oneTypeAllBook), 批量插入db-comic-stats表 失败"}) // 返回错误
			return
		}

		// 打印结果
		onePageOkTotal = len(onePageBookArr) // 每页成功条数
		okTotal += onePageOkTotal            // 总成功条数
		log.Infof("爬取某个分类allBook, 第%d页, 爬取成功, 插入%d条 book 数据", i+1, onePageOkTotal)
	}
	log.Infof("爬取某个分类allBook,, 爬取成功, 插入%d条 book 数据", okTotal)

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条数据")

}

// 爬某一本书所有章节 - V1.5版本实现方式 V2方法实现: 通过配置 实现代码,并把switch 去掉
/*
步骤：
	0. 初始化
	1. 获取传参。实现方式: gjson框架 Get()实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走)
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值
		步骤5: 验证爬取 数据准确性
		步骤6: 数据库插入
	6. 返回结果

	哪些步骤可以组合成1个方法
*/
func DispatchApi_OneBookAllChapter_V1_5_V2(c *gin.Context) {
	// v2 写法, 把switch 去掉,并通过配置 实现代码
	// 0. 初始化
	okTotal := 0 // 成功条数

	// 1. 校验传参
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 找到应该爬哪个网站
	// 读取 JSON Body --
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
		return
	}

	// gjson 读取 前端 JSON 里 spiderTag -> website字段 --
	website := gjson.Get(string(data), "spiderTag.website").String()
	bookId := gjson.Get(string(data), "bookId").Int() // bookId字段
	// 校验 --
	// bookId !=0
	if bookId == 0 {
		c.JSON(500, gin.H{"error": "爬取1本书所有章节 失败, bookId为0, 不执行后续步骤"}) // 返回错误
		return
	}
	log.Info("bookId = ", bookId)

	// -- 根据网站字段，使用不同的爬虫 ModelMapping映射表
	webCfg := config.CfgSpiderYaml.Websites[website]
	if webCfg == nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 %s 的配置", website)}) // 返回错误
		return
	}
	log.Debug("webCfg = ", webCfg)

	// 获取 one_book_all_chapter 阶段配置
	stageCfg := webCfg.Stages["one_book_all_chapter"]
	if stageCfg == nil {
		c.JSON(400, gin.H{"error": "func=爬 爬oneTypeAllBookArr V1.5, 配置文件里没有找到 one_book_all_chapter 阶段的配置"}) // 返回错误
		return
	}

	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterMappingForSpiderKxmanhuaByHTML,
	}
	mapping := mappingFactory[website]

	// 2. 爬取 chapter
	// -- 请求html页面
	chapterArr := GetOneBookAllChapterByCollyMapping[models.ChapterSpider](data, mapping.(map[string]models.ModelHtmlMapping))
	// -- 插入前数据校验
	if chapterArr == nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db")
		c.JSON(400, gin.H{"error": "爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db"}) // 返回错误
		return                                                                                    // 直接结束
	}

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

	// 3. upsert chapter
	err = db.DBUpsertBatch(db.DBComic, chapterArr, stageCfg.Insert.UniqueKeys, stageCfg.Insert.UpdateKeys)

	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml(分发api- /spider/oneBookAllChapterByHtml), 批量插入db chapter 失败, err: ", err)
		c.JSON(500, gin.H{"error": "批量插入db chapter 失败"}) // 返回错误
	}

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
	okTotal = totalChapterDbRealUpsert

	// -- 更新 comic_spider_stats
	log.Info("update comicSpiderStats = ", comicSpiderStats)
	log.Info("update comicSpiderStats.ComicId = ", comicSpiderStats.ComicId)
	log.Info("update comicSpiderStats.LatestChapterId = ", *comicSpiderStats.LatestChapterId)
	err = db.DBUpdate(db.DBComic, &comicSpiderStats, stageCfg.UpdateParentStats.UniqueKeys, stageCfg.UpdateParentStats.UpdateKeys)
	if err != nil {
		log.Error("func= DispatchApi_OneBookAllChapterByHtml, 更新comic_spider_stats失败, err: ", err)
		c.JSON(500, gin.H{"error": "更新comic_spider_stats失败"}) // 返回错误
		return
	}

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapter数据")
}

// 爬某一本书所有章节 - V1.5版本实现方式 V3方法实现
// 主要改动：基于V2,将代码分成小方法，容易看，整洁。要不一个方法120行，看着乱
// 基于V2实现: 通过配置 实现代码,并把switch 去掉
/*
步骤：
	0. 初始化
	1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
		- shouldBIndJson已经包含 validator验证了
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
	6. 返回结果

	哪些步骤可以组合成1个方法
*/
func DispatchApi_OneBookAllChapter_V1_5_V3(c *gin.Context) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "爬OneBookAllChapter"
	var funcErr error

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	var req SpiderOneBookAllChapterReqV15
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 失败, 获取前端传参失败: %v", funcName, err)})
		return
	}
	websiteName := req.SpiderTag.Website
	bookId := req.BookId
	log.Infof("func=%v, 要爬的bookId = %v", funcName, bookId)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterMappingForSpiderKxmanhuaByHTML,
	}
	mapping := mappingFactory[websiteName]

	// 2. 爬取 chapter
	// -- 请求html页面
	chapterArr, err := GetOneBookAllChapterByCollyMappingV1_5[models.ChapterSpider](mapping.(map[string]models.ModelHtmlMapping), bookId)
	// -- 插入前数据校验
	if chapterArr == nil || err != nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		c.JSON(400, gin.H{"error": "爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到"}) // 返回错误
		return                                                                                                                        // 直接结束
	}

	// 4. 执行核心逻辑 - 插入部分
	if okTotal, funcErr = SpiderOneBookAllChapter_UpsertPart(websiteName, bookId, chapterArr); funcErr != nil {
		c.JSON(500, gin.H{"error": "爬取失败"})
	}

	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapter数据")
}

// 爬某多本书所有章节 - V1.5版本实现方式 V1方法实现
// 主要改动：将代码分成小方法，容易看，整洁。要不一个方法120行，看着乱
// 基于 DispatchApi_OneBookAllChapter_V1_5_V2 实现: 通过配置 实现代码,并把switch 去掉
/*
步骤：
	0. 初始化
	1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
		- shouldBIndJson已经包含 validator验证了
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤1: 找到目标网站
		步骤2: 爬取
		步骤3: 提取数据
		步骤4: 数据清洗/ 未爬到的字段赋值
		步骤5: 验证爬取数据 准确性
		步骤6: 数据库插入
	6. 返回结果

	哪些步骤可以组合成1个方法
*/
func DispatchApi_ManyBookAllChapter_V1_5_V1(c *gin.Context) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "爬ManyBookAllChapter"
	var funcErr error

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	var req SpiderManyBookAllChapterReqV15
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 失败, 获取前端传参失败: %v", funcName, err)})
		return
	}
	websiteName := req.SpiderTag.Website
	bookIdArr := req.BookIdArr
	log.Infof("func=%v, 要爬的bookId = %v", funcName, bookIdArr)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	// -- 从mapping 工厂了拿数据
	var mappingFactory = map[string]any{
		"kxmanhua": ChapterMappingForSpiderKxmanhuaByHTML,
	}
	mapping := mappingFactory[websiteName]

	// 2. 爬取 chapter
	// -- 请求html页面
	manyBookChapterArrMap, err := GetManyBookAllChapterByCollyMappingV1_5[models.ChapterSpider](mapping.(map[string]models.ModelHtmlMapping), websiteName, bookIdArr)
	chapterNamePreviewCount = 0 // ！！！！重要,必有，重置计数器。chapter中 name包含"Preview"次数
	// -- 插入前数据校验
	if manyBookChapterArrMap == nil || err != nil {
		log.Error("爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db。可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到")
		c.JSON(400, gin.H{"error": "爬取 OneBookAllChapterByHtml失败, chapterArr 为空, 拒绝进入下一步: 插入db可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到"}) // 返回错误
		return                                                                                                                        // 直接结束
	}

	// 4. 执行核心逻辑 - 插入部分
	if okTotal, funcErr = SpiderManyBookAllChapter_UpsertPart(websiteName, manyBookChapterArrMap); funcErr != nil {
		c.JSON(500, gin.H{"error": "爬取失败"})
	}

	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapter数据")
}

// 爬某多本书所有章节.！！！解耦： gin.Context和实际爬取、插入逻辑， - V1.5版本实现方式 V2方法实现
// 主要改动：将代码分成小方法，容易看，整洁。要不一个方法120行，看着乱
// 基于 DispatchApi_OneBookAllChapter_V1_5_V2 实现: 通过配置 实现代码,并把switch 去掉
/*
步骤：
	0. 初始化
	1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
		- shouldBIndJson已经包含 validator验证了
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤5.1: 爬取 + 插入 + 更新book 字段：spider_sub_chapter_end_status


			// 步骤1: 找到目标网站
			// 步骤2: 爬取
			// 步骤3: 提取数据
			// 步骤4: 数据清洗/ 未爬到的字段赋值
			// 步骤5: 验证爬取数据 准确性
			// 步骤6: 数据库插入
	6. 返回结果

	哪些步骤可以组合成1个方法
*/
func DispatchApi_ManyBookAllChapter_V1_5_V2(c *gin.Context) {
	// 0. 初始化
	funcName := "爬ManyBookAllChapter"

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	var req SpiderManyBookAllChapterReqV15
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 失败, 获取前端传参失败: %v", funcName, err)})
		return
	}
	websiteName := req.SpiderTag.Website
	bookIdArr := req.BookIdArr
	log.Infof("func=%v, 要爬的bookId = %v", funcName, bookIdArr)

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	okTotal, err := SpiderManyBookAllChapter2DB(websiteName, bookIdArr) // 成功条数
	if err != nil {
		c.JSON(500, "爬取成功, reason: 插入db失败")
	}

	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapter数据")
}

// 爬某多本书所有章节 - V1.5版本实现方式 V2方法实现：自动负载均衡 (根据传参: 每次最多请求几个book)
// 主要改动：将代码分成小方法，容易看，整洁。要不一个方法120行，看着乱
// 基于 DispatchApi_OneBookAllChapter_V1_5_V2 实现: 通过配置 实现代码,并把switch 去掉
/*
功能：
	- 能断点续爬
	- 给的bookIdArr重新排序
	- 自动负载均衡 (根据传参: 每次最多请求几个book)

建议:
	- 建议 mysql 最多1000-5000条, 就是一次操作 10-40个book (每个book 100个章节)

步骤：
	0. 初始化
	1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
		- shouldBIndJson已经包含 validator验证了
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤5.1: 给传的数组重新排序，从小到大
		步骤5.2: 根据传的bookIdArr，根据负载均衡配置，分组
		步骤5.3：for循环, 调接口。
		调接口:DispatchApi_ManyBookAllChapter_V1_5_V1
			步骤1: 找到目标网站
			步骤2: 爬取
			步骤3: 提取数据
			步骤4: 数据清洗/ 未爬到的字段赋值
			步骤5: 验证爬取数据 准确性
			步骤6: 数据库插入
	6. 返回结果

	哪些步骤可以组合成1个方法
*/
func DispatchApi_ManyBookAllChapter_V1_5_V3(c *gin.Context) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "爬ManyBookAllChapter V3"

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	var req SpiderManyBookAllChapterReqV15V2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 失败, 获取前端传参失败: %v", funcName, err)})
		return
	}
	websiteName := req.SpiderTag.Website
	bookIdArr := req.BookIdArr
	bookEverytimeMax := req.BookEverytimeMax

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// 步骤5.1: 给传的数组重新排序，从小到大。slices.Sort() -> 默认从小到大
	slices.Sort(bookIdArr) // 默认从小到大,修改原数组

	// 步骤5.2: 根据传的bookIdArr，根据负载均衡配置，分组
	spilitBookIdArr2D := util.SplitIntArr(bookIdArr, bookEverytimeMax) // 二维数组
	// 步骤5.3：for循环, 调接口
	for i, spilitBookIdArr := range spilitBookIdArr2D {
		// 步骤5.3.1 查这个id是否需要爬 --
		spilitBookIdArrNeedCrawl, err := DBGetIdsNeedCrawlByFiled[models.ComicSpider](db.DBComic, spilitBookIdArr, "spider_sub_chapter_end_status", 0)
		if err != nil {
			log.Errorf("func=%v, 第%v个数组, 爬取前查询是否需要爬 失败, reason: %v", funcName, i, err)
			c.JSON(500, fmt.Sprint("爬取失败, 爬取前查询是否需要爬 失败, 出错数组=", spilitBookIdArr))
			return
		}

		// 如果不需要爬，跳过
		if len(spilitBookIdArrNeedCrawl) == 0 {
			log.Warnf("func=%v, 第%v个数组, 不需要爬, 已经是爬取过的, arr= %v", funcName, i, spilitBookIdArr)
			continue
		}

		// 步骤5.3.2 爬取并插入db --
		log.Infof("func=%v, 要爬的bookIdArr = %v, 过滤需要爬的bookIdArr = %v", funcName, bookIdArr, spilitBookIdArrNeedCrawl)
		okTotalOneArr, err := SpiderManyBookAllChapter2DB(websiteName, spilitBookIdArrNeedCrawl) // 成功条数
		if err != nil {
			log.Errorf("func=%v, 第%v个数组, 爬取失败, reason: %v", funcName, i, err)
			c.JSON(500, fmt.Sprint("爬取失败, 出错数组=", spilitBookIdArr))
			return
		}
		log.Infof("func=%v, 第%v个数组, 爬取成功,total: %v, arr= %v", funcName, i, okTotalOneArr, spilitBookIdArr)
		okTotal += okTotalOneArr
	}

	// 6. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapter数据")
}

// 爬某一章节所有内容 - V1.5版本实现方式
/*
步骤：
	0. 初始化
	1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	2. 校验传参。用validator，需要提前定义 请求结构体(包含校验规则：必有、必须>0等等) -》实现这个之后，再说通过写配置实现这个结构体,而不是总改代码
		- shouldBIndJson已经包含 validator验证了
	3. 前端传参, 数据清洗
	4. 业务逻辑 需要的数据校验 +清洗
	5. 执行核心逻辑 (6步走) : 爬取 | 插入 可以分成2个方法
		步骤5.1: 给传的数组重新排序，从小到大
		步骤5.2: 根据传的bookIdArr，根据负载均衡配置，分组
		步骤5.3：for循环, 调接口。
		调接口:DispatchApi_ManyBookAllChapter_V1_5_V1
			步骤1: 找到目标网站
			步骤2: 爬取
			步骤3: 提取数据
			步骤4: 数据清洗/ 未爬到的字段赋值
			步骤5: 验证爬取数据 准确性
			步骤6: 数据库插入
	6. 返回结果
*/
func DispatchApi_ManyChapterAllContent_V1_5_V1(c *gin.Context) {
	// 0. 初始化
	okTotal := 0 // 成功条数
	funcName := "爬ManyChapterAllContent V1"

	// 1. 获取传参。实现方式: c.ShouldBindJSON(请求结构体)实现
	var req SpiderManyChapterAllContentReqV15V3 // 可以传 start end 表示数组
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 失败, 获取前端传参失败: %v", funcName, err)})
		return
	}
	websiteId := req.SpiderTag.WebsiteId
	websiteName := req.SpiderTag.Website
	chapterIdArr := req.ChapterIdArr
	ChapterEverytimeMax := req.ChapterEverytimeMax

	// 2. 校验传参。用validator，上面shouldBIndJson已经包含 validator验证了
	// 如果 chapterIdArr 传空数组，就使用 start end 章节号码，代替chapterIdArr --
	if len(chapterIdArr) == 0 {
		if req.ChapterIdStart == nil || req.ChapterIdEnd == nil {
			log.Errorf("func=%v 前端传参错误: chapterIdArr 传空数组，请传 start end 章节号码", funcName)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 前端传参错误: chapterIdArr 传空数组，请传 start end 章节号码", funcName)})
			return
		} else if *req.ChapterIdStart > *req.ChapterIdEnd {
			log.Errorf("func=%v 前端传参错误: start章节号码 应小于 end章节号码", funcName)
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=%v 前端传参错误: start章节号码 应小于 end章节号码", funcName)})
			return
		} else { // 应使用 start end 章节号码，代替chapterIdArr
			chapterIdArr = util.GenArr(*req.ChapterIdStart, *req.ChapterIdEnd, 1)
			log.Info("使用start、end生成 chapterId数组 = ", chapterIdArr)
		}

	}

	// 3. 前端传参, 数据清洗
	// 4. 业务逻辑 需要的数据校验 +清洗

	// 5. 执行核心逻辑 (6步走)
	// 步骤5.1: 给传的数组重新排序，从小到大。slices.Sort() -> 默认从小到大
	slices.Sort(chapterIdArr) // 默认从小到大,修改原数组

	// 步骤5.2: 根据传的bookIdArr，根据负载均衡配置，分组
	spilitChapterIdArr2D := util.SplitIntArr(chapterIdArr, ChapterEverytimeMax) // 二维数组
	// 步骤5.3：for循环, 调接口
	for i, spilitChapterIdArr := range spilitChapterIdArr2D {
		// 步骤5.3.1 查这个id是否需要爬 --
		spilitChapterIdArrNeedCrawl, err := DBGetIdsNeedCrawl[models.ComicSpider](db.DBComic, spilitChapterIdArr, "spider_end_status", 0)
		if err != nil {
			log.Errorf("func=%v, 第%v个数组, 爬取前查询是否需要爬 失败, reason: %v", funcName, i, err)
			c.JSON(500, fmt.Sprint("爬取失败, 爬取前查询是否需要爬 失败, 出错数组=", spilitChapterIdArr))
			return
		}

		// 如果查询数据库空的，跳过
		if len(spilitChapterIdArrNeedCrawl) == 0 {
			log.Warnf("func=%v, 第%v个数组, 不需要爬, 原因: 已经是爬取过的/无此id, arr= %v", funcName, i, spilitChapterIdArr)
			continue
		}

		// 步骤5.3.2 爬取并插入db --
		log.Debugf("func=%v, 要爬的chapterIdArr = %v, 过滤需要爬的chapterIdArr = %v", funcName, chapterIdArr, spilitChapterIdArrNeedCrawl)
		okTotalOneArr, err := SpiderManyChapterAllContent2DB(websiteId, websiteName, spilitChapterIdArrNeedCrawl) // 成功条数
		if err != nil {
			log.Errorf("func=%v, 第%v个数组, 爬取失败, 可能原因:1 爬取url不对 2 目标网站挂了 3 爬取逻辑错了,没爬到. err = %v", funcName, i, err)
			c.JSON(500, fmt.Sprint("爬取失败, 出错数组=", spilitChapterIdArr))
			return
		}
		log.Infof("func=%v, 第%v个数组, 爬取成功,total: %v, arr= %v", funcName, i, okTotalOneArr, spilitChapterIdArr)
		okTotal += okTotalOneArr
	}

	// 6. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条chapterContent数据")
}
