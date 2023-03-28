package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	// 声明静态文件目录并将其指向我们刚刚创建的目录
	staticFileDirectory := http.Dir("./assets/")
	// 声明处理程序，将请求路由到它们各自的文件名。
	// 文件服务器被包装在 `stripPrefix` 方法中，因为我们想在查找文件时去掉“/assets/”前缀。
	// 例如，如果我们在浏览器中键入“/assets/index.html”，文件服务器将仅在上面声明的目录中查找“index.html”。
	// 如果我们不去除前缀，文件服务器将查找“./assets/assets/index.html”，并产生错误
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// “PathPrefix”方法充当匹配器，匹配所有以“/assets/”开头的路由
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	// 添加名人列表路由处理函数
	r.HandleFunc("/items", getItemsHandler).Methods("GET")
	r.HandleFunc("/items", createItemHandler).Methods("POST")

	return r
}

func main() {
	// 现在通过调用上面定义的 newRouter 构造函数来创建路由器。
	// 其余代码保持不变
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
