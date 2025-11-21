package models

// 定义网站模型
type Website struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	NameId    int    `gorm:"not null;unique"`
	Name      string `gorm:"not null"`
	Domain    string `gorm:"not null"` // 域名，如：www.google.com
	NeedProxy bool   `gorm:"not null"` // 是否需要翻墙
	IsHttps   bool   `gorm:"not null"` // 网站是否Https前缀，如果是false,默认就是 http头
}
