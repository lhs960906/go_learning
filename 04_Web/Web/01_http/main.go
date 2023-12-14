package main

import (
	"fmt"

	"github.com/lhs960906/Go-Learning/04_Web/Web/01_http/config"
)

func main() {
	// 加载配置并获取配置的并发数
	vip := config.Load()
	vip.Unmarshal(config.Cfg)
	concurrency := config.Cfg.GetConcurrency()

	fmt.Println(concurrency)
}
