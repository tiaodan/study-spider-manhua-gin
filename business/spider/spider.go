package spider

import (
	"fmt"
	"study-spider-manhua-gin/errorutil"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/util/langutil"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func Spider(context *gin.Context) {
	log.Debug("爬虫开始-------------------------------")
	c := colly.NewCollector()

	// 设置请求限制（每秒最多2个请求）
	// spider.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*",
	// 	Parallelism: 2,
	// 	RandomDelay: 5 * time.Second,
	// })

	// 修改为解析所有元素并检查文本内容
	// spider.OnHTML("*", func(e *colly.HTMLElement) { // 修改选择器为"*"匹配所有元素
	// 	text := e.Text
	// 	if strings.Contains(text, "还没有看过的漫画") { // 新增文本判断逻辑
	// 		log.Debug("找到目标文本: %s", text)
	// 	}
	// 	link := e.Attr("href")
	// 	log.Debug("发现链接：%s", link)
	// })

	// 获取html内容,每成功匹配一次, 就执行一次逻辑。这个标签选只匹配一次的
	c.OnHTML(".cate-comic-list", func(e *colly.HTMLElement) {
		// log.Debug("html 元素名称Name= ", e.Name) // 一般显示div a span 等标签名
		// log.Debug("html 元素DOM= ", e.DOM) // 就显示地址 &{[0xc0002fa540] 0xc00060e4e0 <nil>}
		// log.Debug("html 元素内容Text= ", e.Text)
		// baiduBtn := e.Attr("value")
		// names := e.ChildText(".comic-feature")
		// e.ForEach(".common-comic-item", func(i int, h *colly.HTMLElement) {
		// 	// names := h.ChildText(".comic__title")
		// 	log.Debug("第%d个漫画名称= %s", i, h.Attr("comic__title"))
		// })
		// 漫画名
		// 在选一个可以foreach匹配的子标签
		log.Debug("匹配到.cate-comic-list")
		e.ForEach(".common-comic-item", func(i int, element *colly.HTMLElement) {
			// 思路：
			// 1. 爬名字
			// 2. 爬更新到 ?集
			// 3. 爬人气
			// 4. 爬封面链接
			// 5. 爬漫画链接
			// 5.1. 去重
			// 5.2 看是否需要转简体
			// 6. 把参数赋值给 comic对象
			// 7. 插入数据库

			// 1. 爬名字
			comicName := element.ChildText(".comic__title")
			// 转换为简体中文
			simplifiedName, err := langutil.TraditionalToSimplified(comicName)
			if err != nil {
				log.Errorf("转换为简体中文失败: %v", err)
				simplifiedName = comicName // 如果转换失败，使用原名称
			}

			// 2. 爬更新到 ?集
			// <p class="comic-update">更至：<a class="hl" href="/chapter/82314" target="_blank">休刊公告</a></p>
			// link := element.Attr("a.hl@href")
			// updateNum := element.ChildText("a.h1")
			updateNum := element.DOM.ChildrenFiltered("p .comic-update").ChildrenFiltered("a")
			log.Debug(updateNum)
			log.Debugf("当前%d, 漫画名称转简体= %s -> %s", i+1, comicName, simplifiedName)
			// log.Debugf("当前%d, 漫画更新至 %s 集", i+1, upDateNum)
		})

	})

	// 修改本地文件访问路径为file协议格式
	err := c.Visit("http://localhost:8080/test/index.html") // 使用file协议访问本地文件

	errorutil.ErrorPrint(err, "访问失败")
	log.Debug("爬虫结束-------------------------------")
}

// 爬百度主页, 熟悉练手
func SpiderBaiduTest(context *gin.Context) {
	c := colly.NewCollector()

	// 获取 "百度一下" 按钮文本
	c.OnHTML("#su", func(e *colly.HTMLElement) {
		baiduBtn := e.Attr("value")
		fmt.Println(baiduBtn)
	})

	// 开始访问
	err := c.Visit("http://www.baidu.com/") // 百度 测试
	errorutil.ErrorPrint(err, "访问失败")
}
