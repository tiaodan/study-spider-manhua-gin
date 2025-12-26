/*
功能: v2-spider-config.yaml 配置文件 仅结构体, 不提供get/set方法。 V1.5实现
1. 提供 2种获取配置文件 tag方式
  - yaml.v3 库 + yaml tag, 实现起来较简单，不能处理复杂情况: 比如命令行替换 某个配置. 核心代码：yaml.Unmarshal()
  - viper库 + mapstructure tag, 实现起来较复杂.能处理复杂情况: 比如命令行替换 某个配置. 核心代码：viper.Unmarshal()
  - 一般推荐2种写法都用，用户想用哪种就用哪种。聪明人不做选择，我都要!!!
*/

package config

// 配置结构体定义, v2-spider-config.yaml
type SpiderConfig struct {
	Websites     map[string]*WebsiteConfig `yaml:"websites" mapstructure:"websites" `
	TransformLib map[string]*TransformDef  `yaml:"transform_library" mapstructure:"transform_library" `
}

type WebsiteConfig struct {
	Meta          *MetaConfig                    `yaml:"meta" mapstructure:"meta" `
	Crawl         *CrawlConfig                   `yaml:"crawl" mapstructure:"crawl"`
	Extract       *ExtractConfig                 `yaml:"extract" mapstructure:"extract"`
	Clean         *CleanConfig                   `yaml:"clean" mapstructure:"clean"`
	Validate      *ValidateConfig                `yaml:"validate" mapstructure:"validate"`
	Insert        *InsertConfig                  `yaml:"insert" mapstructure:"insert"`
	RelatedTables map[string]*RelatedTableConfig `yaml:"related_tables" mapstructure:"related_tables"` // 关联表配置
}

type MetaConfig struct {
	Name    string `yaml:"name" mapstructure:"name"`
	BaseURL string `yaml:"base_url" mapstructure:"base_url"`
	Table   string `yaml:"table" mapstructure:"table"`
}

type CrawlConfig struct {
	Type      string         `yaml:"type" mapstructure:"type"`           // html/json/xml/api
	Selectors map[string]any `yaml:"selectors" mapstructure:"selectors"` // HTML选择器，支持嵌套结构
	DataPath  string         `yaml:"data_path" mapstructure:"data_path"` // JSON数据路径
}

type ExtractConfig struct {
	Mappings map[string]*FieldMapping `yaml:"mappings" mapstructure:"mappings"`
}

type FieldMapping struct {
	Selector   string   `yaml:"selector" mapstructure:"selector"`     // HTML选择器
	Path       string   `yaml:"path" mapstructure:"path"`             // JSON路径
	Type       string   `yaml:"type" mapstructure:"type"`             // content/attr
	Transforms []string `yaml:"transforms" mapstructure:"transforms"` // Transform函数名列表
}

type CleanConfig struct {
	ForeignKeys map[string]string `yaml:"foreign_keys" mapstructure:"foreign_keys"` // 外键字段映射
	Defaults    map[string]any    `yaml:"defaults" mapstructure:"defaults"`         // 默认值
}

type ValidateConfig struct {
	Rules map[string][]*ValidateRule `yaml:"rules" mapstructure:"rules"`
}

type ValidateRule struct {
	Name   string         `yaml:"name" mapstructure:"name"`     // 验证器名称
	Params map[string]any `yaml:"params" mapstructure:"params"` // 验证器参数
}

type InsertConfig struct {
	Strategy   string   `yaml:"strategy" mapstructure:"strategy"` // insert/update/upsert
	UniqueKeys []string `yaml:"unique_keys" mapstructure:"unique_keys"`
	UpdateKeys []string `yaml:"update_keys" mapstructure:"update_keys"`
}

// 关联表配置
type RelatedTableConfig struct {
	Table      string        `yaml:"table" mapstructure:"table"`             // 表名
	Source     string        `yaml:"source" mapstructure:"source"`           // 数据来源：main/field
	SourcePath string        `yaml:"source_path" mapstructure:"source_path"` // 数据来源路径，如 "stats" 表示从主表的 Stats 字段获取
	Insert     *InsertConfig `yaml:"insert" mapstructure:"insert"`           // 插入配置
}

type TransformDef struct {
	Type        string         `yaml:"type" mapstructure:"type"` // string/number/enum/validator
	Description string         `yaml:"description" mapstructure:"description"`
	Params      map[string]any `yaml:"params" mapstructure:"params"` // 默认参数
}
