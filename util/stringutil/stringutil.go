// 处理字符串工具
package stringutil

import "strings"

// 处理model对象, 如 models.Website。预留，暂时没用
func TrimSpaceObj(nouse string) {
}

// 处理 map[string]interface{} 这种键值对 对象。obj原始对象数据，会被修改
func TrimSpaceMap(obj map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	for key, value := range obj {
		if str, ok := value.(string); ok {
			obj[key] = strings.TrimSpace(str) // 去除首尾空格
		}
	}
}
