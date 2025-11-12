package models

// 爬取请求body 结构体
//
//	{
//	    "url": "/category/list/2/page/",
//	    "websitePrefix": "https://9y01.xyz"
//	    "needTcp": 0,
//	    "needHttps": 0,
//	    "endNum": 32
//	}
type SpiderRequestBody struct {
	Url           string `json:"url" `           // 每页请求链接, 不带具体数字
	WebsitePrefix string `json:"websitePrefix" ` // 请求前缀，网站url
	NeedTcp       int    `json:"needTcp" `       // 完整请求，是否需要带 http / https
	NeedHttps     int    `json:"needHttps" `     // 完整请求，是否需要带  https
	EndNum        int    `json:"endNum" `        // 尾页号码
	CountryId     int    `json:"countryId" `     // 国家id
	WebsiteId     int    `json:"websiteId" `     // 网站id
	SexTypeId     int    `json:"sexTypeId" `     // 总分类id
	TypeId        int    `json:"typeId" `        // 类型id
	End           int    `json:"end" `           // 是否完结
}
