package controller

import (
	"gin-demo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	courseSvc *service.CourseService
}

func NewCourseController(courseSvc *service.CourseService) *CourseController {
	return &CourseController{courseSvc: courseSvc}
}

func (c *CourseController) CourseRegister(ctx *gin.Context) {
	var req service.CourseRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := c.courseSvc.CourseRegister(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}

func (c *CourseController) GetCourseProfile(ctx *gin.Context) {
	// 从中间件获取用户ID
	ID, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	course, err := c.courseSvc.GetCourseProfile(ctx.Request.Context(), ID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, course)
}
