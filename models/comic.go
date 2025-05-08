// 拼多多订单数据模型, 存数据用的
package models

import (
	"time"

	"gorm.io/gorm"
)

// 漫画数据
// 分类顺序 国家-> 网站 -> 总分类 -> 类型
// 不用写column,因为系统会自动关联
type Comic struct {
	// 不写column写法
	Id   uint   `json:"id" gorm:"primaryKey;autoIncrement"`                          // 数据库id,主键、自增
	Name string `json:"name" gorm:"not null; uniqueIndex:idx_comic_unique;size:150"` // 漫画名 组合索引字段
	// 外键
	CountryId  int `json:"countryId" gorm:"not null; uniqueIndex:idx_comic_unique" `  // 国家id-外键 组合索引字段
	WebsiteId  int `json:"websiteId" gorm:"not null; uniqueIndex:idx_comic_unique" `  // 网站id-外键 组合索引字段
	CategoryId int `json:"categoryId" gorm:"not null; uniqueIndex:idx_comic_unique" ` // 总分类id-最高级-外键 组合索引字段
	TypeId     int `json:"typeId" gorm:"not null; uniqueIndex:idx_comic_unique" `     // 类型id-外键 组合索引字段

	// 其它
	Update         string  `json:"update" gorm:"not null" `         // 更新到多少集, 字符串
	Hits           uint    `json:"hits" gorm:"not null" `           // 人气
	ComicUrl       string  `json:"comicUrl" gorm:"not null"`        // 漫画链接
	CoverUrl       string  `json:"coverUrl" gorm:"not null" `       // 封面链接
	BriefShort     string  `json:"briefShort" gorm:"not null" `     // 简介-短
	BriefLong      string  `json:"briefLong" gorm:"not null" `      // 简介-长
	End            int     `json:"end" gorm:"not null" `            // 漫画是否完结,如果完结是1
	Star           float64 `json:"star" gorm:"not null"`            // 评分
	NeedTcp        int     `json:"needTcp" gorm:"not null" `        // 漫画是否需要http 或者https前缀,如果链接有了tcp,就应该是0
	CoverNeedTcp   int     `json:"coverNeedTcp" gorm:"not null" `   // 封面是否需要http 或者https前缀,如果链接有了tcp,就应该是0
	SpiderEnd      int     `json:"spiderEnd" gorm:"not null" `      // 爬取结束
	DownloadEnd    int     `json:"downloadEnd" gorm:"not null" `    // 下载结束
	UploadAwsEnd   int     `json:"uploadAwsEnd" gorm:"not null" `   // 是否上传到aws
	UploadBaiduEnd int     `json:"uploadBaiduEnd" gorm:"not null" ` // 是否上传到baidu网盘

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除

	// 关联外键写法
	Country  Country  `gorm:"foreignKey:CountryId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 可选：级联操作
	Website  Website  `gorm:"foreignKey:WebsiteId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Category Category `gorm:"foreignKey:CategoryId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Type     Type     `gorm:"foreignKey:TypeId;references:NameId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// 写column写法
	/*
		Id   uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`                       // 数据库id,主键、自增
		Name string `json:"name" gorm:"not null; uniqueIndex:name_unique;size:150;column:name"` // 漫画名 数据库唯一索引
		// 外键
		CountryId  int `json:"countryId" gorm:"not null;column:country_id"`   // 国家id-外键
		WebsiteId  int `json:"websiteId" gorm:"not null;column:website_id"`   // 网站id-外键
		CategoryId int `json:"categoryId" gorm:"not null;column:category_id"` // 总分类id-最高级-外键
		TypeId     int `json:"typeId" gorm:"not null;column:type_id"`         // 类型id-外键

		// 其它
		Update         string  `json:"update" gorm:"not null;column:update"`                   // 更新到多少集, 字符串
		Hits           uint    `json:"hits" gorm:"not null;hits"`                              // 人气
		ComicUrl       string  `json:"comicUrl" gorm:"not null;column:comic_url"`              // 漫画链接
		CoverUrl       string  `json:"coverUrl" gorm:"not null;column:cover_url"`              // 封面链接
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
		Category Category `gorm:"foreignKey:CategoryId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
		Type     Type     `gorm:"foreignKey:TypeId; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	*/
}
