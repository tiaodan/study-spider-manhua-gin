/*
功能: v2-spider-config.yaml 配置文件 仅结构体, 不提供get/set方法。 V1.5实现
1. 提供 2种获取配置文件 tag方式
  - yaml.v3 库 + yaml tag, 实现起来较简单，不能处理复杂情况: 比如命令行替换 某个配置. 核心代码：yaml.Unmarshal()
  - viper库 + mapstructure tag, 实现起来较复杂.能处理复杂情况: 比如命令行替换 某个配置. 核心代码：viper.Unmarshal()
  - 一般推荐2种写法都用，用户想用哪种就用哪种。聪明人不做选择，我都要!!!
*/
package config

// 配置结构体定义, v2-spider-config.yaml
type SpiderConfigV15 struct {
	Websites     map[string]*WebsiteConfigV15 `yaml:"websites" mapstructure:"websites" `
	TransformLib map[string]*TransformDefV15  `yaml:"transform_library" mapstructure:"transform_library" `
}

// 网站配置结构体
type WebsiteConfigV15 struct {
	Stages map[string]*StageConfigV15 `yaml:"stages" mapstructure:"stages"` // 通用阶段配置
}

// 阶段结构体： 爬取某一类所有书籍 -> one_type_all_book
// 阶段结构体： 爬取某一本书所有章节 -> one_book_all_chapter
// 阶段结构体： 爬取某一章节所有内容 -> one_chapter_all_content
type StageConfigV15 struct {
	Meta              *MetaConfigV15                    `yaml:"meta" mapstructure:"meta" `
	Crawl             *CrawlConfigV15                   `yaml:"crawl" mapstructure:"crawl"`
	Extract           *ExtractConfigV15                 `yaml:"extract" mapstructure:"extract"`
	Clean             *CleanConfigV15                   `yaml:"clean" mapstructure:"clean"`
	Validate          *ValidateConfigV15                `yaml:"validate" mapstructure:"validate"`
	Insert            *InsertConfigV15                  `yaml:"insert" mapstructure:"insert"`
	RelatedTables     map[string]*RelatedTableConfigV15 `yaml:"related_tables" mapstructure:"related_tables"`           // 关联表配置
	UpdateParentStats *UpdateParentStatsConifgV15       `yaml:"update_parent_stats" mapstructure:"update_parent_stats"` // 更新 父stats 统计表
}

type MetaConfigV15 struct {
	Name    string `yaml:"name" mapstructure:"name"`
	BaseURL string `yaml:"base_url" mapstructure:"base_url"`
	Table   string `yaml:"table" mapstructure:"table"`
}

type CrawlConfigV15 struct {
	Type      string         `yaml:"type" mapstructure:"type"`           // html/json/xml/api
	Selectors map[string]any `yaml:"selectors" mapstructure:"selectors"` // HTML选择器，支持嵌套结构
	DataPath  string         `yaml:"data_path" mapstructure:"data_path"` // JSON数据路径
}

type ExtractConfigV15 struct {
	Mappings map[string]*FieldMapping `yaml:"mappings" mapstructure:"mappings"`
}

type FieldMappingV15 struct {
	Selector   string   `yaml:"selector" mapstructure:"selector"`     // HTML选择器
	Path       string   `yaml:"path" mapstructure:"path"`             // JSON路径
	Type       string   `yaml:"type" mapstructure:"type"`             // content/attr
	Transforms []string `yaml:"transforms" mapstructure:"transforms"` // Transform函数名列表
}

type CleanConfigV15 struct {
	ForeignKeys map[string]string `yaml:"foreign_keys" mapstructure:"foreign_keys"` // 外键字段映射
	Defaults    map[string]any    `yaml:"defaults" mapstructure:"defaults"`         // 默认值
}

type ValidateConfigV15 struct {
	Rules map[string][]*ValidateRule `yaml:"rules" mapstructure:"rules"`
}

type ValidateRuleV15 struct {
	Name   string         `yaml:"name" mapstructure:"name"`     // 验证器名称
	Params map[string]any `yaml:"params" mapstructure:"params"` // 验证器参数
}

type InsertConfigV15 struct {
	Strategy   string   `yaml:"strategy" mapstructure:"strategy"` // insert/update/upsert
	UniqueKeys []string `yaml:"unique_keys" mapstructure:"unique_keys"`
	UpdateKeys []string `yaml:"update_keys" mapstructure:"update_keys"`
}

// 关联表配置
type RelatedTableConfigV15 struct {
	Table      string           `yaml:"table" mapstructure:"table"`             // 表名
	Source     string           `yaml:"source" mapstructure:"source"`           // 数据来源：main/field
	SourcePath string           `yaml:"source_path" mapstructure:"source_path"` // 数据来源路径，如 "stats" 表示从主表的 Stats 字段获取
	Insert     *InsertConfigV15 `yaml:"insert" mapstructure:"insert"`           // 插入配置
}

type TransformDefV15 struct {
	Type        string         `yaml:"type" mapstructure:"type"` // string/number/enum/validator
	Description string         `yaml:"description" mapstructure:"description"`
	Params      map[string]any `yaml:"params" mapstructure:"params"` // 默认参数
}

// 更新 父stats 统计表 UpdateParentStatsConifg
type UpdateParentStatsConifgV15 struct {
	Strategy   string   `yaml:"strategy" mapstructure:"strategy"` // insert/update/upsert
	UniqueKeys []string `yaml:"unique_keys" mapstructure:"unique_keys"`
	UpdateKeys []string `yaml:"update_keys" mapstructure:"update_keys"`
}
