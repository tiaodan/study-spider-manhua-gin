package models

// 定义书籍类型，玄幻、悬疑、科幻等
/*
唯一索引：Name
默认数据：
1级分类：
	ID   Name
	1    待分类
	2    韩漫
	3    日漫
	4    真人漫画
	5    3D漫画
	6    欧美漫画
	7    同性
	8    同人志
	9    出版漫画

2级分类:

*/
type Type struct {
	Id     int    `gorm:"primaryKey;autoIncrement"`          // 主键
	Name   string `gorm:"not null;unique;check:name <> ''" ` // 唯一索引
	Level  int    `gorm:"not null"`                          // 1:一级分类，2:二级分类，3:三级分类
	Parent int    `gorm:"not null"`                          // 父级分类id
}
