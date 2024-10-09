package log

import (
	"io"
	"log"
	"os"
)

func InitLogger(logName string) *log.Logger {
	// 创建一个日志文件
	logFile, err := os.OpenFile(logName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	// 创建一个多写入器，将日志写入文件和标准输出
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	// 设置日志输出为多写入器
	return log.New(multiWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
}
