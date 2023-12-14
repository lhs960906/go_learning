package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 发送一个 head 请求
	req, err := http.NewRequest(http.MethodHead, "https://download.jetbrains.com/go/goland-2020.2.2.exe", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	header1 := resp.Header.Get("Accept-Ranges")
	header2 := resp.Header.Get("Content-Length")
	status := resp.StatusCode

	// 如果资源支持范围下载, 则进行范围下载
	fmt.Println(resp.Header)
	fmt.Printf("StatusCode: %d, Accept-Ranges: %s, Content-Length: %s\n", status, header1, header2)
}
