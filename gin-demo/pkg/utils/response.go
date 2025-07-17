package utils

import (
	"github.com/gin-gonic/gin"
)

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, gin.H{
		"data": payload,
	})
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
	})
}
