package main

import (
	"log"
	"os"

	"github.com/liuyifan1996/course-selection-system/cmd"
)

func main() {
	// 设置环境变量
	os.Setenv("JWT_SECRET_KEY", "3a1f8d7e4c9b2a5f6e8c3d0a7b4e5f2d1c8e3f6a9d2b5c4e7f8a1d3e6c9b2a5")
	os.Setenv("DSN", "root:root@tcp(127.0.0.1:3306)/cousle_sys?charset=utf8mb4&parseTime=True&loc=Local")

	r := cmd.Setup()

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Printf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}
