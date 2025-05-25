// 对象操作工具
package objutil

// 拷贝一个对象，用于单元测试。因为默认 增删改查操作，会去除原始对象的空格
// 参数1：原始obj 返回新 拷贝的obj
func Copy(originObj interface{}) interface{} {
	newObj := originObj
	return newObj
}
