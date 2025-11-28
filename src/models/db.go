package models

// 1. 放一些db 会用到的 通用的 struct

// 表字段的 ”爬取“映射关系 结构，写通用爬虫方法时，只要实现这个结构，就能用通用爬虫方法爬取数据
/*
作用简单说：
	- 爬取json数据，并映射到数据库字段

作用详细说:

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑 - 爬取 - 插入db
	5. 返回结果

参数：
	1.

返回：
	无

注意：

使用方式：
type FieldDef struct {
	Path string // gjson path
	Type string // string,int,float,array...
}

var BookFieldMap = map[string]FieldDef{
	"websiteId":   {Path: "websiteId", Type: "int"},
	"pornTypeId":  {Path: "pornTypeId", Type: "int"},
	"countryId":   {Path: "countryId", Type: "int"},
	"typeId":      {Path: "typeId", Type: "int"},
	"bookName":    {Path: "adult.100.meta.title", Type: "string"},
	"update":      {Path: "adult.100.lastUpdated.episodeTitle", Type: "string"},
	"hits":        {Path: "adult.100.meta.viewCount", Type: "int"},
	"comicUrlApiPath":    {Path: "adult.100.id", Type: "string"},
	"coverUrlApiPath":    {Path: "adult.100.thumbnail.standard", Type: "string"},
	"briefLong":   {Path: "adult.100.meta.description", Type: "string"},
	"star":        {Path: "adult.100.meta.rating", Type: "float"},
}

// 用的时候，写一个通用 字段提取方法
func ExtractFields(data []byte, fieldMap map[string]FieldDef) map[string]interface{} {
	result := make(map[string]interface{})

	for key, def := range fieldMap {
		v := gjson.GetBytes(data, def.Path)

		switch def.Type {
		case "string":
			result[key] = v.String()
		case "int":
			result[key] = v.Int()
		case "float":
			result[key] = v.Float()
		case "bool":
			result[key] = v.Bool()
		case "array":
			result[key] = v.Array()
		default:
			result[key] = v.Value() // fallback
		}
	}
	return result
}
*/
type ModelMapping struct {
	// 下面2个 是相互对应关系。比如： JsonFieldName = name, ModelFiledType="string"
	GetFieldPath string              // gjson 提取字段路径。获取这个字段, gjson path。提取的是json数据里的字段
	FiledType    string              // 字段类型。如 "string","int“,”float“,”array“ ..
	Transform    func(value any) any // 转换函数。提取到字段后，转换成数据库字段类型
}

// 定义一个泛型接口
type ModelIface interface {
	*Website | *Country | *PornType | *Type | *ComicSpider |
		*WebsiteType | *Process | *Author
}
