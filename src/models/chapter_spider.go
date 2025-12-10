/*
*
章节表 模板
*/
package models

import (
	"strings"
	"study-spider-manhua-gin/src/util/langutil"
	"time"

	"gorm.io/gorm"
)

// 章节表 - 爬取的
/*
唯一索引: 父id-章节号码-章节子号码
*/
type ChapterSpider struct {
	Id                   int       `json:"id" gorm:"primary_key;autoIncrement" spider:"id" `                                  // 主键, 数据库id,
	ChapterNum           int       `json:"chapterNum" gorm:"not null;uniqueIndex:idx_comic_unique" spider:"chapterNum"`       // 章节编号,从1开始。组合索引
	ChapterSubNum        int       `json:"chapterSubNum" gorm:"not null;uniqueIndex:idx_comic_unique" spider:"chapterSubNum"` // 章节子编号,从1/0 ?开始(从0吧，因为int默认0)。目的:如果爬的有问题，比如少几章，可以人为插入,还不影响顺序。组合索引
	ChapterRealSortNum   int       `json:"chapterRealSortNum" gorm:"not null;" spider:"chapterRealSortNum"`                   // 章节真实排序号-真实显示也用的它,可以是负数，负数在 有声书中，可以表示 花絮，前言、介绍等非正式语音
	Name                 string    `json:"name" gorm:"not null;check:name <> ''" spider:"name" `                              // 章节名称.不能为空字符串？
	UrlApiPath           string    `json:"urlApiPath" gorm:"not null;check:url_api_path <> ''" spider:"urlApiPath" `          // 章节api路径.不能为空字符串
	ReleaseDate          time.Time `json:"releaseDate" gorm:"not null;" spider:"releaseDate"`                                 // 发布时间
	SpiderEndStatus      int       `json:"spiderEndStatus" gorm:"not null;" spider:"spiderEndStatus"`                         // 爬取结束状态,0:未爬取,1:已爬取,2:爬取失败
	DownloadEndStatus    int       `json:"downloadEndStatus" gorm:"not null;" spider:"downloadEndStatus"`                     // 下载结束状态,0:未爬取,1:已爬取,2:爬取失败
	UploadAwsEndStatus   int       `json:"uploadAwsEndStatus" gorm:"not null;" spider:"uploadAwsEndStatus"`                   // 上传aws结束状态,0:未爬取,1:已爬取,2:爬取失败
	UploadBaiduEndStatus int       `json:"uploadBaiduEndStatus" gorm:"not null;" spider:"uploadBaiduEndStatus"`               // 上传baidu结束状态,0:未爬取,1:已爬取,2:爬取失败

	// 外键
	ParentId int `json:"parent_id" gorm:"not null;uniqueIndex:idx_comic_unique" spider:"parentId"` // 父id, 如：comic表、audiobook表。组合索引

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除
}

// ------ func
// -- 实现各种数据清洗接口

// 实现stringutils 里 处理空格接口
func (c *ChapterSpider) TrimSpaces() {
	// string 类型
	c.Name = strings.TrimSpace(c.Name)
	c.UrlApiPath = strings.TrimSpace(c.UrlApiPath)
}

// 实现stringutils 里 繁体转简体 接口
func (c *ChapterSpider) Trad2Simple() {
	// string 类型
	c.Name, _ = langutil.TraditionalToSimplified(c.Name)
	c.UrlApiPath, _ = langutil.TraditionalToSimplified(c.UrlApiPath)
}

// 实现 业务数据清理接口
func (c *ChapterSpider) BusinessDataClean() {
	// ChapterRealSortNum -》 真正用的序号， = ChapterNum (爬取到的)
	c.ChapterRealSortNum = c.ChapterNum

	// -- 日期类型，不让有0000-00-00 00:00:00 / 0001-01-01 00:00:00
	if c.ReleaseDate.IsZero() {
		// -- 考虑健壮性，而且有的字段 not null ,不让时间用nil。所以设置一个 sql支持的默认值。如：1001-01-01 00:00:00
		c.ReleaseDate = time.Date(1001, 1, 1, 0, 0, 0, 0, time.UTC) // 默认 1001-01-01 00:00:00
	}
}

// 实现stringutils 里 中文字符转英文 接口
func (c *ChapterSpider) ChChar2En() {
	c.Name = langutil.ChineseCharToEnglish(c.Name)
}

// 实现 数据清理统一入口
func (c *ChapterSpider) DataClean() {
	c.TrimSpaces()        // 处理空格
	c.Trad2Simple()       // 繁体转简体
	c.ChChar2En()         // 中文字符转英文,并去除无用字符
	c.BusinessDataClean() // 业务数据清理
}
