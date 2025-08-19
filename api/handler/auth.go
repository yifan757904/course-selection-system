package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	IDCard   string `json:"id_card" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=student teacher"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := service.RegisterInput{
		IDCard:   req.IDCard,
		Name:     req.Name,
		Password: req.Password,
		Role:     req.Role,
	}

	user, err := h.authService.Register(input)
	if err != nil {
		status := http.StatusBadRequest
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"id_card": user.IDCard,
		"name":    user.Name,
		"role":    user.Role,
	})
}

type LoginRequest struct {
	IDCard   string `json:"id_card" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := service.LoginInput{
		IDCard:   req.IDCard,
		Password: req.Password,
	}

	token, err := h.authService.Login(input)
	if err != nil {
		status := http.StatusUnauthorized
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": 3600,
		"token_type": "Bearer",
	})
}
