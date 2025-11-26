/*
*
网站类型表。如：漫画、小说、有声书、影视
默认写死数据：
0 - 待分类
1 - 漫画
2 - 小说
3 - 有声书
4 - 视频
5 - 音乐
6 - 网盘
*/
package models

type WebsiteType struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`
	NameId int    `gorm:"not null;unique"`            // 类型id,唯一索引
	Name   string `gorm:"not null;check:name <> ''" ` // 类型名称，如漫画、小说、有声书、影视
}
