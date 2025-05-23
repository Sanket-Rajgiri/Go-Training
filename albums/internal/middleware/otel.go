package middleware

import (
	definedMetrics "albums/internal/metrics"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func OtelMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		definedMetrics.HttpRequestDuration.Record(
			c,
			duration,
			metric.WithAttributeSet(
				attribute.NewSet(
					attribute.String("method", c.Request.Method),
					attribute.String("path", c.FullPath()),
					semconv.HTTPResponseStatusCode(c.Writer.Status()),
				),
			),
		)
		if c.FullPath() == "/albums" && c.Request.Method == http.MethodPost && c.Writer.Status() == http.StatusCreated {
			definedMetrics.AlbumsAdded.Add(c.Request.Context(), 1)
		}
	}
}
