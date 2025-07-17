package routes

import (
	"database/sql"
	"gin-demo/internal/controller"
	"gin-demo/internal/middleware"
	"gin-demo/internal/repository"
	"gin-demo/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// 初始化依赖
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSvc)

	// 公共路由
	public := r.Group("/api")
	{
		public.POST("/register", userCtrl.Register)
	}

	// 需要认证的路由
	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/profile", userCtrl.GetProfile)
	}

	return r
}
