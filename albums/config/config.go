package config

import (
	_ "albums/docs" // this line is REQUIRED for Swagger to find the docs package
	customlogs "albums/internal/customlogs"
	"albums/internal/database"
	"albums/internal/handlers"
	"albums/internal/metrics"
	"albums/internal/middleware"
	"albums/internal/service"
	"albums/internal/traces"
	"albums/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	dbTypeMap = map[string]string{
		"dev":   "mysql",
		"local": "sqlite",
	}
	serviceName = semconv.ServiceNameKey.String("gin-app")
)

func initGRPCConn(endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func loadEnv() (map[string]string, error) {
	envVariables := make(map[string]string)
	envName, exists := os.LookupEnv("ENV")
	if !exists {
		envName = "local"
		log.Println("setting envName = local")
	}
	envVariables["DB_TYPE"] = dbTypeMap[envName]
	if envName != "local" {
		requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_NAME", "COLLECTOR_ENDPOINT"}
		for _, key := range requiredVars {
			value := os.Getenv(key)
			if value == "" {
				return nil, fmt.Errorf("missing required environment variable: %s", key)
			}
			envVariables[key] = value
		}
	}

	return envVariables, nil
}

func initDB(envVars map[string]string) (*gorm.DB, error) {
	var db *gorm.DB

	if envVars["DB_TYPE"] == "local" {
		sqlite, err := database.SqliteConnect()
		if err != nil {
			return nil, err
		}
		err = database.SqliteInit(sqlite)
		if err != nil {
			return nil, err
		}
		db = sqlite
	} else {
		mysql, err := database.MysqlConnect(envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"])
		if err != nil {
			return nil, err
		}
		db = mysql
	}
	return db, nil
}
func RouterSetup(ctx context.Context) (*gin.Engine, *gorm.DB) {
	envVars, err := loadEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	db, err := initDB(envVars)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// metrics.InitPrometheusMetrics()

	grpcConn, err := initGRPCConn(envVars["COLLECTOR_ENDPOINT"])
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
