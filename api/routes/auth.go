package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		IDCard string `json:"id_card" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Rule   string `json:"rule" binding:"required,oneof=student teacher"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已存在
	var existing model.User
	if h.db.Where("id_card = ?", input.IDCard).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	user := model.User{
		IDCard: input.IDCard,
		Name:   input.Name,
		Rule:   input.Rule,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"id_card": user.IDCard,
		"name":    user.Name,
		"rule":    user.Rule,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		IDCard string `json:"id_card" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := h.db.Where("id_card = ?", input.IDCard).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 简化版token生成
	token := "generated-token-" + user.IDCard

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
			"rule": user.Rule,
		},
	})
}
