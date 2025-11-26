/*
*
进度表
0 - 待分类
1 - 连载
2 - 完结
*/
package models

type Process struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`   // 数据库id
	NameId int    `gorm:"not null;unique"`            // 进度id,唯一索引
	Name   string `gorm:"not null;check:name <> ''" ` // 进度名称，如连载、完结
}
