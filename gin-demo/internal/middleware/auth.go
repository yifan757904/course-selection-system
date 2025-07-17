package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简化版认证，实际项目应使用JWT等
		token := c.GetHeader("Authorization")
		if token != "secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// 模拟用户ID，实际应从token解析
		c.Set("userID", int64(1))
		c.Next()
	}
}
