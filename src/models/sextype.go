package models

// 类别: 如普通漫画、色漫
type PornType struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`
	NameId int    `gorm:"not null;unique"`
	Name   string `gorm:"not null;check:name <> ''" `
}
