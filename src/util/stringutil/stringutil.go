// 处理字符串工具
package stringutil

import (
	"regexp"
	"strconv"
	"strings"
	"study-spider-manhua-gin/src/util/langutil"
)

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

// 通过正则去除 http 头前缀. http:// https:// 都能去除
func TrimHttpPrefix(str string) string {
	// 匹配并移除 http:// 或 https:// 前缀
	re := regexp.MustCompile(`^https?://`)
	return re.ReplaceAllString(str, "")
}

// 转换点击数量 字符串 比如：5.8w 5.7千 5.6亿
/*
可传参：
	1. 带单位   比如：5.8w 5.7千 5.6亿
	2. 不带单位 比如：5600
*/
func ParseHitsStr(hitsStr string) int {
	// 1. 预处理：去空格，繁体转简体
	hitsStr = strings.TrimSpace(hitsStr)
	hitsStr, _ = langutil.TraditionalToSimplified(hitsStr)

	// 2. 正则匹配：同时匹配数字和单位（支持中英文单位）
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([万千亿kw]?)`)
	matches := re.FindStringSubmatch(hitsStr)

	if len(matches) < 2 {
		return 0 // 没有找到数字，返回0
	}

	// 3. 解析数字部分
	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil || num < 0 {
		return 0
	}

	// 4. 根据单位进行乘法转换
	multiplier := 1
	if len(matches) >= 3 && matches[2] != "" {
		switch matches[2] {
		case "亿":
			multiplier = 100000000
		case "万", "w":
			multiplier = 10000
		case "千", "k":
			multiplier = 1000
		}
	}

	return int(num * float64(multiplier))
}
