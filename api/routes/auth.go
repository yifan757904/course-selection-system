package routes

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/pkg"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type Registerinput struct {
	IDCard   string `json:"id_card" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=student teacher"`
}

func (h *AuthHandler) Register(c *gin.Context) {

	var input Registerinput

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

	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码必须包含至少一个大写字母、一个小写字母和一个数字"})
		return
	}

	user := model.User{
		IDCard:   input.IDCard,
		Name:     input.Name,
		Password: input.Password,
		Role:     input.Role,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"id_card": user.IDCard,
		"name":    user.Name,
		"role":    user.Role,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		IDCard   string `json:"id_card" binding:"required"`
		Password string `json:"password" binding:"required"`
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

	if err := h.db.Where("id_card = ? and password = ?", input.IDCard, input.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := pkg.GenerateToken(user.IDCard, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "系统错误"})
		return
	}

	c.JSON(200, gin.H{
		"token":      token,
		"expires_in": 3600,
		"token_type": "Bearer",
	})
}
