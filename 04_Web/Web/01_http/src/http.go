package http_lib

import (
	"fmt"
	"net/http"
)

// 创建一个 HTTP Server
func HTTPServer() {
	// 第一种方法: 使用默认的
	// http.ListenAndServe("", nil)

	// 第二种方法
	server := &http.Server{
		Addr:    "",
		Handler: nil,
	}
	server.ListenAndServe()
}

// 创建一个 HTTPS Server
func HTTPSServer() {

}

// 使用单个处理器
func OneHandler() {

}

// 使用多个处理器
type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

type WorldHandler struct{}

func (h *WorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World!")
}

func MultiHandler() {
	hello := HelloHandler{}
	world := WorldHandler{}

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.Handle("/hello", &hello)
	http.Handle("/world", &world)

	server.ListenAndServe()
}

// 处理器函数
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func World(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World!")
}

func HandlerFunction() {
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/world", World)
}

func init() {

}
