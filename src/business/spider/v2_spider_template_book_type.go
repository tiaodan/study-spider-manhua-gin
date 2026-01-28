/*
V2 版本：考虑把爬取流程做成,更一劳永逸的方式。思路：一般一劳永逸就是: 配置驱动+?(策略模式)

思路：
程序可以分为：
- 配置 -》 配置驱动一切
- 变量/数据结构
- 策略调度/执行策略
- 算法
*/

package spider

import (
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/tidwall/gjson"
)

// GetOneTypeAllBookUseCollyByMappingV2 v2版本的HTML爬取函数，直接接受选择器参数
/*
参数：
	websiteSpiderCfg: 网站爬取相关配置，从配置文件 v2-spider-config.yaml 读取
*/
func GetOneTypeAllBookUseCollyByMappingV2[T any](requestBodyData []byte, mapping map[string]models.ModelHtmlMapping, spiderUrlArr []string, bookArrCssSelector string, bookArrItemCssSelector string, websiteSpiderCfg *config.WebsiteConfig) [][]T {
	/*
		思路：参考V1逻辑，还用colly+queue方式，只是有些内容 从配置读取
		0. 读取前端传参
		1. 建一个爬虫对象
		2. 设置并发数，和爬取限制
		3. 写c.OnHTML 处理逻辑
		4. 创建队列 queue
		5. 添加任务到队列
		6. 启动队列
	*/
	// 使用传入的选择器参数，而不是从JSON解析
	log.Debug("v2爬取html, bookArrCssSelector = ", bookArrCssSelector)
	log.Debug("v2爬取html, bookArrItemCssSelector = ", bookArrItemCssSelector)

	// 0. 读取前端传参
	// 1. gjson 读取 前端 JSON 里 spiderTag -> website字段 --
	website := gjson.Get(string(requestBodyData), "spiderTag.website").String() // websiteTag - website
	table := gjson.Get(string(requestBodyData), "spiderTag.table").String()     // websiteTag - table

	websiteId := gjson.Get(string(requestBodyData), "websiteId").Int()               // 网站id
	pornTypeId := gjson.Get(string(requestBodyData), "pornTypeId").Int()             // 色情类型id
	countryId := gjson.Get(string(requestBodyData), "countryId").Int()               // 国家id
	typeId := gjson.Get(string(requestBodyData), "typeId").Int()                     // 类型id
	processId := gjson.Get(string(requestBodyData), "processId").Int()               // 进程：完结状态 id
	authorConcatType := gjson.Get(string(requestBodyData), "authorConcatType").Int() // 作者拼接方式 id
	needTcp := gjson.Get(string(requestBodyData), "needTcp").Bool()                  // 是否需要tcp 头
	coverNeedTcp := gjson.Get(string(requestBodyData), "coverNeedTcp").Bool()        // 封面链接是否需要tcp 头
	endNum := gjson.Get(string(requestBodyData), "endNum").Int()                     // 结束页 号码
	// adultArrGjsonResult := gjson.GetBytes(requestBodyData, "adult").Array()      // 数组 - adult 内容 - html 用不到, 一会删
	// bookArrCssSelector := gjson.Get(string(requestBodyData), "bookArrCssSelector").String()         // 获取某页所有书 用的CSS选择器。不从前端拿了，从配置拿
	// bookArrItemCssSelector := gjson.Get(string(requestBodyData), "bookArrItemCssSelector").String() // 获取某本书 用的CSS选择器。不从前端拿了，从配置拿

	log.Info("爬取html,前端传参= ", string(requestBodyData))
	log.Debug("爬取html,前端传参. piderTag.website = ", website)
	log.Debug("爬取html,前端传参. piderTag.table = ", table)
	log.Debug("爬取html,前端传参. websiteId = ", websiteId)
	log.Debug("爬取html,前端传参. pronTypeId = ", pornTypeId)
	log.Debug("爬取html,前端传参. countryId = ", countryId)
	log.Debug("爬取html,前端传参. typeId = ", typeId)
	log.Debug("爬取html,前端传参. processId = ", processId)
	log.Debug("爬取html,前端传参. authorConcatType = ", authorConcatType)
	log.Debug("爬取html,前端传参. needTcp = ", needTcp)
	log.Debug("爬取html,前端传参. coverNeedTcp = ", coverNeedTcp)
	log.Debug("爬取html,前端传参. endNum = ", endNum)
	log.Debug("爬取html,前端传参. bookArrCssSelector = ", bookArrCssSelector)
	log.Debug("爬取html,前端传参. bookArrItemCssSelector = ", bookArrItemCssSelector)
	// 1. 建一个爬虫对象
	c := colly.NewCollector()

	// 2. 设置并发数，和爬取限制
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Duration(config.Cfg.Spider.Public.SpiderType.RandomDelayTime) * time.Second, // 请求发送前触发。模仿人类，随机等待几秒，再请求。如果queue同时给了3条URL，那每条url触发请求前，都要随机延迟下
	})

	// 3. 写c.OnHTML 处理逻辑
	var allPageBookArr [][]T // 存放爬好的 obj，因为要返回泛型，所以用T ,以前写法：comicArr := []models.ComicSpider{}. 二维数组，里面存放 onePageBookArr = []models.T
	var mu sync.Mutex        // 添加互斥锁
	// 遍历每一个 bookArr .c.OnHTML() 根据 CSS选择器, 就让触发1次
	c.OnHTML(bookArrCssSelector, func(eBookArr *colly.HTMLElement) {
		log.Debug("-------------- 匹配 bookArr = ", eBookArr.Text)

		// 遍历每一个 bookArrItem, 用forEach. colly，用Html遍历
		var onePageBookArr []T
		eBookArr.ForEach(bookArrItemCssSelector, func(i int, e *colly.HTMLElement) {
			// 1. 获取能获取到的
			// 通过mapping -> 转成1个对象
			// 创建对象comic
			var comicT T
			comicSpiderStats := models.ComicSpiderStats{} // 子表，统计数据
			log.Info("delete comicSpiderStats = ", comicSpiderStats)

			// 通过mapping 爬内容
			rawResult := GetOneObjByCollyMapping(e, mapping)
			if rawResult != nil {
				// 应用v2的transforms --
				processedResult := make(map[string]interface{})
				for fieldName, rawValue := range rawResult {
					// 从config中获取字段配置和transforms
					if fieldConfig, exists := websiteSpiderCfg.Extract.Mappings[fieldName]; exists && len(fieldConfig.Transforms) > 0 {
						// 创建field mapper并应用transforms
						fieldMapper := NewFieldMapper()
						processedValue, err := fieldMapper.transformRegistry.ApplyTransforms(fieldConfig.Transforms, rawValue, fieldMapper.configLoader)
						if err != nil {
							log.Errorf("v2字段 %s 转换失败: %v", fieldName, err)
							processedResult[fieldName] = rawValue // 使用原始值
						} else {
							processedResult[fieldName] = processedValue
						}
					} else {
						processedResult[fieldName] = rawValue
					}
				}
				// 通过 model字段 spider，把爬出来的 map[string]any，转成 model对象 --
				MapByTag(processedResult, &comicT)
				log.Infof("这个是最后插入db前, 最准确的结果. 映射后的comic对象: %+v", comicT)
			}

			// 2. 设置对象值
			// -- T 类型 -》 具体struct 类型
			comic := any(comicT).(models.ComicSpider)

			// -- 进度id逻辑
			if processId == 1 {
				// 如果用户传 1 - 》程序自己判断
				comic.ProcessId = comic.End
			} else {
				// 如果是2/3, 就直接替换赋值
				comic.ProcessId = int(processId)
			}
			// -- 其它直接赋值
			comic.WebsiteId = int(websiteId)               // 网站id
			comic.PornTypeId = int(pornTypeId)             // 色情类型id
			comic.CountryId = int(countryId)               // 国家id
			comic.TypeId = int(typeId)                     // 类型id
			comic.AuthorConcatType = int(authorConcatType) // 作者拼接方式 id

			// 3 数据清洗
			comic.DataClean()

			// 4 把爬好的单个数据，放到数组里，准备插入数据库
			onePageBookArr = append(onePageBookArr, any(comic).(T))
		})

		// 3. 遍历完之后，加到 allPageBookArr 里 - 使用互斥锁保护 (因为要操作多个线程 用的共享对象 -》 allPageBookArr)
		mu.Lock()
		allPageBookArr = append(allPageBookArr, onePageBookArr)
		mu.Unlock()
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		log.Errorf("v2爬取页面出错: %v, URL: %s", err, r.Request.URL)
	})

	// 4. 创建队列 queue
	q, _ := queue.New(config.Cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum,
		&queue.InMemoryQueueStorage{MaxSize: config.Cfg.Spider.Public.SpiderType.QueuePoolMaxnum}) // 最多并发 N 个Url, 根据配置文件来

	// 5. 添加任务到队列
	for i := range spiderUrlArr {
		q.AddURL(spiderUrlArr[i])
	}
	// 6. 启动队列
	q.Run(c)

	return allPageBookArr
}
