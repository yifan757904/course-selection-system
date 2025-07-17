package controller

import (
	"net/http"

	"github.com/liuyifan1996/internal/service"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	studentSvc *service.StudentService
}

func NewStudentController(studentSvc *service.StudentService) *StudentController {
	return &StudentController{studentSvc: studentSvc}
}

func (c *StudentController) StudentRegister(ctx *gin.Context) {
	var req service.StudentRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := c.studentSvc.StudentRegister(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	ctx.JSON(http.StatusCreated, student)
}

func (c *StudentController) GetStudentProfile(ctx *gin.Context) {
	// 从中间件获取用户ID
	ID, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	student, err := c.studentSvc.GetStudentProfile(ctx.Request.Context(), ID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, student)
}
