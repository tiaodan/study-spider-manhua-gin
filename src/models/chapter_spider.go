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
	Id                 int       `json:"id" gorm:"primary_key;autoIncrement" spider:"id" `                         // 主键, 数据库id,
	ChapterNum         int       `json:"chapterNum" gorm:"not null;" spider:"chapterNum"`                          // 章节编号,从1开始
	ChapterSubNum      int       `json:"chapterSubNum" gorm:"not null;" spider:"chapterSubNum"`                    // 章节子编号,从1/0 ?开始(从0吧，因为int默认0)。目的:如果爬的有问题，比如少几章，可以人为插入,还不影响顺序
	ChapterRealSortNum int       `json:"chapterRealSortNum" gorm:"not null;" spider:"chapterRealSortNum"`          // 章节真实排序号-真实显示也用的它,可以是负数，负数在 有声书中，可以表示 花絮，前言、介绍等非正式语音
	Name               string    `json:"name" gorm:"not null;check:name <> ''" spider:"name" `                     // 章节名称.不能为空字符串？
	UrlApiPath         string    `json:"urlApiPath" gorm:"not null;check:url_api_path <> ''" spider:"urlApiPath" ` // 章节api路径.不能为空字符串
	ReleaseDate        time.Time `json:"releaseDate" gorm:"not null;" spider:"releaseDate"`                        // 发布时间
	SpiderStatus       int       `json:"spiderStatus" gorm:"not null;" spider:"spiderStatus"`                      // 爬取状态,0:未爬取,1:已爬取,2:爬取失败

	// 外键
	ParentId int `json:"parent_id" gorm:"not null;index" spider:"parentId"` // 父id, 如：comic表、audiobook表

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
	// 还没想好要清理哪些业务字段
}

// 实现 数据清理统一入口
func (c *ChapterSpider) DataClean() {
	c.TrimSpaces()        // 处理空格
	c.Trad2Simple()       // 繁体转简体
	c.BusinessDataClean() // 业务数据清理
}
