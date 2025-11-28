/**
漫画spider -作者-关系表

*/

package models

// 定义 漫画-作者-关系表 模型
type ComicMyAuthorRelation struct {
	Id uint `gorm:"primaryKey;autoIncrement"` // 数据库id

	// 外键
	ComicSpiderId int `json:"comicSpiderId" gorm:"not null; uniqueIndex:idx_comic_spider_name_id_unique" ` // 漫画-Spider库 id-外键 组合索引字段
	AuthorId      int `json:"authorId" gorm:"not null; uniqueIndex:idx_author_name_id_unique" `            // 作者 id-外键 组合索引字段

	// 关联外键写法，更新时，同步更新，删除时，不让删
	ComicSpider ComicSpider `gorm:"foreignKey:ComicSpiderId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" ` // 级联操作
	Author      Author      `gorm:"foreignKey:AuthorId;references:Id; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" `
}
