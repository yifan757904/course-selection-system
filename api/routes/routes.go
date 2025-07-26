package routes

import (
	"github.com/liuyifan1996/course-selection-system/api/controller"
	"github.com/liuyifan1996/course-selection-system/api/middleware"
	"github.com/liuyifan1996/course-selection-system/api/repository"
	"github.com/liuyifan1996/course-selection-system/api/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo)
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	// 初始化控制器
	authCtrl := controller.NewAuthController(authService)
	courseCtrl := controller.NewCourseController(courseService)

	// 公共路由
	public := r.Group("/api")
	{
		public.POST("/register", authCtrl.Register)
		public.POST("/login", authCtrl.Login)
		public.GET("/courses", courseCtrl.GetAllCourses)
	}

	// 需要认证的路由
	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/courses/my", courseCtrl.GetMyCourses)
		private.DELETE("/courses/:id", courseCtrl.DeleteCourse)
		private.POST("/courses", courseCtrl.CreateCourse)
		private.POST("/courses/:id/enroll", courseCtrl.EnrollCourse)
	}

	return r
}
