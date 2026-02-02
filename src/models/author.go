/**
作者表
*/

package models

import (
	"strings"
	"study-spider-manhua-gin/src/util/langutil"
)

// 定义 作者 模型
/*
默认：
1 - 佚名
*/
type Author struct {
	Id int `gorm:"primaryKey;autoIncrement" json:"id" ` // 主键

	// 外键
	WebsiteId int `gorm:"not null;uniqueIndex:idx_website_id_namer_unique" json:"websiteId" ` // 组合索引

	// 其它
	Name string `gorm:"not null;;uniqueIndex:idx_website_id_namer_unique;size:50;check:name <> ''" json:"name" spider:"name" ` // 作者名字. 唯一索引, 之前写法：gorm:"not null;unique;
}

// 实现stringutils 里 处理空格接口
func (a *Author) TrimSpaces() {
	// 只要是stirng类型，就去掉前后空格
	a.Name = strings.TrimSpace(a.Name)
}

// 实现stringutils 里 繁体转简体 接口
func (a *Author) Trad2Simple() {
	// 只要是string类型，都处理
	a.Name, _ = langutil.TraditionalToSimplified(a.Name)
}

// 实现 业务数据清理接口
func (a *Author) BusinessDataClean() {
	// 暂时没有要处理的
}

// 实现stringutils 里 中文字符转英文 接口
func (a *Author) ChChar2En() {
	a.Name = langutil.ChineseCharToEnglish(a.Name)
}

// 实现 数据清理统一入口
func (a *Author) DataClean() {
	a.TrimSpaces()        // 处理空格
	a.Trad2Simple()       // 繁体转简体
	a.ChChar2En()         // 中文字符转英文,并去除无用字符
	a.BusinessDataClean() // 业务数据清理
}
