package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			return
		}

		// 简化版，实际应该验证JWT token
		userID := uint(1) // 从token解析出用户ID
		c.Set("userID", userID)
		c.Next()
	}
}
