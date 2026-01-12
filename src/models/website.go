package models

// 定义网站模型
/*
唯一索引：Name + Domain
*/
type Website struct {
	Id                                   int    `gorm:"primaryKey;autoIncrement"`
	Name                                 string `gorm:"not null; uniqueIndex:idx_website_unique;size:150;check:name <> ''" `   // 网站名称.唯一组合索引
	Domain                               string `gorm:"not null; uniqueIndex:idx_website_unique;size:500;check:domain <> ''" ` // 域名，如：www.google.com:8888 唯一组合索引. 没有Port参数，要想写port,在host里写。例如：localhost:8888
	NeedProxy                            bool   `gorm:"not null"`                                                              // 是否需要翻墙
	IsHttps                              bool   `gorm:"not null"`                                                              // 网站是否Https前缀，如果是false,默认就是 http头
	IsRefer                              bool   `gorm:"not null"`                                                              // 是否是参考/参照网站，比如电影：参考网站：豆瓣，小说：起点，漫画：toptoon
	CoverURLIsNeedHttps                  bool   `gorm:"not null"`                                                              // 封面URL是否需要https前缀
	ChapterContentURLIsNeedHttps         bool   `gorm:"not null"`                                                              // 章节内容URL是否需要https前缀
	CoverURLConcatRule                   string `gorm:"not null;check:cover_url_concat_rule <> ''" `                           // 封面URL拼接规则。一般实现方式：判断website是否用https + website表-domain + book表->cover_url_api_path
	ChapterContentURLConcatRule          string `gorm:"not null;check:chapter_content_url_concat_rule <> ''" `                 // 章节内容URL拼接规则。一般实现方式：判断website是否用https + website表-domain + book表->cover_url_api_path
	CoverDomain                          string `gorm:"not null;check:cover_domain <> ''" `                                    // 封面域名。也可以填ip
	ChapterContentDomain                 string `gorm:"not null;check:chapter_content_domain <> ''" `                          // 章节内容域名。也可以填ip
	BookCanSpiderType                    string `gorm:"not null;check:book_can_spider_type <> ''" `                            // book可以的爬取方式:. 可选:json/html/both/bothno  json还是html还是2者都行,还是都不行.必须填1个
	ChapterCanSpiderType                 string `gorm:"not null;check:chapter_can_spider_type <> ''" `                         // chapter 可以的爬取方式:. 可选:json/html/both/bothno  json还是html还是2者都行,还是都不行.必须填1个
	BookSpiderReqBodyEgServerFilepath    string `gorm:"not null" `                                                             // 爬取book时,请求体内容示例,后台服务器路径,不要求必须有值,字符串空也可以. 如："爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html"
	ChapterSpiderReqBodyEgServerFilepath string `gorm:"not null" `                                                             // 爬取chapter时,请求体内容示例,后台服务器路径,不要求必须有值,字符串空也可以. 如："爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html"
	StarType                             string `gorm:"not null"`                                                              // 打分方式,自己打/参考其他。没有参考的网站，或者网站评分系统无参考价值，就自己主观打分。 比如：my -> 表示我自己打的; copy_toptoon -> 表示从toptoon网站爬取的

	// -- 关联外键 website_type 表
	WebsiteTypeId int `gorm:"not null"` // website_type 表的外键id
	// 关联 website_type 表 nameId,要求：website_type表，删除时，不让删，更新时，同步更新 --
	WebsiteType WebsiteType `gorm:"foreignKey:WebsiteTypeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
