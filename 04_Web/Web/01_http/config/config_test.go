package config_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lhs960906/Go-Learning/04_Web/Web/01_http/config"
	"github.com/spf13/viper"
)

var (
	vip *viper.Viper
)

func TestConfig(t *testing.T) {
	// 读取配置到 Config 实例中
	if err := vip.Unmarshal(config.Cfg); err != nil {
		panic(err)
	}
	// 将对象，转换成json格式
	data, err := json.Marshal(config.Cfg)

	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	// 最终以json格式输出
	log.Println("data: ", string(data))

	for {
		time.Sleep(1 * time.Second)
	}
}

func init() {
	log.Println("Loading configuration logics...")
	// 加载配置
	vip = config.Load()
	// 动态加载配置
	go config.DynamicReloadConfig()
}
