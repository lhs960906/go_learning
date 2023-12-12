package config_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lhs960906/Go-Learning/07_Config/viper/config"
	"github.com/spf13/viper"
)

var (
	vip *viper.Viper
)

func TestLoad(t *testing.T) {
	err := vip.Unmarshal(config.Cfg)
	if err != nil {
		panic(err)
	}

	// 将对象，转换成json格式
	data, err := json.Marshal(config.Cfg)

	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	// 最终以json格式输出
	fmt.Println("data:\t", string(data))

	for {
		time.Sleep(1 * time.Second)
	}
}
func init() {
	log.Println("Loading configuration logics...")
	// 加载配置
	vip = config.Load()
	go config.DynamicReloadConfig()
}
