/*
功能: config.yaml 配置文件 仅结构体, 不提供get/set方法

注意:
v1实现方式：
  - 结构体：只提供了 mapstructure tag方式
  - 并实现了 mapstructure + viper方式获取配置

V1.5实现方式：
1. 提供 2种获取配置文件 tag方式
  - yaml.v3 库 + yaml tag, 实现起来较简单，不能处理复杂情况: 比如命令行替换 某个配置. 核心代码：yaml.Unmarshal()
  - viper库 + mapstructure tag, 实现起来较复杂.能处理复杂情况: 比如命令行替换 某个配置. 核心代码：viper.Unmarshal()
  - 一般推荐2种写法都用，用户想用哪种就用哪种。聪明人不做选择，我都要!!!
*/
package config

//v1.5 方式 结构体定义, 用2种tag： yaml + mapstructure
// 配置文件 结构体
type Config struct {
	// 日志相关
	Log struct {
		Level string `mapstructure:"level" yaml:"level"`
		Path  string `mapstructure:"path" yaml:"path"`
	}

	// 网络相关
	Network struct {
		XimalayaIIp string `mapstructure:"ximalaya_ip" yaml:"ximalaya_ip"`
	}

	// 数据库相关
	DB struct {
		Name     string `mapstructure:"name" yaml:"name"`
		User     string `mapstructure:"user" yaml:"user"`
		Password string `mapstructure:"password" yaml:"password"`
	}

	// gin api接口框架相关
	Gin struct {
		Mode string `mapstructure:"mode" yaml:"mode"`
	}

	// 爬虫相关
	Spider struct {
		// -- 公用配置
		Public struct {
			// 爬取某一类相关 --
			SpiderType struct {
				RandomDelayTime      int `mapstructure:"random_delay_time" yaml:"random_delay_time"`             // 爬虫对象 colly, 每次请求前，随机延迟时间。单位 = 秒
				QueueLimitConcMaxnum int `mapstructure:"queue_limit_conc_maxnum" yaml:"queue_limit_conc_maxnum"` // 爬虫队列, 爬虫限制最大并发数
				QueuePoolMaxnum      int `mapstructure:"queue_pool_maxnum" yaml:"queue_pool_maxnum"`             // 爬虫队列, 队列池最大数
			}
		}
	}
}

/* v1 方式 结构体定义, 仅用于参考, 没有使用。后续考虑用2种tag： yaml + mapstructure
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
*/
