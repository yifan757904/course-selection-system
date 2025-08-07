package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/routes"
	"github.com/liuyifan1996/course-selection-system/config"
	"github.com/liuyifan1996/course-selection-system/pkg"
)

func main() {
	// 设置环境变量
	os.Setenv("JWT_SECRET_KEY", "3a1f8d7e4c9b2a5f6e8c3d0a7b4e5f2d1c8e3f6a9d2b5c4e7f8a1d3e6c9b2a5")
	os.Setenv("DSN", "root:root@tcp(127.0.0.1:3306)/cousle_sys?charset=utf8mb4&parseTime=True&loc=Local")

	// 初始化数据库
	db, err := config.InitDB()

	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&model.User{}, &model.Course{}, &model.Enrollment{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		os.Exit(1)
	}

	// 初始化处理器
	authHandler := routes.NewAuthHandler(db)
	courseHandler := routes.NewCourseHandler(db)
	enrollHandler := routes.NewEnrollmentHandler(db)

	// 设置路由
	r := gin.Default()

	// 公共路由
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// 需要认证的路由
	auth := r.Group("/").Use(authMiddleware())
	{
		// 课程相关
		auth.POST("/courses/create", courseHandler.CreateCourse)
		auth.GET("/courses", courseHandler.GetCourses)
		auth.DELETE("/courses/:id", courseHandler.DeleteCourse)
		auth.GET("/courses-teacherid/:id", courseHandler.GetTeacherCourses)
		auth.GET("/courses-teachername/:teachername", courseHandler.GetCoursesByTeacherName)
		auth.GET("/courses-coursename/:coursename", courseHandler.GetCoursesByCourseName)
		auth.POST("/courses/update/:id", courseHandler.UpdateCourse)

		// 选课相关
		auth.POST("/courses/:id/enroll", enrollHandler.Enroll)
		auth.GET("/student-courses", enrollHandler.GetStudentCourses)
		auth.DELETE("/courses/:id/enroll", enrollHandler.DeleteEnroll)
	}

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Printf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}

// 简化版认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌格式"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := pkg.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效令牌"})
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.UserRole)
		c.Next()
	}
}
