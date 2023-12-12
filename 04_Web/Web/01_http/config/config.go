package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	GlobalConfig *GlobalConfig
	// StorageConfig *StorageConfig
}

// type StorageConfig struct {
// 	// MySQLConfig  *MySQLConfig
// }

type GlobalConfig struct {
	Concurrency int
	Filepath    string
}

func init() {
	// 配置设置和读取
	viper.SetDefault("concurrency", 3)
	viper.SetDefault("config", "../etc/config.yaml")
	viper.SetConfigFile("../etc/config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")             // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")               // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("../etc/")            // 查找配置文件所在的路径
	viper.AddConfigPath(".")                  // 多次调用以添加多个搜索路径；还可以在工作目录中查找配置
	err := viper.ReadInConfig()               // 查找并读取配置文件
	if err != nil {                           // 处理读取配置文件的错误
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
