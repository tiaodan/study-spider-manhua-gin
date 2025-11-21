// 处理字符串工具
package stringutil

import "strings"

// 定义一个接口，约定所有可清理的对象都要实现这个方法
type TrimAble interface {
	TrimSpaces()
}

// 定义一个接口，约定所有 需要繁体转简体 对象都要实现这个方法
type SimpleAble interface {
	Trad2Simple()
}

// 处理model对象, 前后空格。通用入口：调用实现了接口的对象的 TrimSpaces 方法
/*
核心逻辑：

注意：
处理对象 string字段，前后空格的 工具方法
	一般有2种实现方式：
	1. 用反射实现，缺点：反射比较耗时，写起来麻烦，不容易理解 -- 弃用
	2. 用接口实现，util里 定义一个结果，所有models里的 struct 都实现这个接口 -- 推荐，现在用的这种方法，参考stringutils.go
*/
// TrimSpaceObj 是一个函数，用于处理实现了TrimAble接口的对象的空白字符
// 参数:
//   obj: TrimAble接口类型，表示需要去除空白字符的对象
func TrimSpaceObj(obj TrimAble) {
	obj.TrimSpaces()
}

// 处理model对象，繁体转简体
func SimpleObj(obj SimpleAble) {
	obj.Trad2Simple()
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
