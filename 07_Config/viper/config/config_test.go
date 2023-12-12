package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lhs960906/Go-Learning/07_Config/viper/config"
)

func TestLoad(t *testing.T) {
	// 加载配置
	conf := config.NewConfig().Load()

	// 将对象，转换成json格式
	data, err := json.Marshal(conf)

	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	// 最终以json格式输出
	fmt.Println("data:\t", string(data))
}
