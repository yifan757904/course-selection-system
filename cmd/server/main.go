package main

import (
	"log"
	"os"

	"github.com/liuyifan1996/course-selection-system/cmd"
)

func main() {

	r := cmd.Setup()

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Printf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}
