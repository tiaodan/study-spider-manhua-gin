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
	Id   uint   `gorm:"primaryKey;autoIncrement"`                        // 主键
	Name string `gorm:"not null;unique;check:name <> ''" spider:"name" ` // 作者名字. 唯一索引
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
