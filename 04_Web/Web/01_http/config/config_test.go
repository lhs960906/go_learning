package config_test

import (
	"fmt"
	"testing"

	"github.com/lhs960906/Go-Learning/04_Web/Web/01_http/config"
)

func TestConfig(t *testing.T) {
	var a *config.GlobalConfig = &config.GlobalConfig{
		Concurrency: 3,
		Filepath:    "test",
	}
	fmt.Println(a)
}
