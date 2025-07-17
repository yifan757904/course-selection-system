package controller

import (
	"net/http"

	"github.com/liuyifan1996/internal/service"

	"github.com/gin-gonic/gin"
)

type TeacherController struct {
	teacherSvc *service.TeacherService
}

func NewTeacherController(teacherSvc *service.TeacherService) *TeacherController {
	return &TeacherController{teacherSvc: teacherSvc}
}

func (c *TeacherController) TeacherRegister(ctx *gin.Context) {
	var req service.TeacherRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := c.teacherSvc.TeacherRegister(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	ctx.JSON(http.StatusCreated, teacher)
}

func (c *TeacherController) GetTeacherProfile(ctx *gin.Context) {
	// 从中间件获取用户ID
	ID, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teacher, err := c.teacherSvc.GetTeacherProfile(ctx.Request.Context(), ID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, teacher)
}
