package config

import (
	"os"

	"github.com/spf13/viper"
)

type KafkaCluster struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
	//map类型
	Labels map[string]*NodeServer `yaml:"labels"`
}

type NodeServer struct {
	Address string `yaml:"address"`
	Id      string `yaml:"id"`
	Name    string `yaml:"name"`
	//注意，属性里，如果有大写的话，tag里不能存在空格
	//如yaml: "nodeName" 格式是错误的，中间多了一个空格，不能识别的
	NodeName string `yaml:"nodeName"`
	Role     string `yaml:"role"`
}

type Spec struct {
	Replicas   int          `yaml:"replicas"`
	Name       string       `yaml:"name"`
	Image      string       `yaml:"iamge"`
	Ports      int          `yaml:"ports"`
	Conditions []Conditions `yaml:"conditions"`
}

type Conditions struct {
	ContainerPort int      `yaml:"containerPort"`
	Requests      Requests `yaml:"requests"`
	Limits        Limits   `yaml:"limits"`
}

type Requests struct {
	CPU    string `yaml:"cpu"`
	MEMORY string `yaml:"memory"`
}

type Limits struct {
	CPU    string `yaml:"cpu"`
	MEMORY string `yaml:"memory"`
}

// 读取Yaml配置文件，并转换成KafkaCluster对象  struct结构
func (kafkaCluster *KafkaCluster) GetConf() *KafkaCluster {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip := viper.New()
	vip.AddConfigPath(path + "/" + "../etc/") //设置读取的文件路径
	vip.AddConfigPath(".")                    //设置读取的文件路径,当前路径
	vip.SetConfigName("config")               //设置读取的文件名
	vip.SetConfigType("yaml")                 //设置文件的类型

	//尝试进行配置读取
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}

	err = vip.Unmarshal(&kafkaCluster)
	if err != nil {
		panic(err)
	}

	return kafkaCluster
}
