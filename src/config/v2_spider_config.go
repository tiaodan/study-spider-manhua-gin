// 爬虫配置系统 - 配置驱动架构
// 不修改任何现有代码，创建全新的配置驱动实现

package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v3"
)

// 配置结构体定义
type SpiderConfig struct {
	Websites       map[string]*WebsiteConfig `yaml:"websites"`
	TransformLib   map[string]*TransformDef  `yaml:"transform_library"`
}

type WebsiteConfig struct {
	Meta     *MetaConfig     `yaml:"meta"`
	Crawl    *CrawlConfig    `yaml:"crawl"`
	Extract  *ExtractConfig  `yaml:"extract"`
	Clean    *CleanConfig    `yaml:"clean"`
	Validate *ValidateConfig `yaml:"validate"`
	Insert   *InsertConfig   `yaml:"insert"`
}

type MetaConfig struct {
	Name    string `yaml:"name"`
	BaseURL string `yaml:"base_url"`
	Table   string `yaml:"table"`
}

type CrawlConfig struct {
	Type      string            `yaml:"type"`       // html/json/xml/api
	Selectors map[string]string `yaml:"selectors"` // HTML选择器
	DataPath  string            `yaml:"data_path"`  // JSON数据路径
}

type ExtractConfig struct {
	Mappings map[string]*FieldMapping `yaml:"mappings"`
}

type FieldMapping struct {
	Selector   string                   `yaml:"selector"`   // HTML选择器
	Path       string                   `yaml:"path"`       // JSON路径
	Type       string                   `yaml:"type"`       // content/attr
	Transforms []string                 `yaml:"transforms"` // Transform函数名列表
}

type CleanConfig struct {
	ForeignKeys map[string]string `yaml:"foreign_keys"` // 外键字段映射
	Defaults    map[string]interface{} `yaml:"defaults"` // 默认值
}

type ValidateConfig struct {
	Rules map[string][]*ValidateRule `yaml:"rules"`
}

type ValidateRule struct {
	Name   string                 `yaml:"name"`   // 验证器名称
	Params map[string]interface{} `yaml:"params"` // 验证器参数
}

type InsertConfig struct {
	Strategy   string   `yaml:"strategy"`    // insert/update/upsert
	UniqueKeys []string `yaml:"unique_keys"`
	UpdateKeys []string `yaml:"update_keys"`
}

type TransformDef struct {
	Type        string                 `yaml:"type"`        // string/number/enum/validator
	Description string                 `yaml:"description"`
	Params      map[string]interface{} `yaml:"params"` // 默认参数
}

// 配置加载器
type SpiderConfigLoader struct {
	config *SpiderConfig
	once   sync.Once
}

var loaderInstance *SpiderConfigLoader
var loaderOnce sync.Once

// GetSpiderConfigLoader 获取配置加载器单例
func GetSpiderConfigLoader() *SpiderConfigLoader {
	loaderOnce.Do(func() {
		loaderInstance = &SpiderConfigLoader{}
	})
	return loaderInstance
}

// LoadConfig 加载配置文件
func (l *SpiderConfigLoader) LoadConfig(configPath string) error {
	var err error
	l.once.Do(func() {
		// 读取YAML文件
		data, readErr := ioutil.ReadFile(configPath)
		if readErr != nil {
			err = fmt.Errorf("读取配置文件失败: %v", readErr)
			return
		}

		// 解析YAML
		config := &SpiderConfig{}
		if parseErr := yaml.Unmarshal(data, config); parseErr != nil {
			err = fmt.Errorf("解析配置文件失败: %v", parseErr)
			return
		}

		l.config = config
	})

	return err
}

// GetWebsiteConfig 获取网站配置
func (l *SpiderConfigLoader) GetWebsiteConfig(website string) (*WebsiteConfig, error) {
	if l.config == nil {
		return nil, fmt.Errorf("配置未加载")
	}

	config, exists := l.config.Websites[website]
	if !exists {
		return nil, fmt.Errorf("未找到网站配置: %s", website)
	}

	return config, nil
}

// GetTransformDef 获取Transform定义
func (l *SpiderConfigLoader) GetTransformDef(name string) (*TransformDef, error) {
	if l.config == nil {
		return nil, fmt.Errorf("配置未加载")
	}

	def, exists := l.config.TransformLib[name]
	if !exists {
		return nil, fmt.Errorf("未找到Transform定义: %s", name)
	}

	return def, nil
}

// GetAllWebsites 获取所有支持的网站列表
func (l *SpiderConfigLoader) GetAllWebsites() []string {
	if l.config == nil {
		return nil
	}

	websites := make([]string, 0, len(l.config.Websites))
	for name := range l.config.Websites {
		websites = append(websites, name)
	}
	return websites
}

// ValidateConfig 验证配置完整性
func (l *SpiderConfigLoader) ValidateConfig() error {
	if l.config == nil {
		return fmt.Errorf("配置未加载")
	}

	// 检查必需的配置项
	for name, website := range l.config.Websites {
		if website.Crawl == nil {
			return fmt.Errorf("网站 %s 缺少crawl配置", name)
		}
		if website.Extract == nil {
			return fmt.Errorf("网站 %s 缺少extract配置", name)
		}
		if website.Insert == nil {
			return fmt.Errorf("网站 %s 缺少insert配置", name)
		}
	}

	return nil
}
