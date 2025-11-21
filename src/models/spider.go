package models

import (
	"strings"
)

// 爬取请求body 结构体
/*
	例如：
	{
		"websitePrefix": "www.manhuagui.com",  // 网站前缀，现在想的是最后不带/
		"url": "list/c1-p",                    // 排除前缀后的，url路径，需要带/
		"needTcp": 1,   					   // 完整请求，是否需要带 http / https 因为有的爬取的 book的链接，有的带http，有的不带
		"needHttps": 1, 					   // 完整请求，是否需要带  https
		"endNum": 5, 						   // 尾页号码
		"endJudgeMethod": 2,                   // 完结判断方式 0：全部写成 未完结false 1：全部写成完结true 2：程序自动判断
		"countryId": 1,                        // 国家id
		"websiteId": 3,                        // 网站id
		"pornTypeId": 1,                       // 色情id
		"typeId": 2                            // 类型id
		"End": 0                               // 是否完结
	}
*/
type SpiderRequestBody struct {
	Url           string `json:"url" `           // 每页请求链接, 不带具体数字
	WebsitePrefix string `json:"websitePrefix" ` // 请求前缀，网站url
	NeedTcp       int    `json:"needTcp" `       // 完整请求，是否需要带 http / https
	NeedHttps     int    `json:"needHttps" `     // 完整请求，是否需要带  https
	EndNum        int    `json:"endNum" `        // 尾页号码
	CountryId     int    `json:"countryId" `     // 国家id
	WebsiteId     int    `json:"websiteId" `     // 网站id
	PornTypeId    int    `json:"pornTypeId" `    // 总分类id
	TypeId        int    `json:"typeId" `        // 类型id
	End           int    `json:"end" `           // 是否完结
}

// 通用的，爬取请求body 结构体
/*
	例如：
	{
		"websiteId": 3,                        // 网站id
		"pornTypeId": 1,                       // 色情类型id
		"countryId": 1,                        // 国家id
		"targetSiteTypeId": 1,                 // 目标网站类型id,即：这个网站自己是怎么分类型的，如：悬疑、冒险等
		"typeId": 2                            // 类型id。 我的网站自己是怎么分类型的，如：悬疑、冒险等
		//
		"websitePrefix": "www.manhuagui.com",  // 网站前缀，现在想的是最后不带/
		"url": "list/c1-p",                    // 排除前缀后的，url路径，需要带/
		"needTcp": 1,   					   // 完整请求，是否需要带 http / https 因为有的爬取的 book的链接，有的带http，有的不带
		"needHttps": 1, 					   // 完整请求，是否需要带  https
		"endNum": 5, 						   // 尾页号码
		"endJudgeMethod": 2,                   // 完结判断方式 0：全部写成 未完结false 1：全部写成完结true 2：程序自动判断
	}
*/
type SpiderRequestBody2 struct {
	// 和数据表外键相关的，有顺序区分的
	WebsiteId        int `json:"websiteId" `        // 网站id
	PornTypeId       int `json:"pornTypeId" `       // 色情类型id
	CountryId        int `json:"countryId" `        // 国家id
	TargetSiteTypeId int `json:"targetSiteTypeId" ` // 目标网站类型id,即：这个网站自己是怎么分类型的，如：悬疑、冒险等
	TypeId           int `json:"typeId" `           // 类型id。我的网站类型id,即：我的网站自己是怎么分类型的，如：悬疑、冒险等
	//
	WebsitePrefix  string `json:"websitePrefix" `  // 请求前缀，网站url，最后不带/
	ApiPath        string `json:"apiPath" `        // 接口路径。最前面需要带/ 如: /category/comic
	NeedTcp        bool   `json:"needTcp" `        // 完整请求，是否需要带 http / https
	NeedHttps      bool   `json:"needHttps" `      // 完整请求，是否需要带  https
	EndNum         int    `json:"endNum" `         // 尾页号码
	EndJudgeMethod int    `json:"endJudgeMethod" ` // 完结判断方式 0：全部写成 未完结false 1：全部写成完结true 2：程序自动判断

	// 实现时，只能map里传一个key
	ParamArr []map[string]string `json:"paramArr"` // 请求参数 例如: ?后面的就是 https://kxmanhua.com/manga/library?type=0&complete=1&page=1&orderby=1. 保证是有序结构，能确保传的值对应上type/complete
}

// 实现 stringutils 里处理空格 接口
func (s *SpiderRequestBody2) TrimSpaces() {
	s.WebsitePrefix = strings.TrimSpace(s.WebsitePrefix) // 网站前缀 -- 域名
	s.ApiPath = strings.TrimSpace(s.ApiPath)             // 接口路径

	// 请求参数
	// 遍历数组
	for _, param := range s.ParamArr {
		// 遍历map
		for k, v := range param {
			param[k] = strings.TrimSpace(v)
		}
	}
}
