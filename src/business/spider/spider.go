/*
*
爬取核心处理逻辑
*/
package spider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"study-spider-manhua-gin/src/util/langutil"
	"study-spider-manhua-gin/src/util/stringutil"
)

// -- 初始化 ------------------------------------------------------------------------------
// -- 接口
// 定义一个接口-业务数据清理，约定所有 要清理业务数据 的对象都要实现这个方法
type BusinessDataCleanable interface {
	BusinessDataClean() // 业务数据清理
}

// 定义一个统一：数据清理入口-包含 TrimSpaces() Trad2Simple() BusinessDataClean() 等方法啊
type DataCleanable interface {
	/*
		DataClean 真实逻辑，包含以下内容：
		TrimSpaces()        // 去除空格
		Trad2Simple()       // 繁体转简体
		BusinessDataClean() // 业务数据清理
	*/
	DataClean() // 数据清理
}

// -- 批量更新用到
// comic 表 --
var tableComicUniqueIndexArr = []string{"Name", "CountryId", "WebsiteId", "pornTypeId", "TypeId", "authorConcat"} // 唯一索引字段
var tableComicUpdateColArr = []string{"comic_url_api_path", "cover_url_api_path", "brief_short", "brief_long", "end",
	"spider_end_status", "download_end_status", "upload_aws_end_status", "upload_baidu_end_status", "release_date",
	"updated_at",
	"website_id", "porn_type_id", "country_id", "type_id", "process_id",
	"author_concat", "author_concat_type"} // 要更新的字段。要传updated_at ，upsert必须传, UPDATE()方法不用传，会自动改
// author 表 --
var tableAuthorUniqueIndexArr = []string{"Id"} // 唯一索引字段,用 models 里 字段
var tableAuthorUpdateColArr = []string{"name"} // 要更新的字段。要传updated_at ，upsert必须传, UPDATE()方法不用传，会自动改

// -- 爬漫画用 mapping
// 表映射，爬 https:/www.toptoon.net (台湾服务器) 用，爬的JSON数据 - 只能爬1个
var ComicMappingForSpiderToptoonByJSON = map[string]models.ModelMapping{
	"name":       {GetFieldPath: "adult.%d.meta.title", FiledType: "string"}, // adult.100.meta.title 这样能获取第100个 的内容
	"websiteId":  {GetFieldPath: "websiteId", FiledType: "int"},
	"pornTypeId": {GetFieldPath: "pornTypeId", FiledType: "int"},
	"countryId":  {GetFieldPath: "countryId", FiledType: "int"},
	"typeId":     {GetFieldPath: "typeId", FiledType: "int"},
	"processId":  {GetFieldPath: "processId", FiledType: "int"},
	"comicUrlApiPath": {GetFieldPath: "adult.%d.id", FiledType: "string",
		Transform: func(v any) any {
			id := v.(string)
			return "/comic/epList/" + id
		}}, // Template 表示模板：能实现拼接"/comic/epList/" + id
	"coverUrlApiPath":  {GetFieldPath: "adult.%d.thumbnail.standard", FiledType: "string"},
	"briefLong":        {GetFieldPath: "adult.%d.meta.description", FiledType: "string"},
	"releaseDate":      {GetFieldPath: "adult.%d.lastUpdated.pubDate", FiledType: "time"},
	"authorConcat":     {GetFieldPath: "adult.%d.meta.author.authorString", FiledType: "string"},
	"authorConcatType": {GetFieldPath: "authorConcatType", FiledType: "int"},
	"authorArr":        {GetFieldPath: "adult.%d.meta.author.authorData", FiledType: "array"}, // []any 表示数组

	"stats.latestChapterName":         {GetFieldPath: "adult.%d.lastUpdated.episodeTitle", FiledType: "string"},
	"stats.hits":                      {GetFieldPath: "adult.%d.meta.viewCount", FiledType: "int"},
	"stats.star":                      {GetFieldPath: "adult.%d.meta.rating", FiledType: "float"},
	"stats.totalChapter":              {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},
	"stats.lastestChapterReleaseDate": {GetFieldPath: "adult.%d.lastUpdated.pubDate", FiledType: "time"},
}

// 表映射，爬 https:/www.toptoon.net (台湾服务器) 用，爬的 Html 数据 - 只能爬1个
var ComicMappingForSpiderToptoonByHtml = map[string]models.ModelHtmlMapping{
	// 注意: GetFieldPath 假如要传2个值，就用|分隔. 比如: GetFieldPath: ".product__item__pic set-bg|onclick"
	// 表都有哪些数据
	"name": {GetFieldPath: ".product__item__text", GetHtmlType: "content", FiledType: "string",
		Transform: func(v any) string {
			// 爬出来 = 舞蹈系学姊们
			// 思路： 爬出来都是string类型，必须先清洗: 去空格，繁体转简体; 再做其他转换
			// 1. 去空格
			name := strings.TrimSpace(v.(string))
			// 2. 繁体转简体
			name, _ = langutil.TraditionalToSimplified(name)

			// 3. 返回
			return name
		}},

	// 外键相关
	"websiteId":       {GetFieldPath: "websiteId", GetHtmlType: "只能人工给", FiledType: "int"},       // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	"pornTypeId":      {GetFieldPath: "pornTypeId", GetHtmlType: "只能人工给", FiledType: "int"},      // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	"countryId":       {GetFieldPath: "countryId", GetHtmlType: "只能人工给", FiledType: "int"},       // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	"typeId":          {GetFieldPath: "typeId", GetHtmlType: "只能人工给", FiledType: "int"},          // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	"processId":       {GetFieldPath: "processId", GetHtmlType: "只能人工给", FiledType: "int"},       // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	"latestChapterId": {GetFieldPath: "latestChapterId", GetHtmlType: "只能人工给", FiledType: "int"}, // 没想好怎么获取, 可能要注释掉, 根据前端传参给赋值
	// "authorArr":       {GetFieldPath: "adult.%d.meta.author.authorData", FiledType: "array"}, // []any 表示数组 // 爬不到-----------

	// 其他
	"comicUrlApiPath": {GetFieldPath: ".product__item__pic.set-bg|onclick", GetHtmlType: "attr", FiledType: "string",
		Transform: func(v any) string {
			// 爬出来 = location.href='/manga/2722';
			// 思路： 爬出来都是string类型，必须先清洗: 去空格，繁体转简体; 再做其他转换
			// 1. 去空格
			value := strings.TrimSpace(v.(string))
			// 2. 繁体转简体
			value, _ = langutil.TraditionalToSimplified(value)

			// 3. 提取 location.href=' 引号中内容，re 正则获取
			re := regexp.MustCompile(`'(.+?)'`)
			value = re.FindStringSubmatch(value)[1]

			// 4. 返回
			return value
		}}, // Template 表示模板：能实现拼接"/comic/epList/" + id ->>>>>>>>>>>>>>>>>>>>> 好像不对，还需要 location.href='/manga/4015'; 去除些内容
	"coverUrlApiPath": {GetFieldPath: ".product__item__pic.set-bg|data-setbg", GetHtmlType: "attr", FiledType: "string",
		Transform: func(v any) string {
			// 爬出来 = https://img.imh99.top/webtoon/cover-image/618_1748423944636.webp
			// 思路： 爬出来都是string类型，必须先清洗: 去空格，繁体转简体; 再做其他转换
			// 1. 去空格
			value := strings.TrimSpace(v.(string))
			// 2. 繁体转简体
			value, _ = langutil.TraditionalToSimplified(value)

			// 3. 去除协议头+域名，https://img.imh99.to,只留后面内容,re 正则获取
			re := regexp.MustCompile(`https://img.imh99.to`)
			value = re.ReplaceAllString(value, "")

			// 4. 返回
			return value
		}}, // 还需要方法，去除一些东西

	// "briefShort":           {GetFieldPath: "adult.%d.meta.description", FiledType: "string"}, // 爬不到-----------
	// "briefLong":            {GetFieldPath: "adult.%d.meta.description", FiledType: "string"}, // 爬不到-----------
	"end": {GetFieldPath: ".epgreen", GetHtmlType: "content", FiledType: "int",
		Transform: func(v any) string {
			// 爬出来 = 完结 / 连载
			// 思路： 爬出来都是string类型，必须先清洗: 去空格，繁体转简体; 再做其他转换
			// 1. 去空格
			value := strings.TrimSpace(v.(string))
			// 2. 繁体转简体
			value, _ = langutil.TraditionalToSimplified(value)

			// 3. 把 中文的结束状态 转成 数字-》对应数据库中  未知1 连载2 完结3
			switch value {
			case "完结":
				value = "3"
			case "连载":
				value = "2"
			default:
				value = "1"
			}

			// 4. 返回
			return value
		}},
	// "spiderEndStatus":      {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},             // 爬不到
	// "downloadEndStatus":    {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},             // 爬不到
	// "uploadAwsEndStatus":   {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},             // 爬不到
	// "uploadBaiduEndStatus": {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},             // 爬不到
	// "releaseDate":          {GetFieldPath: "adult.%d.lastUpdated.pubDate", FiledType: "time"},        // 爬不到-----------
	// "authorConcat":         {GetFieldPath: "adult.%d.meta.author.authorString", FiledType: "string"}, // 爬不到-----------
	// "authorConcatType":     {GetFieldPath: "authorConcatType", FiledType: "int"},                     // 爬不到-----------

	// 子表相关
	// "stats.latestChapterName":         {GetFieldPath: "adult.%d.lastUpdated.episodeTitle", FiledType: "string"}, // 爬不到-----------
	"stats.hits": {GetFieldPath: ".view", GetHtmlType: "content", FiledType: "int",
		Transform: func(v any) string {
			// 爬出来 = 5.3w
			// 思路： 爬出来都是string类型，必须先清洗: 去空格，繁体转简体; 再做其他转换
			// 1. 去空格
			value := strings.TrimSpace(v.(string))
			// 2. 繁体转简体
			value, _ = langutil.TraditionalToSimplified(value)

			// 3. 把 带字母/中文的 “访问量” 转成 数字, 判断逻辑: 看字符末尾是由有 k/千w/万, 正则实现
			value = strconv.Itoa(stringutil.ParseHitsStr(value)) // stringutil.ParseHitsStr(value) 返回的是int

			// 4. 返回
			return value
		}},
	// "stats.star":                      {GetFieldPath: "adult.%d.meta.rating", FiledType: "float"},        // 爬不到-----------
	// "stats.totalChapter":              {GetFieldPath: "adult.%d.meta.epTotalCnt", FiledType: "int"},      // 爬不到-----------
	// "stats.lastestChapterReleaseDate": {GetFieldPath: "adult.%d.lastUpdated.pubDate", FiledType: "time"}, // 爬不到-----------
}

// 表映射，爬 https:/www.toptoon.net (台湾服务器)  - 作者相关 用，爬的JSON数据
var AuthorMappingForSpiderToptoonByJSON = map[string]models.ModelMapping{
	"name": {GetFieldPath: "adult.%d.meta.author.authorData.%d.name", FiledType: "string"}, // 参考 /doc/F12找到的JSON/comic项目/类别/任一json
}

// -- 爬漫画章节用 mapping
var ComicChapterMappingForSpiderToptoonByJSON = map[string]models.ModelMapping{
	// ---- 下面都是copy 的 comic的，还要改 ！！！！！！！！！！！！！！
	// content 表示内容, 爬的时候用 element.ChildText(".comic__title")
	"chapterNum":         {GetFieldPath: "adult.%d.meta.title", FiledType: "content"},
	"chapterSubNum":      {GetFieldPath: "websiteId", FiledType: "int"},
	"chapterRealSortNum": {GetFieldPath: "pornTypeId", FiledType: "int"},
	"name":               {GetFieldPath: ".comic__title", FiledType: "content"}, // content 表示内容，不转换
	"urlApiPath":         {GetFieldPath: "typeId", FiledType: "string"},
	"releaseDate":        {GetFieldPath: "processId", FiledType: "time"},
	"SpiderStatus":       {GetFieldPath: "adult.%d.thumbnail.standard", FiledType: "int"},

	// 下面是参考 - 一会删
	"不知道": {GetFieldPath: "adult.%d.id", FiledType: "string",
		Transform: func(v any) any {
			id := v.(string)
			return "/comic/epList/" + id
		}}, // Template 表示模板：能实现拼接"/comic/epList/" + id
}

func init() {
	// 为防止控制台警告，看着烦，临时写点日志打印。后续可以删除
	fmt.Println("-------------------- func=init() . 为防止控制台警告，看着烦，临时打印", tableAuthorUniqueIndexArr)
	fmt.Println("-------------------- func=init() . 为防止控制台警告，看着烦，临时打印", tableAuthorUpdateColArr)
}

// -- 初始化 ------------------------------------------- end -----------------------------------

// -- 方法 ------------------------------------------------------------------------------
// mapping赋值。把 带%d 的mapping (内容不固定)，给%d赋值
/*
参数：
	1. mapping map[string]any  // 带%d 的mapping (内容不固定)，给%d赋值
	1. indices ...int  // 要赋的值，支持多个值

返回：
	1. mapping map[string]any  // 赋完值的mapping

作用简单说：

作用详细说:

核心思路：
	1. 支持多个%d占位符的替换

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑
	5. 返回结果

注意：

使用方式：
	- 对于单个%d占位符：mappingAssign(mapping, index)
	- 对于多个%d占位符：mappingAssign(mapping, index1, index2, ...)
*/
func mappingAssign(mapping map[string]models.ModelMapping, indices ...int) map[string]models.ModelMapping {
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
			// 根据传入的参数数量，使用不同的fmt.Sprintf调用方式
			if len(indices) == 1 {
				v.GetFieldPath = fmt.Sprintf(v.GetFieldPath, indices[0]) // 单个参数
			} else if len(indices) == 2 {
				v.GetFieldPath = fmt.Sprintf(v.GetFieldPath, indices[0], indices[1]) // 两个参数
			} else if len(indices) == 3 {
				v.GetFieldPath = fmt.Sprintf(v.GetFieldPath, indices[0], indices[1], indices[2]) // 三个参数
			} else {
				// 对于更多参数的情况，使用反射来动态处理
				args := make([]interface{}, len(indices))
				for i, idx := range indices {
					args[i] = idx
				}
				v.GetFieldPath = fmt.Sprintf(v.GetFieldPath, args...)
			}
		}
		mapping[k] = v
	}
	// log.Debug("-------------------- func=mappingAssign(给带的mapping赋值). 赋值后 mapping: ", mapping["name"].GetFieldPath)  // 这样写不通用，仅仅适用于 comic表
	// log.Debug("-------------------- func=mappingAssign(给带的mapping赋值). 赋值后 mapping index=: ", index)

	// 5. 返回结果
	return mapping
}

// 实现接口。处理 models对象，清洗 业务数据
func BusinessDataCleanObj(obj BusinessDataCleanable) {
	obj.BusinessDataClean()
}

// 实现接口: 数据清理，统一入口。处理 models 对象
func DataCleanObj(obj DataCleanable) {
	obj.DataClean() // 业务数据清理
}

// -- 方法 ------------------------------------------- end -----------------------------------
