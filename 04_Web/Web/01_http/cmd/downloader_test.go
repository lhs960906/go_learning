package cmd_test

import (
	"log"
	"testing"

	"github.com/lhs960906/Go-Learning/04_Web/Web/01_http/cmd"
)

func TestExecute(t *testing.T) {
	err := cmd.Downloader.Execute()
	if err != nil {
		log.Println(err)
	}
}
