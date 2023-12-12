package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lhs960906/Go-Learning/07_Config/viper/config"
)

func TestGetConf(t *testing.T) {
	var c *config.KafkaCluster
	//读取yaml配置文件, 将yaml配置文件，转换struct类型
	conf := c.GetConf()

	//将对象，转换成json格式
	data, err := json.Marshal(conf)

	if err != nil {
		fmt.Println("err:\t", err.Error())
		return
	}

	//最终以json格式，输出
	fmt.Println("data:\t", string(data))
}
