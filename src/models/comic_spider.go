// ComicSpider 数据模型, 存数据用的
package models

import (
	"strings"
	"study-spider-manhua-gin/src/util/langutil"
	"time"

	"gorm.io/gorm"
)

// 漫画数据 - 存爬取数据用的
// 分类顺序 国家-> 网站 -> 总分类 -> 类型
// 不用写column,因为系统会自动关联
type ComicSpider struct {
	// 不写column写法

	// 本来用的uint，改成int了。因为1) uint有风险，传负数-》变成很大的正数。 2) 默认建数据库表也不用uint，用的int
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`                                                          // 数据库id,主键、自增.
	Name string `json:"name" gorm:"not null; uniqueIndex:idx_comic_unique;size:150;check:name <> ''" spider:"name" ` // 漫画名 组合索引字段
	// 外键
	WebsiteId  int      `json:"websiteId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"websiteId" `   // 网站id-外键 组合索引字段
	PornTypeId int      `json:"pornTypeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"pornTypeId" ` // 总分类id-最高级-外键 组合索引字段
	CountryId  int      `json:"countryId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"countryId" `   // 国家id-外键 组合索引字段
	TypeId     int      `json:"typeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"typeId" `         // 类型id-外键 组合索引字段
	ProcessId  int      `json:"process" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"ProcessId" `     // 进度id-外键 组合索引字段
	AuthorArr  []Author `gorm:"many2many:comic_spider_authors;" spider:"authorArr" `                           // 多对多关联

	// 其它
	ComicUrlApiPath      string    `json:"comicUrlApiPath" gorm:"not null;check:comic_url_api_path <> ''" spider:"comicUrlApiPath" `                             // 漫画链接.不能是空字符串
	CoverUrlApiPath      string    `json:"coverUrlApiPath" gorm:"not null;check:cover_url_api_path <> ''" spider:"coverUrlApiPath" `                             // 封面链接.不能是空字符串
	BriefShort           string    `json:"briefShort" gorm:"not null" spider:"briefShort" `                                                                      // 简介-短.可以是空字符串
	BriefLong            string    `json:"briefLong" gorm:"not null" spider:"briefLong" `                                                                        // 简介-长.可以是空字符串
	End                  bool      `json:"end" gorm:"not null" spider:"end" `                                                                                    // 漫画是否完结,如果完结是1
	SpiderEndStatus      int       `json:"spiderEndStatus" gorm:"not null" spider:"spiderEndStatus" `                                                            // 爬取结束
	DownloadEndStatus    int       `json:"downloadEndStatus" gorm:"not null" spider:"downloadEndStatus" `                                                        // 下载结束
	UploadAwsEndStatus   int       `json:"uploadAwsEndStatus" gorm:"not null" spider:"uploadAwsEndStatus" `                                                      // 是否上传到aws
	UploadBaiduEndStatus int       `json:"uploadBaiduEndStatus" gorm:"not null" spider:"uploadBaiduEndStatus" `                                                  // 是否上传到baidu网盘
	ReleaseDate          time.Time `json:"releaseDate" gorm:"not null" spider:"releaseDate" `                                                                    // 发布日期.可以是空字符串
	AuthorConcat         string    `json:"authorConcat" gorm:"not null;uniqueIndex:idx_comic_unique;size:500; check:author_concat <> ''" spider:"authorConcat" ` // 作者.不能是空字符串。组合索引
	AuthorConcatType     int       `json:"authorConcatType" gorm:"not null" spider:"authorConcatType" `                                                          // 作者拼接方式，不能空。：0 默认，按爬取顺序拼接，1: 按字母升序拼接 2:按我的意愿拼接 3: 参考最权威的网站拼接(b比如有声书，参考喜马拉雅，韩漫参考toptoon，小说参考 起点-建议0 /3

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除

	// 关联外键写法，更新时，同步更新，删除时，不让删
	// 注意：references 写主表id, foreignKey 写从表id
	Country  Country          `gorm:"foreignKey:CountryId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 可选：级联操作
	Website  Website          `gorm:"foreignKey:WebsiteId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PornType PornType         `gorm:"foreignKey:PornTypeId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Type     Type             `gorm:"foreignKey:TypeId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Process  Process          `gorm:"foreignKey:ProcessId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Stats    ComicSpiderStats `gorm:"foreignKey:ComicID;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 漫画 统计

	// 写column写法
	/*
		Id   uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`                       // 数据库id,主键、自增
		Name string `json:"name" gorm:"not null; uniqueIndex:name_unique;size:150;column:name"` // 漫画名 数据库唯一索引
		// 外键
		CountryId  int `json:"countryId" gorm:"not null;column:country_id"`   // 国家id-外键
		WebsiteId  int `json:"websiteId" gorm:"not null;column:website_id"`   // 网站id-外键
		PornTypeId int `json:"pornTypeId" gorm:"not null;column:porn_type_id"` // 总分类id-最高级-外键
		TypeId     int `json:"typeId" gorm:"not null;column:type_id"`         // 类型id-外键

		// 其它
		Update         string  `json:"update" gorm:"not null;column:update"`                   // 更新到多少集, 字符串
		Hits           uint    `json:"hits" gorm:"not null;hits"`                              // 人气
		ComicUrlApiPath       string  `json:"comicUrlApiPath" gorm:"not null;column:comic_url_api_path"`              // 漫画链接
		CoverUrlApiPath       string  `json:"coverUrlApiPath" gorm:"not null;column:cover_url_api_path"`              // 封面链接
		BriefShort     string  `json:"briefShort" gorm:"not null;column:brief_short"`          // 简介-短
		BriefLong      string  `json:"briefLong" gorm:"not null;column:brief_long"`            // 简介-长
		End            int     `json:"end" gorm:"not null;column:end"`                         // 漫画是否完结,如果完结是1
		Star           float64 `json:"star" gorm:"not null;column:star"`                       // 评分
		NeedTcp        int     `json:"needTcp" gorm:"not null;column:need_tcp"`                // 漫画是否需要http 或者https前缀,如果链接有了tcp,就应该是0
		CoverNeedTcp   int     `json:"coverNeedTcp" gorm:"not null;column:cover_need_tcp"`     // 封面是否需要http 或者https前缀,如果链接有了tcp,就应该是0
		SpiderEnd      int     `json:"spiderEnd" gorm:"not null;column:spider_end"`            // 爬取结束
		DownloadEnd    int     `json:"downloadEnd" gorm:"not null;column:download_end"`        // 下载结束
		UploadAwsEnd   int     `json:"uploadAwsEnd" gorm:"not null;column:upload_aws_end"`     // 是否上传到aws
		UploadBaiduEnd int     `json:"uploadBaiduEnd" gorm:"not null;column:upload_baidu_end"` // 是否上传到baidu网盘

		// 关联外键写法
		Country  Country  `gorm:"foreignKey:CountryId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 可选：级联操作
		Website  Website  `gorm:"foreignKey:WebsiteId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
		PornType PornType `gorm:"foreignKey:PornTypeId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
		Type     Type     `gorm:"foreignKey:TypeId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	*/
}

// 统计状态。频繁更新的字段放这个表里，相当于 comic 子表，通过外键关联上
type ComicSpiderStats struct {
	ID      int `json:"id" gorm:"primaryKey;autoIncrement"`
	ComicID int `gorm:"uniqueIndex"` // 外键，唯一索引保证一对一

	// 频繁更新字段
	Star          float64 `json:"star" gorm:"not null" spider:"star" `                   // 评分
	LatestChapter string  `json:"latestChapter" gorm:"not null" spider:"latestChapter" ` // 更新到多少集, 字符串,最新章节.可以是空字符串
	Hits          int     `json:"hits" gorm:"not null" spider:"hits" `                   // 人气
}

// 实现stringutils 里 处理空格接口
func (c *ComicSpider) TrimSpaces() {
	// 只要是stirng类型，就去掉前后空格
	c.Name = strings.TrimSpace(c.Name)
	c.ComicUrlApiPath = strings.TrimSpace(c.ComicUrlApiPath)
	c.CoverUrlApiPath = strings.TrimSpace(c.CoverUrlApiPath)
	c.BriefShort = strings.TrimSpace(c.BriefShort)
	c.BriefLong = strings.TrimSpace(c.BriefLong)
	c.AuthorConcat = strings.TrimSpace(c.AuthorConcat)
}

// 实现stringutils 里 处理空格接口 - 处理子表 -统计数据
func (c *ComicSpiderStats) TrimSpaces() {
	c.LatestChapter = strings.TrimSpace(c.LatestChapter)

}

// 实现stringutils 里 繁体转简体 接口
func (c *ComicSpider) Trad2Simple() {
	// 只要是string类型，都处理
	c.Name, _ = langutil.TraditionalToSimplified(c.Name)
	c.ComicUrlApiPath, _ = langutil.TraditionalToSimplified(c.ComicUrlApiPath)
	c.CoverUrlApiPath, _ = langutil.TraditionalToSimplified(c.CoverUrlApiPath)
	c.BriefShort, _ = langutil.TraditionalToSimplified(c.BriefShort)
	c.BriefLong, _ = langutil.TraditionalToSimplified(c.BriefLong)
	c.AuthorConcat, _ = langutil.TraditionalToSimplified(c.AuthorConcat)
}

// 实现stringutils 里 繁体转简体 接口 - 处理子表 -统计数据
func (c *ComicSpiderStats) Trad2Simple() {
	c.LatestChapter, _ = langutil.TraditionalToSimplified(c.LatestChapter)
}

// 表字段的 ”爬取“映射关系 结构，写通用爬虫方法时，只要实现这个结构，就能用通用爬虫方法爬取数据
type ComicSpiderFieldMapping struct {
	// 下面2个 是相互对应关系。比如： JsonFieldName = name, ModelFiledType="string"
	GetFieldPath string              // gjson 提取字段路径。获取这个字段, gjson path。提取的是json数据里的字段
	FiledType    string              // 字段类型。如 "string","int“,”float“,”array“ ..
	Transform    func(value any) any // 转换函数。提取到字段后，转换成数据库字段类型

	// 使用方式：举例
	/*
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
}
