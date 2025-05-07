package main

import (
	"io"
	"os"
	"study-spider-manhua-gin/business/comic"
	"study-spider-manhua-gin/business/order"
	"study-spider-manhua-gin/business/spider"
	"study-spider-manhua-gin/config"
	"study-spider-manhua-gin/db"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 初始化, 默认main会自动调用本方法
func init() {
	// 1. 读取配置文件， (如果配置文件不填, 自动会有默认值)
	cfg := config.GetConfig(".", "config", "yaml")

	// 2. 根据配置文件,设置日志相关,现在用logrus框架
	log.InitLog()

	// 获取日志实例
	log := log.GetLogger()

	// 设置日志级别
	switch cfg.Log.Level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// 创建一个文件用于写入日志
	file, err := os.OpenFile(cfg.Log.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // os.OpenFile("app.log"
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
	}

	// 使用 io.MultiWriter 实现多写入器功能
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	// 打印配置
	log.Debug("[log] 相关")
	log.Debug("log.level: ", cfg.Log.Level)
	log.Debug("log.path: ", cfg.Log.Path)
	log.Debug("[network] 相关---")
	log.Debug("network.ximalayaIIp_ip: ", cfg.Network.XimalayaIIp)
	log.Debug("[db] 相关")
	log.Debug("db.name: ", cfg.DB.Name)
	log.Debug("db.user: ", cfg.DB.User)
	log.Debug("db.password: ", cfg.DB.Password)
	log.Debug("[gin] 相关")
	log.Debug("gin.mode: ", cfg.Gin.Mode)

	// 初始化数据库连接
	db.InitDB("mysql", cfg.DB.Name, cfg.DB.User, cfg.DB.Password)

	// 自动迁移表结构
	db.DB.AutoMigrate(&models.Website{}, &models.Country{}, &models.Category{}, &models.Type{}, &models.Comic{}) // 有几个表, 写几个参数

	// 插入默认数据
	db.InsertDefaultData()
}

/*
思路:
 1. 读取配置文件， (如果配置文件不填, 自动会有默认值)
 2. 设置日志级别, 默认info
 3. 统一调用错误打印, 封装函数
 4. 封装restful api
*/
func main() {

	// 1. 读取配置文件， (如果配置文件不填, 自动会有默认值)
	// 2. 设置日志级别, 默认info
	// 3. 统一调用错误打印, 封装函数
	// 4. 封装restful api

	// 等会再用 ----------------------------- start
	gin.SetMode(gin.ReleaseMode) // 关键代码：切换到 release 模式
	r := gin.Default()
	r.Use(cors.Default()) // 允许所有跨域

	// 封装api
	//---------------------------- 一会再弄这个
	r.POST("/orders", order.OrderAdd)
	r.DELETE("/orders/:id", order.OrderDelete)
	r.PUT("/orders", order.OrderUpdate)
	r.GET("/orders", order.OrdersPageQuery) // 分页查询

	r.POST("/comics", comic.ComicAdd)
	r.DELETE("/comics/:id", comic.ComicDelete)
	r.PUT("/comics", comic.ComicUpdateByIdOmitIndex)
	r.GET("/comics", comic.ComicsPageQuery) // 分页查询

	// 爬虫
	// 爬虫思路：
	// 1. 爬某一类漫画所有内容
	// 2. 爬某个漫画的所有章节，更新该漫画具体内容

	// 流程：爬完漫画（spider_end）-》爬章节-》修改漫画-》 存章节-》下载漫画(download_end)-》下载章节-》下载完，上传aws章节(upload_aws_end)-》传完，更新漫画标志位
	r.POST("/spider/oneCategory", spider.Spider)

	r.Run(":8888") // 启动服务
	// 等会再用 ----------------------------- end
}
