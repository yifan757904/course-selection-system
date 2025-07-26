package controller

import (
	"net/http"

	"github.com/liuyifan1996/course-selection-system/api/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterRequest struct {
	Name   string `json:"name" binding:"required"`
	IDCard string `json:"id_card" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=student teacher"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Register(req.Name, req.IDCard, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	IDCard string `json:"id_card" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, user, err := c.authService.Login(req.IDCard)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token := rule
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
