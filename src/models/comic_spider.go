// 拼多多订单数据模型, 存数据用的
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
	WebsiteId  int `json:"websiteId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"websiteId" `   // 网站id-外键 组合索引字段
	PornTypeId int `json:"pornTypeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"pornTypeId" ` // 总分类id-最高级-外键 组合索引字段
	CountryId  int `json:"countryId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"countryId" `   // 国家id-外键 组合索引字段
	TypeId     int `json:"typeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"typeId" `         // 类型id-外键 组合索引字段
	ProcessId  int `json:"process" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"ProcessId" `     // 进度id-外键 组合索引字段

	// 其它
	LatestChapter        string    `json:"latestChapter" gorm:"not null" spider:"latestChapter" `                                    // 更新到多少集, 字符串,最新章节.可以是空字符串
	Hits                 int       `json:"hits" gorm:"not null" spider:"hits" `                                                      // 人气
	ComicUrlApiPath      string    `json:"comicUrlApiPath" gorm:"not null;check:comic_url_api_path <> ''" spider:"comicUrlApiPath" ` // 漫画链接.不能是空字符串
	CoverUrlApiPath      string    `json:"coverUrlApiPath" gorm:"not null;check:cover_url_api_path <> ''" spider:"coverUrlApiPath" ` // 封面链接.不能是空字符串
	BriefShort           string    `json:"briefShort" gorm:"not null" spider:"briefShort" `                                          // 简介-短.可以是空字符串
	BriefLong            string    `json:"briefLong" gorm:"not null" spider:"briefLong" `                                            // 简介-长.可以是空字符串
	End                  bool      `json:"end" gorm:"not null" spider:"end" `                                                        // 漫画是否完结,如果完结是1
	Star                 float64   `json:"star" gorm:"not null" spider:"star" `                                                      // 评分
	SpiderEndStatus      int       `json:"spiderEndStatus" gorm:"not null" spider:"spiderEndStatus" `                                // 爬取结束
	DownloadEndStatus    int       `json:"downloadEndStatus" gorm:"not null" spider:"downloadEndStatus" `                            // 下载结束
	UploadAwsEndStatus   int       `json:"uploadAwsEndStatus" gorm:"not null" spider:"uploadAwsEndStatus" `                          // 是否上传到aws
	UploadBaiduEndStatus int       `json:"uploadBaiduEndStatus" gorm:"not null" spider:"uploadBaiduEndStatus" `                      // 是否上传到baidu网盘
	ReleaseDate          time.Time `json:"releaseDate" gorm:"not null" spider:"releaseDate" `                                        // 发布日期.可以是空字符串

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除

	// 关联外键写法，更新时，同步更新，删除时，不让删
	Country  Country  `gorm:"foreignKey:CountryId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 可选：级联操作
	Website  Website  `gorm:"foreignKey:WebsiteId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PornType PornType `gorm:"foreignKey:PornTypeId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Type     Type     `gorm:"foreignKey:TypeId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Process  Process  `gorm:"foreignKey:TypeId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

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

// 实现stringutils 里 处理空格接口
func (c *ComicSpider) TrimSpaces() {
	// 只要是stirng类型，就去掉前后空格
	c.Name = strings.TrimSpace(c.Name)
	c.LatestChapter = strings.TrimSpace(c.LatestChapter)
	c.ComicUrlApiPath = strings.TrimSpace(c.ComicUrlApiPath)
	c.CoverUrlApiPath = strings.TrimSpace(c.CoverUrlApiPath)
	c.BriefShort = strings.TrimSpace(c.BriefShort)
	c.BriefLong = strings.TrimSpace(c.BriefLong)
}

// 实现stringutils 里 繁体转简体 接口
func (c *ComicSpider) Trad2Simple() {
	// 只要是string类型，都处理
	c.Name, _ = langutil.TraditionalToSimplified(c.Name)
	c.LatestChapter, _ = langutil.TraditionalToSimplified(c.LatestChapter)
	c.ComicUrlApiPath, _ = langutil.TraditionalToSimplified(c.ComicUrlApiPath)
	c.CoverUrlApiPath, _ = langutil.TraditionalToSimplified(c.CoverUrlApiPath)
	c.BriefShort, _ = langutil.TraditionalToSimplified(c.BriefShort)
	c.BriefLong, _ = langutil.TraditionalToSimplified(c.BriefLong)
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
