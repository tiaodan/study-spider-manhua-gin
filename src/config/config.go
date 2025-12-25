package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// 配置文件 结构体
type Config struct {
	// 日志相关
	Log struct {
		Level string `mapstructure:"level"`
		Path  string `mapstructure:"path"`
	}

	// 网络相关
	Network struct {
		XimalayaIIp string `mapstructure:"ximalaya_ip"`
	}

	// 数据库相关
	DB struct {
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	// gin api接口框架相关
	Gin struct {
		Mode string `mapstructure:"mode"`
	}

	// 爬虫相关
	Spider struct {
		// -- 公用配置
		Public struct {
			// 爬取某一类相关 --
			SpiderType struct {
				RandomDelayTime      int `mapstructure:"random_delay_time"`       // 爬虫对象 colly, 每次请求前，随机延迟时间。单位 = 秒
				QueueLimitConcMaxnum int `mapstructure:"queue_limit_conc_maxnum"` // 爬虫队列, 爬虫限制最大并发数
				QueuePoolMaxnum      int `mapstructure:"queue_pool_maxnum"`       // 爬虫队列, 队列池最大数
			}
		}
	}
}

var (
	Cfg  *Config   // 全局变量,让其他包可以访问. 对应 根目录 cofig.yaml这个文件
	once sync.Once //保证单例初始化
)

// GetConfig 获取配置实例（单例）
/*
思路:
   1. 初始化Viper
   2. 设置默认值, 防止用户没配置, 读取到空值
   3. 读取配置文件
   4. 将配置文件解析到结构体
   5. 返回配置指针

参数:
	1. path string 配置文件搜索路径（当前目录）
	2. name string 配置文件名（不含扩展名）
	3. ext string 配置文件扩展名（.ini、.yaml 等）

使用方式：
	如main.go调用
	// 获取配置实例（首次调用时触发初始化）
	cfg := config.GetConfig(".", "config", "yaml")

	// 使用配置
	log.Println("network.ximalayaIIp_ip: ", cfg.Network.XimalayaIIp)

	// 读取配置文件，并设置为日志级别, 默认info
	switch cfg.Log.Level {
	case "debug":
		logger.SetLogLevel(logger.LevelDebug)
	case "info":
		logger.SetLogLevel(logger.LevelInfo)
	case "warn":
		logger.SetLogLevel(logger.LevelWarn)
	case "error":
		logger.SetLogLevel(logger.LevelError)
	default:
		logger.SetLogLevel(logger.LevelInfo)
	}
*/
func GetConfig(path, name, ext string) *Config {
	once.Do(func() {
		// 1. 初始化Viper
		viper.AddConfigPath(path) // 配置文件搜索路径（当前目录），如 “.”
		viper.SetConfigName(name) // 配置文件名（不含扩展名）, 如 "config"
		viper.SetConfigType(ext)  // 文件类型（yaml、json 等）, 如 “ini”

		// 2. 设置默认值, 防止用户没配置, 读取到空值
		// 设置默认值 [log] 相关
		viper.SetDefault("log.level", "info")       // 设置默认info级别
		viper.SetDefault("log.path", "log/app.log") // 默认日志文件名

		// 设置默认值 [network] 相关
		viper.SetDefault("network.ximalaya_ip", "www.ximalaya.com")

		// 设置默认值 [db] 相关
		viper.SetDefault("db.name", "comic")
		viper.SetDefault("db.user", "root")
		viper.SetDefault("db.password", "123456")

		// 设置默认值 [gin] 相关
		viper.SetDefault("gin.mode", "release") // 设置默认release模式

		// 设置默认值 [spider] 相关
		// 公共配置,只要涉及到爬虫的 --
		viper.SetDefault("spider.public.spider_type.random_delay_time", 5)       // 爬虫对象 colly, 每次请求前，随机延迟时间。单位 = 秒
		viper.SetDefault("spider.public.spider_type.queue_limit_conc_maxnum", 5) // 爬虫队列, 爬虫限制最大并发数
		viper.SetDefault("spider.public.spider_type.queue_pool_maxnum", 100)     // 爬虫队列, 队列池最大数

		// 3. 读取配置文件
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("读取配置文件失败,err: ", err)
		}

		// 4. 将配置文件解析到结构体
		Cfg = &Config{}
		if err := viper.Unmarshal(Cfg); err != nil {
			log.Fatalln("解析配置文件失败,err: ", err)
		}
	})

	// 5. 返回配置指针
	return Cfg
}

// WriteConfig4Blank 将配置写入文件, 不带注释, 4个空格缩进
func WriteConfig4Blank(cfg *Config) error {
	// 将结构体转换为 map[string]interface{}
	var newCfg map[string]interface{}
	if err := mapstructure.Decode(cfg, &newCfg); err != nil {
		return fmt.Errorf("结构体转 Map 失败:: %v", err)
	}

	// 将 Map 合并到 Viper
	viper.MergeConfigMap(newCfg)

	// 将Viper配置写入文件
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	return nil
}

// WriteConfig2Blank 将配置写入文件, 不带注释, 2个空格缩进
/*
参数:
	1. path string 配置文件搜索路径（当前目录）
	2. name string 配置文件名（不含扩展名）
	3. ext string 配置文件扩展名（.ini、.yaml 等）

调用方式：
if err := config.WriteConfig2Blank(); err != nil {
	log.Fatalln("写入配置文件失败,err: ", err)
}
*/
func WriteConfig2Blank(path, name, ext string) error {
	// 获取 Viper 的所有配置（map 格式）
	cfgMap := viper.AllSettings()

	// 创建 YAML 编码器并设置缩进
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2) // 关键：设置缩进为 2 个空格

	// 将配置写入文件（或输出流）
	filepathStr := path + "/" + name + "." + ext
	log.Println("filepathStr= ", filepathStr)
	file, err := os.Create(filepathStr) // "config.yaml"
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建文件编码器
	fileEncoder := yaml.NewEncoder(file)
	fileEncoder.SetIndent(2)

	// 编码并写入文件
	if err := fileEncoder.Encode(cfgMap); err != nil {
		panic(fmt.Sprintf("YAML 编码失败: %v", err))
	}

	fmt.Println("配置文件已生成（缩进 2 空格）")
	return err
}
