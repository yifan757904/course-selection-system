package utils

import "log"

// LogFatal 记录致命错误并退出程序
func LogFatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// LogError 记录错误但不退出
func LogError(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}
