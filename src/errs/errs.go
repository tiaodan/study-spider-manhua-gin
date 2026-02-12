// 统一定义错误变量 错误
/*
知识：
1. 不用写成单例。因为创建时机是：项目启动时 / 包被第一次导入时
2. 在程序启动时（main函数执行前），Go 语言会自动初始化这个包的所有包级变量
*/

package errs

import "errors"

var (
	// 参数错误
	ErrParams         = errors.New("params error")                        // 参数错误
	ErrInvalidPageNum = errors.New("参数错误, endPageNum 必须 >= startPageNum") // endPageNum 必须 >= startPageNum
	ErrNull           = errors.New("数据为空")
	ErrNoGetConfig    = errors.New("没有获取到配置") // 没有获取到配置
	ErrInvalidArgs    = errors.New("参数错误")

	// 后面可以写一些 string 类型报错内容
)

// 获取错误
func GetErr(errStr string) error {
	return errors.New(errStr)
}
