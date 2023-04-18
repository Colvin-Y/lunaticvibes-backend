package main

import (
	"log"
	"net/http"
	"os"

	logger "github.com/Colvin-Y/lunaticvibes-backend/common/log"
	processor "github.com/Colvin-Y/lunaticvibes-backend/impl"
)

func main() {
	// 初始化 logger
	logger, err := logger.NewLogger("/var/log/lunaticvibes/lc.log")
	if err != nil {
		os.Exit(1)
	}
	defer logger.Close()

	// 注册路由处理函数
	scoreProcessor := &processor.ScoreProcessor{Logger: logger}
	http.HandleFunc("/score", scoreProcessor.ScoreHandlerSet)

	signUpProcessor := &processor.SignUpProcessor{Logger: logger}
	http.HandleFunc("/signup", signUpProcessor.SignUpHandlerSet)

	// 启动服务
	log.Print("Server started on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
