package traces

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

var (
	sdkTraceProvider *sdktrace.TracerProvider
)

func InitOtelTraces(conn *grpc.ClientConn, ctx context.Context, res *resource.Resource) error {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("failed to create otel trace exporter: %s", err)
	}

	sdkTraceProvider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(sdkTraceProvider)
	return nil
}

func ShutdownOTelTraces(ctx context.Context) error {
	if sdkTraceProvider != nil {
		return sdkTraceProvider.Shutdown(ctx)
	}
	return nil
}

// func initTracer() func(context.Context) error {

// 	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
// 	if len(insecure) > 0 {
// 		secureOption = otlptracegrpc.WithInsecure()
// 	}

// 	exporter, err := otlptrace.New(
// 		context.Background(),
// 		otlptracegrpc.NewClient(
// 			secureOption,
// 			otlptracegrpc.WithEndpoint(collectorURL),
// 		),
// 	)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	resources, err := resource.New(
// 		context.Background(),
// 		resource.WithAttributes(
// 			attribute.String("service.name", serviceName),
// 			attribute.String("library.language", "go"),
// 		),
// 	)
// 	if err != nil {
// 		log.Printf("Could not set resources: ", err)
// 	}

// 	otel.SetTracerProvider(
// 		sdktrace.NewTracerProvider(
// 			sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 			sdktrace.WithBatcher(exporter),
// 			sdktrace.WithResource(resources),
// 		),
// 	)
// 	return exporter.Shutdown
// }
