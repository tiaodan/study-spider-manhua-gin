/*
*
说明：
  - main.go 是项目的入口文件。必须放在根目录，不能放在src目录下。
  - 因为：mian.go读取配置，逻辑：读取 ../config.yaml。
    如果在src打包成main.exe后，放到根目录，没有上级目录，会报错。因为打包成实际目录后，src源代码不会保留
*/
package main

import (
	"study-spider-manhua-gin/src/business/spider"
	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 初始化, 默认main会自动调用本方法
/*
作用简单说：
	- 初始化项目

作用:
	- 初始化配置
		- 读取配置
		- 不重新生成配置文件
	- 初始化日志
		- 创建日志文件
		- 根据配置，设置日志级别
		- 让日志可同时写到终端和日志文件
	- 初始化数据库
		- 连接数据库
		- 创建数据表结构
		- 插入默认数据

思路:
	1. 初始化配置
	2. 初始化日志
	3. 初始化数据库
*/
func init() {
	// 1. 初始化配置
	// -- 读取配置文件， (如果配置文件不填, 自动会有默认值)
	cfg := config.GetConfig(".", "config", "yaml")

	// 2. 初始化日志 (现在用logrus框架)
	// -- 创建日志文件
	log.InitLog(cfg.Log.Path)

	// -- 根据配置，设置日志级别
	// 获取日志对象
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

	// 打印配置。用debug打日志是为了安全
	log.Debug("打印配置: --------------------------------- ")
	log.Debug("[log] 相关")
	log.Debug("log.level = ", cfg.Log.Level)
	log.Debug("log.path = ", cfg.Log.Path)
	log.Debug("[network] 相关---")
	log.Debug("network.ximalayaIIp_ip = ", cfg.Network.XimalayaIIp)
	log.Debug("[db] 相关 ")
	log.Debug("db.name = ", cfg.DB.Name)
	log.Debug("db.user = ", cfg.DB.User)
	log.Debug("db.password = ", cfg.DB.Password)
	log.Debug("[gin] 相关")
	log.Debug("gin.mode = ", cfg.Gin.Mode)
	log.Debug("[spider] 相关")
	log.Debug("[spider] 相关 - 公用配置: ")
	log.Debug("[spider] 相关 - 公用配置 - 爬取某一类 配置: ")
	log.Debug("每次请求前随机延迟 (秒) random_delay_time = ", cfg.Spider.Public.SpiderType.RandomDelayTime)
	log.Debug("爬虫队列, 最大并发数 queue_limit_conc_maxnum = ", cfg.Spider.Public.SpiderType.QueueLimitConcMaxnum)
	log.Debug("爬虫队列, 池最大数 queue_pool_maxnum = ", cfg.Spider.Public.SpiderType.QueuePoolMaxnum)

	log.Debug("打印配置: --------------------------------- end  ")

	// 3. 初始化数据库
	// -- 初始化数据库连接
	db.InitDB("mysql", cfg.DB.Name, cfg.DB.User, cfg.DB.Password)

	// -- 自动迁移表结构
	/*
		创建时有讲究的，一般先创新主表，再创建从表。因为从表要关联主表id，主表id没有会报错。
		比如 &models.ComicSpiderStats{}, &models.ComicSpider{},一定要 ComicSpider主表在前
	*/
	err := db.DBComic.AutoMigrate(&models.Website{}, &models.Country{}, &models.PornType{}, &models.Type{},
		&models.ComicSpider{}, &models.ComicSpiderStats{}, &models.ComicMy{}, &models.ComicMyStats{},
		&models.WebsiteType{}, &models.Process{}, &models.Author{},
		&models.ChapterSpider{}, &models.ChapterMy{},
		&models.ChapterContentSpider{}, &models.ChapterContentMy{}) // 有几个表, 写几个参数
	errorutil.ErrorPanic(err, "自动迁移表结构报错, err = ")

	// -- 插入默认数据
	db.InsertDefaultData()
}

// main 入口函数
/*
作用简单说：
  - 初始化项目
  - 提供API接口(增删改查相关),让前端调用
  - 提供API接口(爬虫相关),让前端调用

作用详细说:

	通过调用init()函数实现
	- 初始化配置
		- 读取配置
		- 不重新生成配置文件
	- 初始化日志
		- 根据配置，设置日志级别
		- 创建日志文件
		- 让日志可同时写到终端和日志文件
	- 初始化数据库
		- 连接数据库
		- 创建数据表结构
		- 插入默认数据

思路:
   1. 通过调用init()函数实现 -》 初始化项目
   2. 提供API接口(增删改查相关)
   3. 提供API接口(爬虫相关)
*/
func main() {

	// 1. 通过调用init()函数实现 -》 初始化项目

	// 2. 封装restful api
	// 关键代码：切换到 release 模式，防止打过多日志 --
	gin.SetMode(gin.ReleaseMode)

	// 创建 gin 实例 --
	r := gin.Default()
	r.Use(cors.Default()) // 允许所有跨域

	// 提供接口 --
	// 这些接口是参考，不用了 --
	/*
		r.POST("/orders", order.OrderAdd)
		r.DELETE("/orders/:id", order.OrderDelete)
		r.PUT("/orders", order.OrderUpdate)
		r.GET("/orders", order.OrdersPageQuery) // 分页查询

		r.POST("/comics", comic.ComicAdd)
		r.DELETE("/comics/:id", comic.ComicDelete)
		r.PUT("/comics", comic.ComicUpdateByIdOmitIndex)
		r.GET("/comics", comic.ComicsQueryByPage) // 分页查询
	*/

	// 3. 提供API接口(爬虫相关)
	// 爬虫
	// 爬虫思路：
	// 1. 爬某一类漫画所有内容
	// 2. 爬某个漫画的所有章节，更新该漫画具体内容

	// 流程：爬完漫画（spider_end）-》爬章节-》修改漫画-》 存章节-》下载漫画(download_end)-》下载章节-》下载完，上传aws章节(upload_aws_end)-》传完，更新漫画标志位
	// r.POST("/spider/oneCategory", spider.Spider) // v0.1 写法，没用通用爬虫模板，弃用
	r.POST("/spider/oneTypeByHtml", spider.BookTemSpiderTypeByHtml)       // v0.2 写法，用通用爬虫模板,推荐。爬html页面
	r.POST("/spider/oneTypeByJson", spider.DispatchApi_OneCategoryByJSON) // v0.2 写法，用通用爬虫模板,推荐。爬F12 目标网站返回的json数据

	// -- html spider 相关 V1 策略实现
	r.POST("/spider/oneTypeAllBookByHtml", spider.DispatchApi_SpiderOneTypeAllBookArr_Template) // 通用模板
	r.POST("/spider/oneBookAllChapterByHtml", spider.DispatchApi_OneBookAllChapterByHtml)       // v0.2 写法，用通用爬虫模板,推荐。爬html页面 - 爬一本书所有章节
	r.POST("/spider/oneChapterAllContentByHtml", spider.DispatchApi_OneChapterAllContentByHtml) // v0.2 写法，用通用爬虫模板,推荐。爬html页面 - 爬章节所有内容 - 没实现一章节所有内容

	// v2 写法 获取配置 - start
	// -- spider 相关 V2 策略实现
	err := spider.InitConfigDrivenSpiderControllerV2(r, "v2-spider-config.yaml")
	if err != nil {
		log.Error("初始化配置驱动爬虫失败, err = ", err)
	}
	// v2 写法 获取配置 - end

	r.Run(":8888") // 启动服务
}
