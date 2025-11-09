package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery 是自定义异常处理中间件，捕获 panic 并返回 500。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
