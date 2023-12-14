package config

import (
	"encoding/json"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	vip *viper.Viper
	Cfg = &Config{}
)

type Config struct {
	Global *Global `json:"global"`
	// StorageConfig *StorageConfig
}

// type StorageConfig struct {
// 	// MySQLConfig  *MySQLConfig
// }

type Global struct {
	Concurrency int `json:"concurrency"`
}

// 返回配置的并发数
func (conf *Config) GetConcurrency() int {
	return Cfg.Global.Concurrency
}

// 将配置加载 viper 实例中
func Load() *viper.Viper {
	vip = viper.New()
	vip.AddConfigPath("etc/")    // 设置读取的文件路径
	vip.AddConfigPath("../etc/") // 设置读取的文件路径
	vip.AddConfigPath(".")       // 设置读取的文件路径,当前路径
	vip.SetConfigName("config")  // 设置读取的文件名
	vip.SetConfigType("yaml")    // 设置文件的类型

	// 尝试进行配置读取
	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return vip
}

// 动态加载配置到 Config 实例中
func DynamicReloadConfig() {
	vip.WatchConfig()
	vip.OnConfigChange(func(event fsnotify.Event) {
		log.Printf("Detect config change: %s \n", event.String())

		// 热加载配置到 Cfg
		vip.Unmarshal(Cfg)

		// 配置重新序列化为 json 并以 json 格式输出
		data, err := json.Marshal(Cfg)
		if err != nil {
			log.Printf("err:, %v\t", err.Error())
			return
		}
		log.Printf("data: %s\t", string(data))
	})
}
