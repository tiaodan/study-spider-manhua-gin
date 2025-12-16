// Transform函数库 - 配置驱动的数据转换
// 支持参数化的数据转换函数

package spider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/util/langutil"
)

// Transform函数接口
type TransformFunc func(value interface{}, params map[string]interface{}) (interface{}, error)

// Transform注册器
type TransformRegistry struct {
	transforms map[string]TransformFunc
	validators map[string]TransformFunc
}

// 单例实例
var registryInstance *TransformRegistry
var registryOnce sync.Once

// GetTransformRegistry 获取Transform注册器单例
func GetTransformRegistry() *TransformRegistry {
	registryOnce.Do(func() {
		registryInstance = &TransformRegistry{
			transforms: make(map[string]TransformFunc),
			validators: make(map[string]TransformFunc),
		}
		registryInstance.registerBuiltInTransforms()
	})
	return registryInstance
}

// RegisterTransform 注册转换函数
func (r *TransformRegistry) RegisterTransform(name string, fn TransformFunc) {
	r.transforms[name] = fn
}

// RegisterValidator 注册验证函数
func (r *TransformRegistry) RegisterValidator(name string, fn TransformFunc) {
	r.validators[name] = fn
}

// ApplyTransform 应用单个转换
func (r *TransformRegistry) ApplyTransform(name string, value interface{}, params map[string]interface{}) (interface{}, error) {
	fn, exists := r.transforms[name]
	if !exists {
		return value, fmt.Errorf("未找到转换函数: %s", name)
	}
	return fn(value, params)
}

// ApplyTransforms 应用转换管道
func (r *TransformRegistry) ApplyTransforms(transformNames []string, value interface{}, configLoader *config.SpiderConfigLoader) (interface{}, error) {
	result := value
	for _, name := range transformNames {
		// 从配置中获取Transform定义和参数
		transformDef, err := configLoader.GetTransformDef(name)
		if err != nil {
			return result, fmt.Errorf("获取Transform定义失败 %s: %v", name, err)
		}

		// 应用转换
		result, err = r.ApplyTransform(name, result, transformDef.Params)
		if err != nil {
			return result, fmt.Errorf("应用Transform %s 失败: %v", name, err)
		}
	}
	return result, nil
}

// Validate 应用验证规则
func (r *TransformRegistry) Validate(ruleName string, value interface{}, params map[string]interface{}) error {
	validator, exists := r.validators[ruleName]
	if !exists {
		return fmt.Errorf("未找到验证器: %s", ruleName)
	}
	_, err := validator(value, params)
	return err
}

// toInt 安全地将interface{}转换为int，支持int和float64两种类型
// YAML解析时，整数可能被解析为int或float64，需要统一处理
func toInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	default:
		return 0
	}
}

// 注册内置转换函数
func (r *TransformRegistry) registerBuiltInTransforms() {
	// 字符串处理转换器
	r.RegisterTransform("trim_space", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			return strings.TrimSpace(str), nil
		}
		return value, nil
	})

	r.RegisterTransform("simplify_chinese", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			// ignoreDigits参数暂时保留，未来可能用于控制数字转换
			// 可以通过 params["ignore_digits"] 访问，目前未使用

			result, err := langutil.TraditionalToSimplified(str)
			if err != nil {
				return str, fmt.Errorf("繁体转简体失败: %v", err)
			}

			// 手动替换一些特殊字符（作为opencc的补充）
			result = strings.ReplaceAll(result, "姊", "姐")

			return result, nil
		}
		return value, nil
	})

	r.RegisterTransform("regex_extract", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			pattern := `(.+)` // 默认提取全部
			if val, exists := params["pattern"]; exists {
				pattern = val.(string)
			}

			group := 1 // 默认提取第一个分组
			if val, exists := params["group"]; exists {
				group = toInt(val) // YAML解析后可能是int或float64
			}

			re, err := regexp.Compile(pattern)
			if err != nil {
				return value, fmt.Errorf("正则表达式编译失败: %v", err)
			}

			matches := re.FindStringSubmatch(str)
			if len(matches) > group {
				return matches[group], nil
			}
			return str, nil // 提取失败返回原值
		}
		return value, nil
	})

	r.RegisterTransform("remove_domain_prefix", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			prefix := "https://" // 默认前缀
			if val, exists := params["prefix"]; exists {
				prefix = val.(string)
			}

			if strings.HasPrefix(str, prefix) {
				return strings.TrimPrefix(str, prefix), nil
			}
		}
		return value, nil
	})

	r.RegisterTransform("add_prefix", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			prefix := "" // 默认前缀
			if val, exists := params["prefix"]; exists {
				prefix = val.(string)
			}
			return prefix + str, nil
		}
		return value, nil
	})

	// V1版本的URL拼接transform方法
	// 参考V1版本v1_spider_template_book_type.go:453-456行的Transform函数
	r.RegisterTransform("url_join", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			baseUrl := "" // 基础URL
			if val, exists := params["base_url"]; exists {
				baseUrl = val.(string)
			}
			// 如果baseUrl为空，尝试使用prefix参数（兼容add_prefix的用法）
			if baseUrl == "" {
				if val, exists := params["prefix"]; exists {
					baseUrl = val.(string)
				}
			}
			return baseUrl + str, nil
		}
		return value, nil
	})

	// 类型转换器
	r.RegisterTransform("to_int", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		switch v := value.(type) {
		case int:
			return v, nil
		case int64:
			return int(v), nil
		case float64:
			return int(v), nil
		case float32:
			return int(v), nil
		case string:
			result, err := strconv.Atoi(v)
			if err != nil {
				return 0, fmt.Errorf("无法转换为int: %v", err)
			}
			return result, nil
		default:
			return 0, fmt.Errorf("不支持的类型转换为int: %T", value)
		}
	})

	r.RegisterTransform("to_float", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		switch v := value.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case string:
			result, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0.0, fmt.Errorf("无法转换为float: %v", err)
			}
			return result, nil
		default:
			return 0.0, fmt.Errorf("不支持的类型转换为float: %T", value)
		}
	})

	r.RegisterTransform("to_string", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		return fmt.Sprintf("%v", value), nil
	})

	r.RegisterTransform("to_bool", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			result, err := strconv.ParseBool(v)
			if err != nil {
				return false, fmt.Errorf("无法转换为bool: %v", err)
			}
			return result, nil
		case int:
			return v != 0, nil
		case float64:
			return v != 0, nil
		default:
			return false, fmt.Errorf("不支持的类型转换为bool: %T", value)
		}
	})

	// 数据转换器
	// parse_hits_number: 解析点击数量字符串，支持V1版本的完整逻辑
	// 参考V1版本v1_spider_template_book_type.go:269-301行的处理逻辑
	r.RegisterTransform("parse_hits_number", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			// 1. 预处理：去空格，繁体转简体（V1版本的逻辑）
			hitsStr := strings.TrimSpace(str)
			hitsStrSimplified, err := langutil.TraditionalToSimplified(hitsStr)
			if err != nil {
				// 如果转换失败，使用原字符串
				hitsStrSimplified = hitsStr
			}

			// 2. 正则匹配：匹配数字和单位（支持中英文单位）
			// V1版本使用: `(\d+\.?\d*)\s*([^\d\s]+)` 匹配数字和单位
			re := regexp.MustCompile(`(\d+\.?\d*)\s*([^\d\s]+)`)
			matches := re.FindStringSubmatch(hitsStrSimplified)

			var hitsNumStr string // 数字部分
			var hitsUnit string   // 单位 如：万、千、亿
			numUnit := 1          // 单位倍数，默认1

			if len(matches) >= 3 {
				hitsNumStr = matches[1] // 匹配数字部分 如 95.2
				hitsUnit = matches[2]   // 匹配单位 如：万
				switch hitsUnit {
				case "亿":
					numUnit = 100000000
				case "万":
					numUnit = 10000
				case "千":
					numUnit = 1000
				case "w", "W":
					numUnit = 10000
				case "k", "K":
					numUnit = 1000
				}
			} else {
				// 重新正则匹配：只匹配数字（不带单位的情况）
				re = regexp.MustCompile(`(\d+\.?\d*)\s*`)
				newMatches := re.FindStringSubmatch(hitsStrSimplified)
				if len(newMatches) >= 2 {
					hitsNumStr = newMatches[1]
				}
			}

			// 3. 计算具体数字 HitsNum * hitsUnit
			if hitsNumStr == "" {
				return 0, nil
			}

			hitsFloat, err := strconv.ParseFloat(hitsNumStr, 64)
			if err != nil || hitsFloat < 0 {
				// 错误或负值设为0
				return 0, nil
			}

			return int(hitsFloat * float64(numUnit)), nil
		}
		return value, nil
	})

	r.RegisterTransform("map_end_status", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			mapping := map[string]int{
				"完结": 3,
				"连载": 2,
			}

			// 合并配置中的映射
			if val, exists := params["mapping"]; exists {
				if configMapping, ok := val.(map[string]interface{}); ok {
					for k, v := range configMapping {
						mapping[k] = toInt(v)
					}
				}
			}

			if status, exists := mapping[str]; exists {
				return status, nil
			}

			defaultVal := 1
			if val, exists := params["default"]; exists {
				defaultVal = toInt(val)
			}
			return defaultVal, nil
		}
		return value, nil
	})

	r.RegisterTransform("map_completion_status", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		// 根据章节数判断完成状态
		completeThreshold := 1000
		if val, exists := params["thresholds"]; exists {
			if thresholds, ok := val.(map[string]interface{}); ok {
				if threshold, ok := thresholds["complete"]; ok {
					completeThreshold = toInt(threshold)
				}
			}
		}

		if num, ok := value.(float64); ok {
			if int(num) >= completeThreshold {
				return 3, nil // 完结
			}
			return 2, nil // 连载
		}
		return 1, nil // 未知
	})

	// 验证器
	r.RegisterValidator("not_empty", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			if strings.TrimSpace(str) == "" {
				return nil, fmt.Errorf("值不能为空")
			}
		}
		return value, nil
	})

	r.RegisterValidator("max_length", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		maxLen := 100 // 默认最大长度
		if val, exists := params["max"]; exists {
			maxLen = toInt(val)
		}

		if str, ok := value.(string); ok {
			if len(str) > maxLen {
				return nil, fmt.Errorf("长度超过最大值 %d", maxLen)
			}
		}
		return value, nil
	})

	r.RegisterValidator("valid_url", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		if str, ok := value.(string); ok {
			if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") && !strings.HasPrefix(str, "/") {
				return nil, fmt.Errorf("无效的URL格式")
			}
		}
		return value, nil
	})

	r.RegisterValidator("min_value", func(value interface{}, params map[string]interface{}) (interface{}, error) {
		minVal := 0 // 默认最小值
		if val, exists := params["min"]; exists {
			minVal = toInt(val)
		}

		if num, ok := value.(float64); ok {
			if int(num) < minVal {
				return nil, fmt.Errorf("值小于最小值 %d", minVal)
			}
		}
		return value, nil
	})
}
