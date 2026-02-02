/*
功能：db 管理api接口 - website 表
*/
package db_manage_api

import (
	"fmt"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	"github.com/gin-gonic/gin"
)

// 增
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
func WebsiteAdd(c *gin.Context) {
	// 初始化
	funcName := "WebsiteAdd"

	//  1. 校验传参
	//  2. 数据清洗
	//  3. 业务逻辑 需要的数据校验 +清洗
	//  4. 执行核心逻辑
	// 	- 读取html内容
	// 	- 通过mapping 爬取字段，赋值给chapter_spider对象
	// 	- 验证业务逻辑，保证稳定性(比如 websiteId是否存在, countryId是否存在等)
	// 	- 插入前, 数据清洗
	// 	- 批量插入db
	var website models.Website
	if err := c.ShouldBindJSON(&website); err != nil {
		log.Error("绑定JSON失败: ", err)
		c.JSON(400, gin.H{"error": "绑定JSON失败"})
		return
	}
	website.DataClean() // 数据清理

	uniqueIndexArr := []string{"Name", "Domain"}
	updateDBColumnRealNameArr := []string{"need_proxy", "is_https", "is_refer", "cover_url_is_need_https", "chapter_content_url_is_need_https",
		"cover_url_concat_rule", "chapter_content_url_concat_rule", "cover_domain", "chapter_content_domain", "book_can_spider_type", "chapter_can_spider_type",
		"book_spider_req_body_eg_server_filepath", "chapter_spider_req_body_eg_server_filepath", "star_type", "website_type_id"}
	err := db.DBUpsert(db.DBComic, &website, uniqueIndexArr, updateDBColumnRealNameArr) // 先这么写，等有了配置，再用配置
	errorutil.ErrorPrint(err, "wegs8te Add 错误, err = ")
	log.Info("插入成功，判断websiteId 是否变化 = ", website.Id)

	// 4.2 回显id, 自动创建content 分表
	if website.Id == 0 {
		log.Infof("func=%v, websiteId 为 0, 说明执行的是更新操作", funcName)
		c.JSON(200, gin.H{"message": "成功", "data": website})
		return
	}

	// 推荐用 4 位补 0，方便排序、扩容、视觉对齐
	tableName := fmt.Sprintf("chapter_content_spider_%04d", website.Id)

	// 关键：强制指定表名再 AutoMigrate
	err = db.DBComic.Table(tableName).AutoMigrate(&models.ChapterContentSpider{})
	if err != nil {
		log.Errorf("分表 %s 迁移失败: %v", tableName, err)
		c.JSON(500, gin.H{"error": "分表迁移失败"})
		return
	}
	log.Infof("分表 %s 创建/迁移成功", tableName)

	//  5. 返回结果
	c.JSON(200, gin.H{"message": "成功", "data": website})
}

// 删

// 改
// 查
