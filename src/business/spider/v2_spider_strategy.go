// 爬虫策略工厂 - 配置驱动的爬取策略
// 支持HTML和JSON两种数据源的爬取

package spider

import (
	"fmt"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/models"
)

// 爬虫策略接口
type SpiderStrategy interface {
	Crawl(data []byte, config *config.WebsiteConfig, urls []string, params map[string]interface{}) (interface{}, error)
}

// HTML爬虫策略
type HtmlSpiderStrategy struct{}

// Crawl 实现HTML爬取
func (s *HtmlSpiderStrategy) Crawl(data []byte, config *config.WebsiteConfig, urls []string, params map[string]interface{}) (interface{}, error) {
	// 构建ComicSpider的mapping配置
	mapping := make(map[string]models.ModelHtmlMapping)

	for fieldName, fieldConfig := range config.Extract.Mappings {
		// 转换配置格式
		htmlMapping := models.ModelHtmlMapping{
			GetFieldPath: fieldConfig.Selector,
			GetHtmlType:  fieldConfig.Type,
			FiledType:    "string", // 默认string，后续可以扩展
		}

		// 添加Transform函数（这里暂时不添加，在后续步骤处理）
		// Transform会在字段映射器中处理

		mapping[fieldName] = htmlMapping
	}

	// 从配置中获取选择器
	bookArrCssSelector := ""     // 初始值
	bookArrItemCssSelector := "" // 初始值

	// 从params中获取target参数，确定使用哪个场景的选择器
	target := ""
	if params != nil {
		if t, ok := params["target"].(string); ok {
			target = t
		}
	}

	// 根据target选择对应的选择器配置
	if config.Crawl.Selectors != nil && target != "" {
		if scenarioConfig, exists := config.Crawl.Selectors[target]; exists {
			if scenarioMap, ok := scenarioConfig.(map[string]interface{}); ok {
				// 获取容器选择器（arr）
				if arrSelector, ok := scenarioMap["arr"].(string); ok {
					bookArrCssSelector = arrSelector
				}
				// 获取项目选择器（item）
				if itemSelector, ok := scenarioMap["item"].(string); ok {
					bookArrItemCssSelector = itemSelector
				}
			}
		}
	}

	// 调用v2版本的HTML爬取函数
	// 还要考虑通用爬取，分爬type \ 爬book \ 爬chapter 3种情况 ------------- 待办 ！！！！！！！！！！
	result := GetOneTypeAllBookUseCollyByMappingV2[models.ComicSpider](data, mapping, urls, bookArrCssSelector, bookArrItemCssSelector, config)
	return result, nil
}

// JSON爬虫策略
type JsonSpiderStrategy struct{}

// Crawl 实现JSON爬取
func (s *JsonSpiderStrategy) Crawl(data []byte, config *config.WebsiteConfig, urls []string, params map[string]interface{}) (interface{}, error) {
	// TODO: 实现JSON爬取策略
	// 需要创建类似HTML的JSON版本
	return nil, fmt.Errorf("JSON爬取策略暂未实现")
}

// 爬虫策略工厂
type SpiderStrategyFactory struct{}

// GetStrategy 根据配置类型获取策略
func (f *SpiderStrategyFactory) GetStrategy(crawlType string) (SpiderStrategy, error) {
	switch crawlType {
	case "html":
		return &HtmlSpiderStrategy{}, nil
	case "json":
		return &JsonSpiderStrategy{}, nil
	default:
		return nil, fmt.Errorf("不支持的爬取类型: %s", crawlType)
	}
}

// 默认工厂实例
var defaultFactory = &SpiderStrategyFactory{}

// GetSpiderStrategy 获取爬虫策略（简化接口）
func GetSpiderStrategy(crawlType string) (SpiderStrategy, error) {
	return defaultFactory.GetStrategy(crawlType)
}
