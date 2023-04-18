package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// 定义日志级别常量
const (
	LOG_LEVEL_INFO    = "[INFO]"
	LOG_LEVEL_WARNING = "[WARNING]"
	LOG_LEVEL_ERROR   = "[ERROR]"
)

// 定义 logger 结构体
type Logger struct {
	logFileName string
	file        *os.File
}

// 初始化 logger 实例
func NewLogger(logFileName string) (*Logger, error) {
	logger := &Logger{}
	// 打开日志文件
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return nil, err
	}
	logger.logFileName = logFileName
	logger.file = file
	return logger, nil
}

// Close 方法关闭日志文件
func (logger *Logger) Close() {
	logger.file.Close()
}

// Write 方法写入日志内容到指定日志文件
func (logger *Logger) Write(level string, msg string, calldepth int) {
	// 获取调用栈信息
	pc, file, line, _ := runtime.Caller(calldepth)
	funcName := runtime.FuncForPC(pc).Name()
	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// 将日志信息写入日志文件
	logLine := fmt.Sprintf("%s %s %s:%d:%s(): %s\n", currentTime, level, file, line, funcName, msg)
	logger.file.WriteString(logLine)
	// 将日志信息打印到控制台
	fmt.Print(logLine)
}

// Info 方法打印 INFO 级别的日志信息
func (logger *Logger) Info(msg string) {
	logger.Write(LOG_LEVEL_INFO, msg, 2)
}

// Warning 方法打印 WARNING 级别的日志信息
func (logger *Logger) Warning(msg string) {
	logger.Write(LOG_LEVEL_WARNING, msg, 2)
}

// Error 方法打印 ERROR 级别的日志信息
func (logger *Logger) Error(msg string) {
	logger.Write(LOG_LEVEL_ERROR, msg, 2)
}
