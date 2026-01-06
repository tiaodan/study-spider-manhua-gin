/*
功能: 处理爬取api, V1.5版本实现方式
*/

package spider

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

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
		log.Info("----- delete spiderUrl = ", spiderUrlArr[i])
	}

	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	switch website {
	case "toptoon-tw":
		log.Info("------ 还没实现")
	case "kxmanhua": // 开心看漫画
		// 0. 使用 kxmanhua 相关配置
		webCfg := config.CfgSpiderYaml.Websites["kxmanhua"]
		if webCfg == nil {
			c.JSON(400, gin.H{"error": "func=爬oneTypeAllBookArr V1.5, 配置文件里没有找到网站 kxmanhua 的配置"}) // 返回错误
			return
		}
		log.Debug("------- webCfg = ", webCfg)

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

	default:
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeByJson), 没找到到应爬哪个网站. 建议: 排查json参数 spiderTag-website"}) // 返回错误
	}

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条数据")

}

// 爬某一本书所有章节 - V1.5版本实现方式
func DispatchApi_OneBookAllChapter_V1_5(c *gin.Context) {}

// 爬某一章节所有内容 - V1.5版本实现方式
func DispatchApi_OneChapterAllContent_V1_5(c *gin.Context) {}
