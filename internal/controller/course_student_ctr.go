package controller

import (
	"net/http"

	"github.com/liuyifan1996/internal/service"

	"github.com/gin-gonic/gin"
)

type Course_StudentController struct {
	course_studentSvc *service.Course_StudentService
}

func NewCourse_StudentController(course_studentSvc *service.Course_StudentService) *Course_StudentController {
	return &Course_StudentController{course_studentSvc: course_studentSvc}
}

func (c *Course_StudentController) Course_StudentRegister(ctx *gin.Context) {
	var req service.Course_StudentRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course_student, err := c.course_studentSvc.Course_StudentRegister(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	ctx.JSON(http.StatusCreated, course_student)
}

func (c *Course_StudentController) GetCourse_StudentProfile(ctx *gin.Context) {
	// 从中间件获取用户ID
	courseid, exists := ctx.Get("courseid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	student_num, err := c.course_studentSvc.GetCourse_StudentProfile(ctx.Request.Context(), courseid.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, student_num)
}
