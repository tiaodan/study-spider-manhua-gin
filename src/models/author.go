/**
作者表
*/

package models

// 定义 作者 模型
type Author struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`   // 数据库id
	NameId int    `gorm:"not null;unique"`            // 唯一索引 author nameid
	Name   string `gorm:"not null;check:name <> ''" ` // 作者名字
}
