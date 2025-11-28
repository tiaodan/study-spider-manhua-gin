package models

// 定义国家模型
/*
1 - 待分类
2 - 中国
3 - 韩国
4 - 欧美
5 - 日本
*/
type Country struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null;unique;check:name <> ''" ` // 国家名称. 唯一索引
}
