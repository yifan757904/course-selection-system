package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuyifan1996/course-selection-system/pkg"
)

// 简化版认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌格式"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := pkg.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效令牌"})
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.UserRole)
		c.Next()
	}
}
