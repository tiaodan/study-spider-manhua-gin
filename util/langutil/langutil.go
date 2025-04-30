// 语言工具包
package langutil

import (
	"log"
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
