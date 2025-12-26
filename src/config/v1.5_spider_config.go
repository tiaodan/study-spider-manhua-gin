/*
功能: 获取 v2-spider-config.yaml 配置文件所有内容. v1.5版本逻辑
 1. 采用2种获取配置文件方式
    - yaml.v3 库 + yaml tag, 实现起来较简单，不能处理复杂情况: 比如命令行替换 某个配置. 核心代码：yaml.Unmarshal()
    - viper库 + mapstructure tag, 实现起来较复杂.能处理复杂情况: 比如命令行替换 某个配置. 核心代码：viper.Unmarshal()
    - 一般推荐2种写法都用，用户想用哪种就用哪种。聪明人不做选择，我都要!!!
*/
package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// ------------------------------------------- 初始化 -------------------------------------------
var (
	CfgSpiderYaml     *SpiderConfigV15 // 全局变量,让其他包可以访问. 对应 根目录 cofig.yaml这个文件
	CfgSpiderYamlOnce sync.Once        //保证单例初始化
)

// ------------------------------------------- 增删改查、读取、写入, 方法 -------------------------------------------
// 加载配置文件, 给全局变量CfgSpiderYaml赋值。 v2-spider-config.yaml. 用yaml.v3库 + yaml tag
/*
思路：
1. 创建单例
2. 读取配置文件 os.ReadFile
3. 解析配置文件 yaml.Unmarshal
4. 返回结果
*/
func LoadSpiderConfigFromYAMLUseTagYaml(configPath string) error {
	// 1. 创建单例
	CfgSpiderYamlOnce.Do(func() {
		// 2. 读取配置文件 os.ReadFile
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("加载配置文件 %v 失败, err= %v", configPath, err) // 会中断程序
			return
		}
		// 3. 解析配置文件 yaml.Unmarshal
		if err := yaml.Unmarshal(data, &CfgSpiderYaml); err != nil {
			log.Fatalf("解析配置文件 %v 失败, err= %v", configPath, err) // 会中断程序
			return
		}
	})

	fmt.Println("------- delete  CfgSpiderYaml = ", CfgSpiderYaml)
	// 4. 返回结果
	return nil // 说明成功
}

// 读取配置文件, 给全局变量CfgSpiderYaml赋值。 v2-spider-config.yaml. 用viper库 + mapstructure tag
/*
思路：
1. 创建单例
2. 初始化 vipder
3. 读取配置文件 viper.ReadInConfig, 并处理错误:文件不存在、格式错误等
4. 解析配置文件 viper.Unmarshal
5. 返回结果
*/
func LoadSpiderConfigFromYAMLUseViperAndTagMapstructure(configPath string) error {
	// 1. 创建单例
	CfgSpiderYamlOnce.Do(func() {
		// 2. 初始化 vipder. 有2种写法
		if configPath != "" {
			viper.SetConfigFile(configPath) // 写带类型的路径就行，如 "config.yaml"
		} else {
			// 这里就写死内容，只是展示下写法
			viper.AddConfigPath(".")      // 配置文件搜索路径（当前目录），如 “.”   写法：viper.AddConfigPath(path)
			viper.SetConfigName("config") // 配置文件名（不含扩展名）, 如 "config"  写法：viper.SetConfigName(name)
			viper.SetConfigType("yaml")   // 文件类型（yaml、json 等）, 如 “ini”   写法：viper.SetConfigType(ext)
		}

		// 3. 读取配置文件 viper.ReadInConfig, 并处理错误:文件不存在、格式错误等
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("读取配置文件 %v 失败, err= %v", configPath, err) // 会中断程序
			return
		}

		// 4. 解析配置文件 viper.Unmarshal
		if err := viper.Unmarshal(&CfgSpiderYaml); err != nil {
			log.Fatalf("解析配置文件 %v 失败, err= %v", configPath, err) // 会中断程序
			return
		}
	})

	// 5. 返回结果
	return nil // 一切成功
}
