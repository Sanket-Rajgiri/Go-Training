package middleware

import (
	"albums/internal/customlogs"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	// "go.opentelemetry.io/otel/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		// log.Printf("[%s] %s %d (%s)\n", method, path, status, duration)
		customlogs.OtelLogger.Info(fmt.Sprintf("[%s] %s %d (%s)\n", method, path, status, duration))

	}
}
