package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/routes"
	"github.com/liuyifan1996/course-selection-system/config"
	"github.com/liuyifan1996/course-selection-system/pkg"
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
		auth.GET("/teacher-courses", courseHandler.GetTeacherCourses)
		auth.POST("/courses/update/:id", courseHandler.UpdateCourse)

		// 选课相关
		auth.POST("/courses/:id/enroll", enrollHandler.Enroll)
		auth.GET("/student-courses", enrollHandler.GetStudentCourses)
		auth.DELETE("/courses/:id/enroll", enrollHandler.DeleteEnroll)
	}

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败: ", err)
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
