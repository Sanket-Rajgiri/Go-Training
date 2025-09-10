package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	AlbumAddCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "album_adds",
			Help: "Total number of album Albums addes",
		},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.15, 0.25, 0.5, 0.75, 1.0},
		},
		[]string{"path", "method", "status"},
	)
)

func InitPrometheusMetrics() {
	prometheus.MustRegister(AlbumAddCounter)
	prometheus.MustRegister(HTTPDuration)
}
