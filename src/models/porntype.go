package models

// 类别: 如普通漫画、色漫
/*
1 - 待分类
2 - 普通漫画
3 - 色漫
*/
type PornType struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null;unique;check:name <> ''" ` // 唯一索引
}
