/*
*
爬取核心处理逻辑
*/
package spider

import (
	"fmt"
	"strings"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
)

// -- 初始化 ------------------------------------------------------------------------------
// -- 批量更新用到
// comic 表
var tableComicUniqueIndexArr = []string{"Name", "CountryId", "WebsiteId", "pornTypeId", "TypeId"} // 唯一索引字段
var tableComicUpdateColArr = []string{"latest_chapter", "hits", "comic_url_api_path", "cover_url_api_path", "brief_short", "brief_long", "end",
	"star", "spider_end_status", "download_end_status", "upload_aws_end_status", "upload_baidu_end_status", "release_date",
	"updated_at",
	"website_id", "porn_type_id", "country_id", "type_id", "process_id"} // 要更新的字段。要传updated_at ，upsert必须传, UPDATE()方法不用传，会自动改

// -- 爬漫画用 mapping
// 表映射，爬 https:/www.toptoon.net (台湾服务器) 用，爬的JSON数据
var ComicMappingForSpiderToptoonByJSON = map[string]models.ModelMapping{
	"name":          {GetFieldPath: "adult.%d.meta.title", FiledType: "string"}, // adult.100.meta.title 这样能获取第100个 的内容
	"websiteId":     {GetFieldPath: "websiteId", FiledType: "int"},
	"pornTypeId":    {GetFieldPath: "pornTypeId", FiledType: "int"},
	"countryId":     {GetFieldPath: "countryId", FiledType: "int"},
	"typeId":        {GetFieldPath: "typeId", FiledType: "int"},
	"processId":     {GetFieldPath: "processId", FiledType: "int"},
	"latestChapter": {GetFieldPath: "adult.%d.lastUpdated.episodeTitle", FiledType: "string"},
	"hits":          {GetFieldPath: "adult.%d.meta.viewCount", FiledType: "int"},
	"comicUrlApiPath": {GetFieldPath: "adult.%d.id", FiledType: "string",
		Transform: func(v any) any {
			id := v.(string)
			return "/comic/epList/" + id
		}}, // Template 表示模板：能实现拼接"/comic/epList/" + id
	"coverUrlApiPath": {GetFieldPath: "adult.%d.thumbnail.standard", FiledType: "string"},
	"briefLong":       {GetFieldPath: "adult.%d.meta.description", FiledType: "string"},
	"star":            {GetFieldPath: "adult.%d.meta.rating", FiledType: "float"},
	"releaseDate":     {GetFieldPath: "adult.%d.lastUpdated.pubDate", FiledType: "time"},
}

// -- 初始化 ------------------------------------------- end -----------------------------------

// -- 方法 ------------------------------------------------------------------------------
// mapping赋值。把 带%d 的mapping (内容不固定)，给%d赋值
/*
参数：
	1. mapping map[string]any  // 带%d 的mapping (内容不固定)，给%d赋值
	1. index int  // 要赋的值

返回：
	1. mapping map[string]any  // 赋完值的mapping

作用简单说：

作用详细说:

核心思路:
	1.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑
	5. 返回结果

注意：

使用方式：
*/
func mappingAssign(mapping map[string]models.ModelMapping, index int) map[string]models.ModelMapping {
	// 1. 校验传参
	// -- 至少校验空
	if mapping == nil {
		log.Error("func=mappingAssign(给带的mapping赋值). mapping不能为空")
		return nil
	}

	// 2. 数据清洗
	// 2. 数据清洗
	// 3. 业务逻辑 需要的数据校验 +清洗

	// 4. 执行核心逻辑
	for k, v := range mapping {
		// ！！！重要。只会带%号的处理，因为如果所有都 赋值，会报错
		if strings.Contains(v.GetFieldPath, "%d") {
			v.GetFieldPath = fmt.Sprintf(v.GetFieldPath, index) // 替换 %d
		}
		mapping[k] = v
	}
	// log.Debug("-------------------- func=mappingAssign(给带的mapping赋值). 赋值后 mapping: ", mapping["name"].GetFieldPath)  // 这样写不通用，仅仅适用于 comic表
	// log.Debug("-------------------- func=mappingAssign(给带的mapping赋值). 赋值后 mapping index=: ", index)

	// 5. 返回结果
	return mapping
}

// -- 方法 ------------------------------------------- end -----------------------------------
