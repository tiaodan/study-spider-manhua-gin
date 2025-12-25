// 配置驱动的API接口 控制器 - 全新实现，不修改任何现有代码
// 提供配置驱动的爬虫API，支持一劳永逸的网站扩展

package spider

import (
	"errors"
	"fmt"
	"io"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// ConfigDrivenAPIController 配置驱动的API
type ConfigDrivenAPIController struct {
	executor *SpiderExecutor
}

// NewConfigDrivenAPIController 创建配置驱动API
func NewConfigDrivenAPIController() *ConfigDrivenAPIController {
	return &ConfigDrivenAPIController{
		executor: NewSpiderExecutor(),
	}
}

// InitConfig 初始化配置（应用启动时调用）
func (api *ConfigDrivenAPIController) InitConfig(configPath string) error {
	loader := config.GetSpiderConfigLoader()
	if err := loader.LoadConfig(configPath); err != nil {
		return fmt.Errorf("加载爬虫配置失败: %v", err)
	}

	if err := loader.ValidateConfig(); err != nil {
		return fmt.Errorf("验证爬虫配置失败: %v", err)
	}

	log.Info("爬虫配置文件加载成功")
	return nil
}

// DispatchApi_OneTypeAllBookByHtml_ConfigDriven 配置驱动的书籍爬取API
// 对应原有的 DispatchApi_OneTypeAllBookByHtml，但使用配置驱动
func (api *ConfigDrivenAPIController) DispatchApi_OneTypeAllBookByHtml_ConfigDriven(c *gin.Context) {
	/*
		思路：
		1. 校验传参
		2. 数据清洗
		3. 业务逻辑 需要的数据校验 +清洗
		4. 执行核心逻辑
		- 读取html内容
		- 通过mapping 爬取字段，赋值给chapter_spider对象
		- 插入前, 数据清洗
		- 批量插入db - comic
		- 批量插入db - comic_stats
		5. 返回结果
	*/

	// 读取请求数据
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "读取请求数据失败: " + err.Error()})
		return
	}

	// 解析请求参数
	website := gjson.Get(string(data), "spiderTag.website").String()
	spiderUrl := gjson.Get(string(data), "spiderUrl").String()
	endNum := gjson.Get(string(data), "endNum").Int()
	websiteId := int(gjson.Get(string(data), "websiteId").Int())

	// 验证必需参数
	if website == "" {
		c.JSON(400, gin.H{"error": "缺少必需参数: spiderTag.website"})
		return
	}
	if spiderUrl == "" {
		c.JSON(400, gin.H{"error": "缺少必需参数: spiderUrl"})
		return
	}
	if endNum <= 0 {
		c.JSON(400, gin.H{"error": "endNum必须大于0"})
		return
	}

	// 生成URL列表
	urls := make([]string, endNum)
	for i := range urls {
		urls[i] = fmt.Sprintf(spiderUrl, i+1)
	}

	// 获取网站配置
	config, err := api.executor.configLoader.GetWebsiteConfig(website)
	if err != nil {
		log.Errorf("获取网站配置失败: %v", err)
		c.JSON(500, gin.H{"error": "获取网站配置失败"})
		return
	}

	// 解析其他参数
	params := make(map[string]interface{})
	params["websiteId"] = int(gjson.Get(string(data), "websiteId").Int())
	params["pornTypeId"] = int(gjson.Get(string(data), "pornTypeId").Int())
	params["countryId"] = int(gjson.Get(string(data), "countryId").Int())
	params["typeId"] = int(gjson.Get(string(data), "typeId").Int())
	params["processId"] = int(gjson.Get(string(data), "processId").Int())
	params["authorConcatType"] = int(gjson.Get(string(data), "authorConcatType").Int())

	// 验证业务逻辑
	// websiteId必须 能从数据库 查到 --
	if websiteId <= 0 {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, websiteId必须大于0"})
		return
	}

	// 使用 DBFindOneByField 查询 website 表，验证 websiteId 是否存在 --
	websiteRecord, err := db.DBFindOneByField[models.Website]("id", websiteId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": fmt.Sprintf("func=爬一本书chapter 失败, websiteId=%d 在数据库中不存在, 请先创建一条数据", websiteId)})
		} else {
			log.Errorf("func=爬一本书chapter 失败, 查询website失败: %v", err)
			log.Error("func=爬一本书chapter 失败, 查询website = ", websiteRecord)
			c.JSON(500, gin.H{"error": "func=爬一本书chapter 失败, 查询website失败"})
		}
		return
	}

	// 设置target参数，指定爬取场景
	params["target"] = "one_type_all_book" // 要求跟配置文件,里 一致

	// 从配置中提取选择器参数（v1版本需要的参数）
	if config.Crawl.Selectors != nil {
		if oneTypeAllBook, ok := config.Crawl.Selectors["one_type_all_book"].(map[string]any); ok {
			if arrSelector, ok := oneTypeAllBook["arr"].(string); ok {
				params["bookArrCssSelector"] = arrSelector
			}
			if itemSelector, ok := oneTypeAllBook["item"].(string); ok {
				params["bookArrItemCssSelector"] = itemSelector
			}
		}
	}

	log.Infof("开始配置驱动爬取: website=%s, urls=%d, target=%s", website, len(urls), params["target"])

	// 执行爬虫流程
	result := api.executor.Execute(website, data, urls, params)

	// 返回结果
	if result.Success {
		c.JSON(200, gin.H{
			"message":    result.Message,
			"processed":  result.ProcessedCount,
			"db_results": result.DBResults,
		})
	} else {
		c.JSON(500, gin.H{
			"error":      "爬取失败",
			"details":    result.ErrorDetails,
			"db_results": result.DBResults,
		})
	}
}

// DispatchApi_OneChapterAllChapterByHtml_ConfigDriven 配置驱动的章节爬取API
// 对应v1版的功能 DispatchApi_OneBOokAllChapterByHtml，但使用配置驱动
func (api *ConfigDrivenAPIController) DispatchApi_OneBookAllChapterByHtml_ConfigDriven(c *gin.Context) {
	/*
		待整理思路：
		1. 读取请求数据
		2. 解析请求参数
		3. 验证必需参数
		4. 生成URL
		5. 获取网站配置
	*/

	// 读取请求数据
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "读取请求数据失败: " + err.Error()})
		return
	}

	// 解析请求参数
	website := gjson.Get(string(data), "spiderTag.website").String()
	spiderUrl := gjson.Get(string(data), "spiderUrl").String()
	parentId := gjson.Get(string(data), "parentId").Int()

	// 验证必需参数
	if website == "" {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, 缺少必需参数: spiderTag.website"})
		return
	}
	if spiderUrl == "" {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, 缺少必需参数: spiderUrl"})
		return
	}
	if parentId <= 0 {
		c.JSON(400, gin.H{"error": "func=爬一本书chapter 失败, parentId必须大于0"})
		return
	}

	// 生成URL 列表
	urls := make([]string, 1)
	for i := range urls {
		urls[i] = fmt.Sprintf(spiderUrl, parentId)
	}

	// 获取网站配置
	config, err := api.executor.configLoader.GetWebsiteConfig(website)

	if err != nil {
		log.Errorf("获取网站配置失败: %v", err)
		c.JSON(500, gin.H{"error": "func=爬一本书chapter 失败, 获取网站配置失败"})
		return
	}

	// 解析其他参数
	params := make(map[string]any)
	// 设置target参数，指定爬取场景
	params["target"] = "one_book_all_chapter"

	// 从配置中提取选择器参数（v1版本需要的参数）
	if config.Crawl.Selectors != nil {
		if oneBookAllChapter, ok := config.Crawl.Selectors["one_book_all_chapter"].(map[string]any); ok {
			if arrSelector, ok := oneBookAllChapter["arr"].(string); ok {
				params["bookArrCssSelector"] = arrSelector
			}
			if itemSelector, ok := oneBookAllChapter["item"].(string); ok {
				params["bookArrItemCssSelector"] = itemSelector
			}
		}
	}
	// 注释的用不着
	// params["websiteId"] = gjson.Get(string(data), "websiteId").Int()
	// params["pornTypeId"] = gjson.Get(string(data), "pornTypeId").Int()
	// params["countryId"] = gjson.Get(string(data), "countryId").Int()
	// params["typeId"] = gjson.Get(string(data), "typeId").Int()
	// params["processId"] = gjson.Get(string(data), "processId").Int()
	// params["authorConcatType"] = gjson.Get(string(data), "authorConcatType").Int()

	return // 测试，这里还没写完-------------

	log.Infof("开始配置驱动爬取: website=%s, urls=%d, target=%s", website, len(urls), params["target"])

	// 执行爬虫流程
	result := api.executor.Execute(website, data, urls, params)

	// 返回结果
	if result.Success {
		c.JSON(200, gin.H{
			"message":    result.Message,
			"processed":  result.ProcessedCount,
			"db_results": result.DBResults,
		})
	} else {
		c.JSON(500, gin.H{
			"error":      "爬取失败",
			"details":    result.ErrorDetails,
			"db_results": result.DBResults,
		})
	}
}

// GetSupportedWebsites 获取支持的网站列表
func (api *ConfigDrivenAPIController) GetSupportedWebsites(c *gin.Context) {
	loader := config.GetSpiderConfigLoader()
	websites := loader.GetAllWebsites()

	c.JSON(200, gin.H{
		"websites": websites,
		"count":    len(websites),
	})
}

// GetWebsiteConfig 获取网站配置信息（调试用）
func (api *ConfigDrivenAPIController) GetWebsiteConfig(c *gin.Context) {
	website := c.Query("website")
	if website == "" {
		c.JSON(400, gin.H{"error": "缺少参数: website"})
		return
	}

	loader := config.GetSpiderConfigLoader()
	config, err := loader.GetWebsiteConfig(website)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"website": website,
		"config":  config,
	})
}

// ValidateConfig 验证配置
func (api *ConfigDrivenAPIController) ValidateConfig(c *gin.Context) {
	loader := config.GetSpiderConfigLoader()
	err := loader.ValidateConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "配置验证失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "配置验证通过"})
}

// RegisterRoutes 注册路由（在main.go中调用）
func (apiC *ConfigDrivenAPIController) RegisterRoutes(router *gin.Engine) {
	// 配置驱动的API路由 V2版本
	v2 := router.Group("/api/v2/spider")
	{
		v2.POST("/oneTypeAllBookByHtml/config", apiC.DispatchApi_OneTypeAllBookByHtml_ConfigDriven)
		v2.POST("/oneBookAllChapterByHtml/config", apiC.DispatchApi_OneBookAllChapterByHtml_ConfigDriven)
		v2.GET("/websites", apiC.GetSupportedWebsites)
		v2.GET("/config", apiC.GetWebsiteConfig)
		v2.POST("/validate", apiC.ValidateConfig)
	}
}

// 示例用法：
// 在main.go中添加：
//
//	api := spider.NewConfigDrivenAPIController()
//	err := api.InitConfig("spider-config.yaml")
//	if err != nil {
//		log.Fatal("初始化爬虫配置失败:", err)
//	}
//	api.RegisterRoutes(router)
//
// 然后前端可以调用：
// POST /api/v2/spider/oneTypeAllBookByHtml/config
//
// 请求体与原有API相同，但处理逻辑完全由配置驱动

// 兼容性说明：
// 1. V2 API的请求体格式与原有API完全相同
// 2. 响应格式略有不同，但包含更多详细信息
// 3. 原有API (/spider/oneTypeAllBookByHtml) 保持不变
// 4. 可以逐步迁移到V2 API，无缝切换
