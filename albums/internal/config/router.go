package config

import (
	_ "albums/docs" // this line is REQUIRED for Swagger to find the docs package
	"albums/internal/config/env"
	customlogs "albums/internal/customlogs"
	"albums/internal/handlers"
	"albums/internal/metrics"
	"albums/internal/middleware"
	"albums/internal/service"
	"albums/internal/traces"
	"albums/routes"
	"context"
	"log"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	serviceName = semconv.ServiceNameKey.String("gin-app")
)

func RouterSetup(ctx context.Context) (*gin.Engine, *gorm.DB) {
	err := env.LoadEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	db, err := InitDB()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// metrics.InitPrometheusMetrics()
	grpcConn, err := InitGRPCConn(env.COLLECTOR_ENDPOINT)
	if err != nil {
		log.Fatalln(err.Error())
	}
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			serviceName,
		),
	)
	if err != nil {
		log.Fatalln("failed to create otel resource: ", err)
	}
	err = customlogs.InitOtelLogger(grpcConn, ctx, res)
	if err != nil {
		log.Println(err)
	}
	err = metrics.InitOTelMetrics(grpcConn, ctx, res)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
	}
	err = traces.InitOtelTraces(grpcConn, ctx, res)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
	}
	router := gin.New()
	router.Use(
		middleware.LoggerMiddleware(),
		//  middleware.PrometheusMiddleware(),
		middleware.OtelMetricsMiddleware(),
		otelgin.Middleware(serviceName.Value.AsString()),
		gin.Recovery(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	albumService := &service.AlbumServiceImpl{DB: db}
	albumHandler := &handlers.AlbumHandler{AlbumService: albumService}
	routes.RegisterAlbumRoutes(router, albumHandler)

	loginService := &service.LoginServiceImpl{DB: db}
	loginHandler := &handlers.LoginHandler{LoginService: loginService}
	routes.RegisterLoginRoutes(router, loginHandler, loginService)
	return router, db
}
