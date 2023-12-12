package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Mysql Mysql
}

func NewConfig() *Config {
	return &Config{}
}

type Mysql struct {
	Active string `mapstructure:"active"` // 测试配置 or 生产配置
	Env    Env    `mapstructure:"env"`
}

type Env struct {
	Pro  Pro  `mapstructure:"pro"`
	Test Test `mapstructure:"test"`
}

type Pro struct {
	Ip     string `mapstructure:"ip"`
	Port   int    `mapstructure:"port"`
	DbName string `mapstructure:"dbName"`
	DbUser string `mapstructure:"dbUser"`
	DbPass string `mapstructure:"dbPass"`
}

type Test struct {
	Ip     string `mapstructure:"ip"`
	Port   int    `mapstructure:"port"`
	DbName string `mapstructure:"dbName"`
	DbUser string `mapstructure:"dbUser"`
	DbPass string `mapstructure:"dbPass"`
}

// 读取Yaml配置文件，并转换成Config对象
func (conf *Config) Load() *Config {
	// 获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip := viper.New()
	vip.AddConfigPath(path + "/../etc/") // 设置读取的文件路径
	vip.AddConfigPath(".")               // 设置读取的文件路径,当前路径
	vip.SetConfigName("config")          // 设置读取的文件名
	vip.SetConfigType("yaml")            // 设置文件的类型

	// 尝试进行配置读取
	err = vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 将配置反序列化到 conf 中
	err = vip.Unmarshal(conf)
	if err != nil {
		panic(err)
	}

	return conf
}
