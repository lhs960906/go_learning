package src

import (
	"errors"
	"fmt"
	"net/http"
)

func Head(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		// 支持部分请求
		fmt.Printf("Status code: %d\n", resp.StatusCode)
		// 获取 http 响应中的 Accept-Ranges 头部值
		fmt.Printf("Accept-Ranges: %s\n", resp.Header.Get("Accept-Ranges"))
		return true, nil
	}

	return false, errors.New("unkown error")
}

// func main() {
// 	url := "https://studygolang.com/dl/golang/go1.16.5.src.tar.gz"
// 	Head(url)
// }
