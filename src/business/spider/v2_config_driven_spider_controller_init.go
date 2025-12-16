/*
配置文件 驱动 爬虫程序的 控制器，初始化功能
*/
// 配置驱动爬虫的初始化 - 不修改main.go的情况下提供初始化功能
// 用法：在main.go中添加一行代码调用InitConfigDrivenSpider()

package spider

import (
	"study-spider-manhua-gin/src/log"

	"github.com/gin-gonic/gin"
)

// InitConfigDrivenSpiderControllerV2 初始化配置驱动爬虫系统/控制器 V2版本
// 在main.go中调用此函数即可启用配置驱动功能
func InitConfigDrivenSpiderControllerV2(router *gin.Engine, configPath string) error {
	log.Info("开始初始化配置驱动爬虫系统...")

	// 创建API 控制器 实例
	apiC := NewConfigDrivenAPIController()

	// 初始化配置
	if err := apiC.InitConfig(configPath); err != nil {
		log.Errorf("配置驱动爬虫初始化失败: %v", err)
		return err
	}

	// 注册路由
	apiC.RegisterRoutes(router)

	log.Info("配置驱动爬虫系统V2初始化完成")
	log.Info("可用API端点:")
	log.Info("  POST /api/v2/spider/oneTypeAllBookByHtml/config - 配置驱动爬取")
	log.Info("  GET  /api/v2/spider/websites - 获取支持的网站")
	log.Info("  GET  /api/v2/spider/config - 获取网站配置")
	log.Info("  POST /api/v2/spider/validate - 验证配置")

	return nil
}

// 示例：在main.go中添加以下代码即可启用：
//
//	import "study-spider-manhua-gin/src/business/spider"
//
//	func main() {
//		// ... 现有代码 ...
//
//		// 初始化配置驱动爬虫V2（添加这一行）
//		err := spider.InitConfigDrivenSpiderV2(router, "spider-config-v2.yaml")
//		if err != nil {
//			log.Fatal("初始化配置驱动爬虫失败:", err)
//		}
//
//		// ... 现有代码 ...
//	}
//
// 这样既不修改任何现有代码，又能启用全新的配置驱动功能！
