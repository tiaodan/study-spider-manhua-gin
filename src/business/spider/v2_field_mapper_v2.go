// 字段映射器 - 配置驱动的字段提取和转换
// 支持HTML和JSON数据源的字段映射

package spider

import (
	"fmt"
	"reflect"
	"strconv"

	"study-spider-manhua-gin/src/config"
)

// 数据提取器接口
type DataExtractor interface {
	Extract(selector string, dataType string) (interface{}, error)
}

// HTML数据提取器
type HtmlDataExtractor struct {
	element interface{} // *colly.HTMLElement 或其他HTML元素
}

// Extract 从HTML中提取数据
func (e *HtmlDataExtractor) Extract(selector string, dataType string) (interface{}, error) {
	// TODO: 实现HTML数据提取
	// 这里需要调用现有的HTML解析逻辑
	return nil, fmt.Errorf("HTML数据提取暂未实现")
}

// JSON数据提取器
type JsonDataExtractor struct {
	data map[string]interface{}
}

// Extract 从JSON中提取数据
func (e *JsonDataExtractor) Extract(path string, dataType string) (interface{}, error) {
	// TODO: 实现JSON路径提取
	// 这里需要实现JSON路径解析逻辑
	return nil, fmt.Errorf("JSON数据提取暂未实现")
}

// 字段映射器
type FieldMapper struct {
	transformRegistry *TransformRegistry
	configLoader      *config.SpiderConfigLoader
}

// NewFieldMapper 创建字段映射器
func NewFieldMapper() *FieldMapper {
	return &FieldMapper{
		transformRegistry: GetTransformRegistry(),
		configLoader:      config.GetSpiderConfigLoader(),
	}
}

// MapFields 映射字段并应用转换
func (m *FieldMapper) MapFields(extractor DataExtractor, fieldMappings map[string]*config.FieldMapping, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target必须是指向结构体的指针")
	}

	targetElem := targetValue.Elem()
	targetType := targetElem.Type()

	for fieldName, mapping := range fieldMappings {
		// 提取原始值
		var rawValue interface{}
		var err error

		if mapping.Selector != "" {
			// HTML选择器方式
			rawValue, err = extractor.Extract(mapping.Selector, mapping.Type)
		} else if mapping.Path != "" {
			// JSON路径方式
			rawValue, err = extractor.Extract(mapping.Path, "json")
		} else {
			continue // 跳过没有选择器或路径的字段
		}

		if err != nil {
			return fmt.Errorf("提取字段 %s 失败: %v", fieldName, err)
		}

		// 应用Transform管道
		processedValue, err := m.transformRegistry.ApplyTransforms(mapping.Transforms, rawValue, m.configLoader)
		if err != nil {
			return fmt.Errorf("转换字段 %s 失败: %v", fieldName, err)
		}

		// 设置到目标结构体
		err = m.setFieldValue(targetElem, targetType, fieldName, processedValue)
		if err != nil {
			return fmt.Errorf("设置字段 %s 值失败: %v", fieldName, err)
		}
	}

	return nil
}

// setFieldValue 设置结构体字段值
func (m *FieldMapper) setFieldValue(targetElem reflect.Value, targetType reflect.Type, fieldName string, value interface{}) error {
	field, found := targetType.FieldByName(fieldName)
	if !found {
		return fmt.Errorf("字段 %s 不存在", fieldName)
	}

	fieldValue := targetElem.FieldByName(fieldName)
	if !fieldValue.CanSet() {
		return fmt.Errorf("字段 %s 不可设置", fieldName)
	}

	// 类型转换
	convertedValue, err := m.convertValue(value, field.Type)
	if err != nil {
		return fmt.Errorf("类型转换失败: %v", err)
	}

	fieldValue.Set(reflect.ValueOf(convertedValue))
	return nil
}

// convertValue 类型转换
func (m *FieldMapper) convertValue(value interface{}, targetType reflect.Type) (interface{}, error) {
	valueType := reflect.TypeOf(value)

	// 如果类型匹配，直接返回
	if valueType.AssignableTo(targetType) {
		return value, nil
	}

	// 处理常见类型转换
	switch targetType.Kind() {
	case reflect.String:
		if value == nil {
			return "", nil
		}
		return fmt.Sprintf("%v", value), nil

	case reflect.Int, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case float64:
			return int(v), nil
		case int:
			return v, nil
		case string:
			// 尝试转换字符串为数字
			if intVal, err := strconv.Atoi(v); err == nil {
				return intVal, nil
			}
			return 0, fmt.Errorf("无法将字符串转换为整数: %s", v)
		}

	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case string:
			if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
				return floatVal, nil
			}
			return 0.0, fmt.Errorf("无法将字符串转换为浮点数: %s", v)
		}

	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			return v != "", nil
		case int:
			return v != 0, nil
		}
	}

	return nil, fmt.Errorf("不支持的类型转换: %v -> %v", valueType, targetType)
}

// MapSingleField 映射单个字段（用于特殊情况）
func (m *FieldMapper) MapSingleField(extractor DataExtractor, fieldName string, mapping *config.FieldMapping) (interface{}, error) {
	// 提取原始值
	var rawValue interface{}
	var err error

	if mapping.Selector != "" {
		rawValue, err = extractor.Extract(mapping.Selector, mapping.Type)
	} else if mapping.Path != "" {
		rawValue, err = extractor.Extract(mapping.Path, "json")
	} else {
		return nil, fmt.Errorf("字段 %s 没有选择器或路径", fieldName)
	}

	if err != nil {
		return nil, err
	}

	// 应用Transform管道
	return m.transformRegistry.ApplyTransforms(mapping.Transforms, rawValue, m.configLoader)
}
