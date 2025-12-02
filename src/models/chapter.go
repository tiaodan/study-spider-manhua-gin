/*
*
章节表 模板
*/
package models

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	Id   int    `json:"id" gorm:"primary_key;autoIncrement"`
	Name string `json:"name"`

	// gorm自带时间更新，软删除
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 启用软删除，并设置索引,加快查询. NULL表示没删除
}
