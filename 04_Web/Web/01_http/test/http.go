package test

import (
	"fmt"
	"net/http"
)

func DoHead() {
	// 看服务器是否支持 Range 下载; 可根据响应中的 Accept-Ranges 是否为 bytes 判断
	req, err := http.NewRequest(http.MethodHead, "http://httpbin.org/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	header1 := resp.Header.Get("Accept-Ranges")
	fmt.Println(header1)

	// 支持范围下载的话, 可以根据 Content-Length 字段获取待下载内容的大小
	header2 := resp.Header.Get("Content-Length")
	fmt.Println(header2)

	// 以10个线程为例下载文件；
	// 1）创建10个文件分片结构，用于下载文件内容
	// type filePart struct {
	// 	Index int    //文件分片的序号
	// 	From  int    //开始byte
	// 	To    int    //解决byte
	// 	Data  []byte //http下载得到的文件内容
	// }
	// 第一个分片: Index=0, From=(文件大小/10)*0, To=(文件大小/10)*1-1, Data=Range下载的文件内容
	// 第二个分片: Index=1, From=(文件大小/10)*1, To=(文件大小/10)*2-1, Data=Range下载的文件内容
	// ...
	// 第十个分片: Index=9, From=(文件大小/10)*9, To=(文件大小/10)*10-1, Data=Range下载的文件内容

	// 设置请求的 Range 为 filePart.From - filePart.To, 开 10 个 gorouting 发起 Range 请求下载文件部分内容

	// 合并文件内容并校验文件 hash

}
