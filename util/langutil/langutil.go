// 语言工具包
package langutil

import (
	"log"
	"strings"
	"sync"

	"github.com/longbridgeapp/opencc"
)

var (
	// converter *opencc.Converter
	converter *opencc.OpenCC
	once      sync.Once
)

// TraditionalToSimplified 将繁体中文转换为简体中文
func TraditionalToSimplified(text string) (string, error) {
	once.Do(func() {
		var err error
		converter, err = opencc.New("t2s")
		if err != nil {
			log.Fatal(err)
		}
	})
	result, err := converter.Convert(text)
	if err != nil {
		return "", err
	}
	return result, nil
}

// IsHTTPOrHTTPS 判断Url是否以 http:// 或 https:// 开头
// 返回值: bool
func IsHTTPOrHTTPS(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// IsHTTPS 判断Url是否以  https:// 开头
// 返回值: bool
func IsHTTPS(url string) bool {
	return strings.HasPrefix(url, "https://")
}
