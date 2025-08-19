package cmd

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/handler"
	"github.com/liuyifan1996/course-selection-system/api/middleware"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
	"github.com/liuyifan1996/course-selection-system/api/service"
	"github.com/liuyifan1996/course-selection-system/config"
)

func Setup() *gin.Engine {
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

	// 初始化仓库
	authrepo := repository.NewGormAuthRepository(db)
	courserepo := repository.NewGormCourseRepository(db)
	enrollmentrepo := repository.NewEnrollmentRepository(db)

	// 初始化服务
	authService := service.NewAuthService(authrepo)
	courseService := service.NewCourseService(courserepo, authrepo)
	enrollmentService := service.NewEnrollmentService(enrollmentrepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	courseHandler := handler.NewCourseHandler(courseService)
	enrollHandler := handler.NewEnrollmentHandler(enrollmentService)

	// 设置路由
	r := gin.Default()

	// 公共路由
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// 需要认证的路由
	auth := r.Group("/").Use(middleware.AuthMiddleware())
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

	return r
}
