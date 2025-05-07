// 拼多多订单数据模型, 存数据用的
package models

// 漫画数据
type Comic struct {
	Id         uint    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`                       // 数据库id,主键、自增
	Name       string  `json:"name" gorm:"not null; uniqueIndex:name_unique;size:150;column:name"` // 漫画名 数据库唯一索引
	Update     string  `json:"update" gorm:"not null;column:update"`                               // 更新到多少集, 字符串
	Hits       uint    `json:"hits" gorm:"not null;hits"`                                          // 人气
	ComicUrl   string  `json:"comicUrl" gorm:"not null;column:comic_url"`                          // 漫画链接
	CoverUrl   string  `json:"coverUrl" gorm:"not null;column:cover_url"`                          // 封面链接
	BriefShort string  `json:"briefShort" gorm:"not null;column:brief_short"`                      // 简介-短
	BriefLong  string  `json:"briefLong" gorm:"not null;column:brief_long"`                        // 简介-长
	End        int     `json:"end" gorm:"not null;column:end"`                                     // 漫画是否完结,如果完结是1
	Star       float64 `json:"star" gorm:"not null;column:star"`                                   // 评分
	NeedTcp    int     `json:"needTcp" gorm:"not null;column:need_tcp"`                            // 漫画是否需要http 或者https前缀,如果链接有了tcp,就应该是0
	// NeedHttps    int    `json:"needHttps" gorm:"not null;column:need_https"`                        // 注释。必须配合websitemodel配合。漫画是否需要https, 1-> https 0->http
	CoverNeedTcp int `json:"coverNeedTcp" gorm:"not null;column:cover_need_tcp"` // 封面是否需要http 或者https前缀,如果链接有了tcp,就应该是0
	// CoverNeedHttps int    `json:"coverNeedHttps" gorm:"not null;column:cover_need_https"`               // 注释。必须配合websitemodel配合。封面是否需要https, 1-> https 0->http,这个是根据网站是否https访问来确定的
	SpiderEnd      int `json:"spiderEnd" gorm:"not null;column:spider_end"`            // 爬取结束
	DownloadEnd    int `json:"downloadEnd" gorm:"not null;column:download_end"`        // 下载结束
	UploadAwsEnd   int `json:"uploadAwsEnd" gorm:"not null;column:upload_aws_end"`     // 是否上传到aws
	UploadBaiduEnd int `json:"uploadBaiduEnd" gorm:"not null;column:upload_baidu_end"` // 是否上传到baidu网盘

}
