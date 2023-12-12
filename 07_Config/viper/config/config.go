package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Cfg *Config = &Config{}
	vip *viper.Viper
)

type Config struct {
	Mysql Mysql `json:"mysql"`
}

type Mysql struct {
	Active string `json:"active"` // 测试配置 or 生产配置
	Env    Env    `json:"env"`
}

type Env struct {
	Pro  Pro  `json:"pro"`
	Test Test `json:"test"`
}

type Pro struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	DbName string `json:"dbName"`
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
}

type Test struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	DbName string `json:"dbName"`
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
}

// 读取Yaml配置文件，并转换成Config对象
func Load() *viper.Viper {
	// 获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip = viper.New()
	vip.AddConfigPath(path + "/../etc/") // 设置读取的文件路径
	vip.AddConfigPath(".")               // 设置读取的文件路径,当前路径
	vip.SetConfigName("config")          // 设置读取的文件名
	vip.SetConfigType("yaml")            // 设置文件的类型

	// 尝试进行配置读取
	err = vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return vip

	// // 将配置反序列化到 conf 中
	// err = vip.Unmarshal(conf)
	// if err != nil {
	// 	panic(err)
	// }

	// return conf
}

func DynamicReloadConfig() {
	vip.WatchConfig()
	vip.OnConfigChange(func(event fsnotify.Event) {
		log.Printf("Detect config change: %s \n", event.String())

		// 热加载配置
		vip.Unmarshal(Cfg)

		// 重新序列化为 json 并以 json 格式输出
		data, err := json.Marshal(Cfg)
		if err != nil {
			log.Printf("err:, %v\t", err.Error())
			return
		}
		log.Printf("data: %s\t", string(data))
	})
}
