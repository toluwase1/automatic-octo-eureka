package response

import (
	"github.com/gin-gonic/gin"
	"time"
)

func JSON(c *gin.Context, status int, data interface{}, errs []string, message string) {
	c.JSON(status, gin.H{
		"status":    status,
		"data":      data,
		"errors":    errs,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}
