// 拼多多订单数据模型, 存数据用的
package models

import (
	"strings"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/util/langutil"
	"study-spider-manhua-gin/src/util/stringutil"
	"time"

	"gorm.io/gorm"
)

// 漫画数据 - 我用的,就是业务真实 用的表
// 分类顺序 国家-> 网站 -> 总分类 -> 类型
// 不用写column,因为系统会自动关联
/*
和 ComicSpider 的区别是：
	1. 多 cover_save_path_api_path 字段 // 封面图片, 保存路径的api
*/
type ComicMy struct {
	// 不写column写法

	// 本来用的uint，改成int了。因为1) uint有风险，传负数-》变成很大的正数。 2) 默认建数据库表也不用uint，用的int
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`                                                          // 数据库id,主键、自增.
	Name string `json:"name" gorm:"not null; uniqueIndex:idx_comic_unique;size:150;check:name <> ''" spider:"name" ` // 漫画名 组合索引字段
	// 外键
	WebsiteId       int      `json:"websiteId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"websiteId" `   // 网站id-外键 组合索引字段
	PornTypeId      int      `json:"pornTypeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"pornTypeId" ` // 总分类id-最高级-外键 组合索引字段
	CountryId       int      `json:"countryId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"countryId" `   // 国家id-外键 组合索引字段
	TypeId          int      `json:"typeId" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"typeId" `         // 类型id-外键 组合索引字段
	ProcessId       int      `json:"process" gorm:"not null; uniqueIndex:idx_comic_unique" spider:"processId" `     // 进度id-外键 组合索引字段
	AuthorArr       []Author `gorm:"many2many:comic_my_authors;" spider:"authorArr" `                               // 多对多关联
	LatestChapterId *int     `json:"latestChapterId" spider:"latestChapterId" `                                     // 最新章节id。可为空，因为爬书的时候，章节表还没有内容。传指针，传nil时，就是null

	// 其它
	ComicUrlApiPath      string    `json:"comicUrlApiPath" gorm:"not null;check:comic_url_api_path <> ''" spider:"comicUrlApiPath" `                             // 漫画链接.不能是空字符串
	CoverUrlApiPath      string    `json:"coverUrlApiPath" gorm:"not null;check:cover_url_api_path <> ''" spider:"coverUrlApiPath" `                             // 封面链接.不能是空字符串
	CoverSavePathApiPath string    `json:"coverSavePathApiPath" gorm:"not null;check:cover_save_path_api_path <> ''" spider:"coverSavePathApiPath" `             // 封面图片, 保存路径的api. .可以是空字符串,因为没上传时，是空的
	BriefShort           string    `json:"briefShort" gorm:"not null" spider:"briefShort" `                                                                      // 简介-短.可以是空字符串
	BriefLong            string    `json:"briefLong" gorm:"not null" spider:"briefLong" `                                                                        // 简介-长.可以是空字符串
	End                  int       `json:"end" gorm:"not null" spider:"end" `                                                                                    // 漫画是否完结,如果 未知1 连载2 完结3 == processId
	SpiderEndStatus      int       `json:"spiderEndStatus" gorm:"not null" spider:"spiderEndStatus" `                                                            // 爬取结束状态
	DownloadEndStatus    int       `json:"downloadEndStatus" gorm:"not null" spider:"downloadEndStatus" `                                                        // 下载结束状态
	UploadAwsEndStatus   int       `json:"uploadAwsEndStatus" gorm:"not null" spider:"uploadAwsEndStatus" `                                                      // 是否上传到aws
	UploadBaiduEndStatus int       `json:"uploadBaiduEndStatus" gorm:"not null" spider:"uploadBaiduEndStatus" `                                                  // 是否上传到baidu网盘
	ReleaseDate          time.Time `json:"releaseDate" gorm:"not null" spider:"releaseDate" `                                                                    // 发布日期.可以是空字符串
	AuthorConcat         string    `json:"authorConcat" gorm:"not null;uniqueIndex:idx_comic_unique;size:500; check:author_concat <> ''" spider:"authorConcat" ` // 作者.不能是空字符串。组合索引
	AuthorConcatType     int       `json:"authorConcatType" gorm:"not null" spider:"authorConcatType" `                                                          // 作者拼接方式，不能空。：0 默认，按爬取顺序拼接，1: 按字母升序拼接 2:按我的意愿拼接 3: 参考最权威的网站拼接(b比如有声书，参考喜马拉雅，韩漫参考toptoon，小说参考 起点-建议0 /3

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除

	// 关联外键写法
	Country       Country      `gorm:"foreignKey:CountryId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"` // 可选：级联操作
	Website       Website      `gorm:"foreignKey:WebsiteId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PornType      PornType     `gorm:"foreignKey:PornTypeId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Type          Type         `gorm:"foreignKey:TypeId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Process       Process      `gorm:"foreignKey:ProcessId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Stats         ComicMyStats `gorm:"foreignKey:ComicID;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" spider:"stats"` // 漫画 统计
	LatestChapter ChapterMy    `gorm:"foreignKey:LatestChapterId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" spider:"latestChapter"`

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
type ComicMyStats struct {
	ID              int  `json:"id" gorm:"primaryKey;autoIncrement"`
	ComicID         int  `gorm:"uniqueIndex"`                               // 外键，唯一索引保证一对一
	LatestChapterId *int `json:"latestChapterId" spider:"latestChapterId" ` // 最新章节id。可为空，因为爬书的时候，章节表还没有内容。冗余1个，为了查询方便，不join影响性能

	// 频繁更新字段
	Star                      float64   `json:"star" gorm:"not null" spider:"star" `                           // 评分
	LatestChapterName         string    `json:"latestChapterName" gorm:"not null" spider:"latestChapterName" ` // 更新到多少集, 字符串,最新章节.可以是空字符串
	Hits                      int       `json:"hits" gorm:"not null" spider:"hits" `                           // 人气
	TotalChapter              int       `json:"totalChapter" gorm:"not null" spider:"totalChapter" `           // 总章节数
	LastestChapterReleaseDate time.Time `json:"lastestChapterReleaseDate" spider:"lastestChapterReleaseDate" ` // 最新章节发布时间，可以空

	// 外键结构
	// LatestChapter Chapter `gorm:"foreignKey:LatestChapterId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" spider:"latestChapter"` // 考虑删除，想着都是冗余了，如果用不到就先删除，用到再说
}

// 实现stringutils 里 处理空格接口
func (c *ComicMy) TrimSpaces() {
	// 只要是stirng类型，就去掉前后空格
	c.Name = strings.TrimSpace(c.Name)
	c.ComicUrlApiPath = strings.TrimSpace(c.ComicUrlApiPath)
	c.CoverUrlApiPath = strings.TrimSpace(c.CoverUrlApiPath)
	c.BriefShort = strings.TrimSpace(c.BriefShort)
	c.BriefLong = strings.TrimSpace(c.BriefLong)
	c.AuthorConcat = strings.TrimSpace(c.AuthorConcat)
}

// 实现stringutils 里 处理空格接口 - 处理子表 -统计数据
func (c *ComicMyStats) TrimSpaces() {
	c.LatestChapterName = strings.TrimSpace(c.LatestChapterName)
}

// 实现stringutils 里 繁体转简体 接口
func (c *ComicMy) Trad2Simple() {
	// 只要是string类型，都处理
	c.Name, _ = langutil.TraditionalToSimplified(c.Name)
	c.ComicUrlApiPath, _ = langutil.TraditionalToSimplified(c.ComicUrlApiPath)
	c.CoverUrlApiPath, _ = langutil.TraditionalToSimplified(c.CoverUrlApiPath)
	c.BriefShort, _ = langutil.TraditionalToSimplified(c.BriefShort)
	c.BriefLong, _ = langutil.TraditionalToSimplified(c.BriefLong)
	c.AuthorConcat, _ = langutil.TraditionalToSimplified(c.AuthorConcat)
}

// 实现stringutils 里 繁体转简体 接口 - 处理子表 -统计数据
func (c *ComicMyStats) Trad2Simple() {
	c.LatestChapterName, _ = langutil.TraditionalToSimplified(c.LatestChapterName)
}

// 实现 业务数据清理接口 - comicMy 表
func (c *ComicMy) BusinessDataClean() {
	// -- string 类型
	// 去除 http 协议头 --
	if len(c.ComicUrlApiPath) > 8 && (strings.HasPrefix(c.ComicUrlApiPath, "http://") || strings.HasPrefix(c.ComicUrlApiPath, "https://")) {
		log.Info("业务数据清理, ComicUrlApiPath 有http前缀, 去除. apiPath= ", c.ComicUrlApiPath)
		c.ComicUrlApiPath = stringutil.TrimHttpPrefix(c.ComicUrlApiPath)
	}

	if len(c.CoverUrlApiPath) > 8 && (strings.HasPrefix(c.CoverUrlApiPath, "http://") || strings.HasPrefix(c.CoverUrlApiPath, "https://")) {
		log.Info("业务数据清理, CoverUrlApiPath 有http前缀, 去除. apiPath= ", c.CoverUrlApiPath)
		c.CoverUrlApiPath = stringutil.TrimHttpPrefix(c.CoverUrlApiPath)
	}

	// -- int 类型
	// end 完结状态 --
	/* 判断逻辑
		  processId是人为传的，
	            - processId = 1, 表示待分类，end应该 == processId = 1
	                - end 爬不到, == 1
	                - end 爬到了， == 爬到的值 (2或3)
	            - processId = 2, 表示连载，  end应该 == processId = 2
	            - processId = 3, 表示完结，  end应该 == processId = 3
	*/
	c.End = 1             // 默认是1 - 待分类 / 不知道。反正不让是0. 要和processId 保持一致
	if c.ProcessId == 1 { // 不知道，需要机器自行判断。除非人 特别确认是完结/连载，否则前端传参，都传1
		if strings.Contains(c.Stats.LatestChapterName, "休刊公告") || strings.Contains(c.Stats.LatestChapterName, "后记") {
			c.End = 3 //完结
		}
		// 连载不知道咋判断
	} else {
		c.End = c.ProcessId
	}
}

// 实现 业务数据清理接口 - comicMyStats 表
func (c *ComicMyStats) BusinessDataClean() {
	// -- int 类型
	// 评分 --
	// 评分超过10，就置为0。可能是人为设置错了。0代表未设置
	if c.Star > 10 {
		log.Infof("进行业务数据清洗, c.star=%v >10, 重置为0", c.Star)
		c.Star = 0
	}
}

// 实现 数据清理统一入口 - comicMy 表
func (c *ComicMy) DataClean() {
	c.TrimSpaces()
	c.Trad2Simple()
	c.BusinessDataClean()
}

// 实现 数据清理统一入口 - comicMyStats 表
func (c *ComicMyStats) DataClean() {
	c.TrimSpaces()
	c.Trad2Simple()
	c.BusinessDataClean()
}
