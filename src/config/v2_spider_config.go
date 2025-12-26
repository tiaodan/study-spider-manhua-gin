// 爬虫配置系统 - 配置驱动架构
// 不修改任何现有代码，创建全新的配置驱动实现

package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v3"
)

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
