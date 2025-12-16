// 数据验证器 - 配置驱动的数据验证
// 支持多种验证规则的批量验证

package spider

import (
	"fmt"
	"reflect"
	"strings"

	"study-spider-manhua-gin/src/config"
)

// 验证结果
type ValidationResult struct {
	IsValid bool
	Errors  map[string][]string `json:"errors"` // field -> errors
}

// 数据验证器
type DataValidator struct {
	transformRegistry *TransformRegistry
}

// NewDataValidator 创建数据验证器
func NewDataValidator() *DataValidator {
	return &DataValidator{
		transformRegistry: GetTransformRegistry(),
	}
}

// Validate 验证单个对象
func (v *DataValidator) Validate(data interface{}, rules map[string][]*config.ValidateRule) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  make(map[string][]string),
	}

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	if dataValue.Kind() != reflect.Struct {
		result.IsValid = false
		result.Errors["general"] = []string{"验证对象必须是结构体或结构体指针"}
		return result
	}

	// 遍历所有验证规则
	for fieldName, fieldRules := range rules {
		// 兼容配置里的小写/下划线字段名，尝试转换为导出字段名
		targetFieldName, fieldValue, found := v.findField(dataValue, fieldName)
		if !found {
			result.Errors[fieldName] = []string{fmt.Sprintf("字段 %s 不存在", fieldName)}
			result.IsValid = false
			continue
		}
		fieldInterfaceValue := fieldValue.Interface()

		// 应用每个验证规则
		for _, rule := range fieldRules {
			err := v.transformRegistry.Validate(rule.Name, fieldInterfaceValue, rule.Params)
			if err != nil {
				result.IsValid = false
				if result.Errors[targetFieldName] == nil {
					result.Errors[targetFieldName] = []string{}
				}
				result.Errors[targetFieldName] = append(result.Errors[targetFieldName], err.Error())
			}
		}
	}

	return result
}

// findField 根据配置字段名查找结构体字段，支持小写/下划线到导出字段的转换
func (v *DataValidator) findField(dataValue reflect.Value, fieldName string) (string, reflect.Value, bool) {
	// direct match
	if f, ok := dataValue.Type().FieldByName(fieldName); ok {
		return f.Name, dataValue.FieldByName(f.Name), true
	}

	// convert snake_case or lowerCamel to Camel (exported)
	camel := toExportedCamel(fieldName)
	if f, ok := dataValue.Type().FieldByName(camel); ok {
		return f.Name, dataValue.FieldByName(f.Name), true
	}

	// 尝试嵌套Stats（例如 hits 在 Stats.Hits）
	if statsField, ok := dataValue.Type().FieldByName("Stats"); ok {
		statsVal := dataValue.FieldByName(statsField.Name)
		statsType := statsVal.Type()

		if f, ok := statsType.FieldByName(fieldName); ok {
			return "Stats." + f.Name, statsVal.FieldByName(f.Name), true
		}
		if f, ok := statsType.FieldByName(camel); ok {
			return "Stats." + f.Name, statsVal.FieldByName(f.Name), true
		}
	}

	return fieldName, reflect.Value{}, false
}

// toExportedCamel 将 name 或 comic_url_api_path 转换为 Name、ComicUrlApiPath
func toExportedCamel(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	joined := strings.Join(parts, "")
	if joined == "" {
		return name
	}
	// 已经是小驼峰的情况（没有下划线），也会变成首字母大写
	return strings.ToUpper(joined[:1]) + joined[1:]
}

// ValidateBatch 批量验证
func (v *DataValidator) ValidateBatch(dataList interface{}, rules map[string][]*config.ValidateRule) []*ValidationResult {
	results := []*ValidationResult{}

	dataListValue := reflect.ValueOf(dataList)
	if dataListValue.Kind() != reflect.Slice {
		result := &ValidationResult{
			IsValid: false,
			Errors: map[string][]string{
				"general": {"输入必须是切片类型"},
			},
		}
		return []*ValidationResult{result}
	}

	for i := 0; i < dataListValue.Len(); i++ {
		item := dataListValue.Index(i).Interface()
		result := v.Validate(item, rules)
		results = append(results, result)
	}

	return results
}

// ValidateWithConfig 根据配置验证
func (v *DataValidator) ValidateWithConfig(data interface{}, config *config.WebsiteConfig) *ValidationResult {
	if config.Validate == nil {
		// 没有验证配置，默认为有效
		return &ValidationResult{IsValid: true, Errors: make(map[string][]string)}
	}

	return v.Validate(data, config.Validate.Rules)
}

// ValidateBatchWithConfig 批量验证（带配置）
func (v *DataValidator) ValidateBatchWithConfig(dataList interface{}, config *config.WebsiteConfig) []*ValidationResult {
	if config.Validate == nil {
		// 没有验证配置，默认为有效
		count := 0
		if dataListValue := reflect.ValueOf(dataList); dataListValue.Kind() == reflect.Slice {
			count = dataListValue.Len()
		}

		results := make([]*ValidationResult, count)
		for i := range results {
			results[i] = &ValidationResult{IsValid: true, Errors: make(map[string][]string)}
		}
		return results
	}

	return v.ValidateBatch(dataList, config.Validate.Rules)
}

// FilterValid 过滤出有效的对象
func (v *DataValidator) FilterValid(dataList interface{}, validationResults []*ValidationResult) interface{} {
	dataListValue := reflect.ValueOf(dataList)
	if dataListValue.Kind() != reflect.Slice {
		return dataList
	}

	validSlice := reflect.MakeSlice(dataListValue.Type(), 0, dataListValue.Len())

	for i := 0; i < dataListValue.Len() && i < len(validationResults); i++ {
		if validationResults[i].IsValid {
			validSlice = reflect.Append(validSlice, dataListValue.Index(i))
		}
	}

	return validSlice.Interface()
}

// GetValidationSummary 获取验证摘要
func (v *DataValidator) GetValidationSummary(results []*ValidationResult) map[string]int {
	summary := map[string]int{
		"total":   len(results),
		"valid":   0,
		"invalid": 0,
		"errors":  0,
	}

	for _, result := range results {
		if result.IsValid {
			summary["valid"]++
		} else {
			summary["invalid"]++
			for _, fieldErrors := range result.Errors {
				summary["errors"] += len(fieldErrors)
			}
		}
	}

	return summary
}
