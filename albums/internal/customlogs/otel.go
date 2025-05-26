package customlogs

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var OtelLogger *zap.Logger
var otelLogProvider *otellog.LoggerProvider

func InitOtelLogger(conn *grpc.ClientConn, ctx context.Context, res *resource.Resource) error {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("failed to create otel trace exporter: %s", err)
	}
	processor := otellog.NewBatchProcessor(exporter)
	otelLogProvider = otellog.NewLoggerProvider(otellog.WithResource(res), otellog.WithProcessor(processor))
	OtelLogger = zap.New(
		zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
			otelzap.NewCore("gin-app", otelzap.WithLoggerProvider(otelLogProvider)),
		),
	)
	return nil
}

func ShutdownLogger(ctx context.Context) error {
	if otelLogProvider != nil {
		return otelLogProvider.Shutdown(ctx)
	}
	return nil
}
