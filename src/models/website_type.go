/*
*
网站类型表。如：漫画、小说、有声书、影视
默认写死数据：
1 - 待分类
2 - 漫画
3 - 小说
4 - 有声书
5 - 视频
6 - 音乐
7 - 网盘
8 - 综合娱乐
*/
package models

type WebsiteType struct {
	Id   int    `gorm:"primaryKey;autoIncrement"`          // id
	Name string `gorm:"unique;not null;check:name <> ''" ` // 类型名称,唯一索引，如漫画、小说、有声书、影视
}
