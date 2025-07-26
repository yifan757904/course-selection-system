package main

import (
	"log"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/routes"
	"github.com/liuyifan1996/course-selection-system/config"
)

func main() {
	// 初始化数据库
	db, err := config.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&model.User{}, &model.Course{}, &model.Enrollment{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 设置路由
	r := routes.SetupRouter(db)

	// 启动服务器
	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
