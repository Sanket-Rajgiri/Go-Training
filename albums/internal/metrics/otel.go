package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

var (
	AlbumsAdded         metric.Int64Counter
	HttpRequestDuration metric.Float64Histogram
	sdkMeterProvider    *sdkmetric.MeterProvider
)

func InitOTelMetrics(conn *grpc.ClientConn, ctx context.Context, res *resource.Resource) error {

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("failed to create otel metrics exporter: %s", err)
	}

	views := []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{
				Name: "http_request_duration_seconds",
				Kind: sdkmetric.InstrumentKindHistogram,
			},
			sdkmetric.Stream{
				Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
					Boundaries: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0},
					NoMinMax:   true,
				},
			},
		),
	}

	sdkMeterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithView(views...),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(sdkMeterProvider)
	meter := sdkMeterProvider.Meter("gin-app")
	HttpRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialise meter : %s", err)
	}
	AlbumsAdded, err = meter.Int64Counter(
		"albums_added_count",
		metric.WithDescription("Counter of Albums Added"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialise meter : %s", err)
	}
	return nil
}

func ShutdownOTelMetrics(ctx context.Context) error {
	if sdkMeterProvider != nil {
		return sdkMeterProvider.Shutdown(ctx)
	}
	return nil
}
