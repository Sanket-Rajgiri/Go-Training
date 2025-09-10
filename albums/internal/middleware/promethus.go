package middleware

import (
	"albums/internal/metrics"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPDuration.WithLabelValues(
			c.FullPath(),
			c.Request.Method,
			status,
		).Observe(duration)

		if c.FullPath() == "/albums" && c.Request.Method == http.MethodPost && c.Writer.Status() == http.StatusCreated {
			metrics.AlbumAddCounter.Inc()
		}
	}
}
