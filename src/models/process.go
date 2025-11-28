/*
*
进度表

唯一索引：Name
*/
package models

// 定义进度模型
/*
	1 - 待分类
	2 - 连载
	3 - 完结
*/
type Process struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`          // 主键
	Name string `gorm:"not null;unique;check:name <> ''" ` // 进度名称，唯一索引. 如连载、完结
}
