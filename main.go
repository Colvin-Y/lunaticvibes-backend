package main

import (
    "net/http"
	"log"
	"github.com/Colvin-Y/lunaticvibes-backend/impl"
)

func main(){
	// 注册路由处理函数
	http.HandleFunc("/", processor.Handler)

	// 启动服务
	log.Print("Server started on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}