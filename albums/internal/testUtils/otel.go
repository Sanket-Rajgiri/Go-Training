package testutils

import (
	"albums/internal/customlogs"
	definedMetrics "albums/internal/metrics"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	noopMetric "go.opentelemetry.io/otel/metric/noop"
	noopTrace "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func InitOtel() {
	otel.SetTracerProvider(noopTrace.NewTracerProvider())
	customlogs.OtelLogger = otelzap.New(zap.NewNop())
	definedMetrics.AlbumsAdded = noopMetric.Int64Counter{}
}
