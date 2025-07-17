package main

import (
	"gin-demo/internal/config"
	"gin-demo/internal/routes"
	"gin-demo/pkg/database"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := database.InitMySQL(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 创建Gin实例并注册路由
	r := routes.SetupRouter(db)

	// 启动服务
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
