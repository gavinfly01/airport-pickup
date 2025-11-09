package http

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 是自定义日志中间件，记录请求信息。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s | %d | %v", method, path, status, latency)
	}
}
