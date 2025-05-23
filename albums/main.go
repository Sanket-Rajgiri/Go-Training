package main

import (
	"albums/config"
	definedMetrics "albums/internal/metrics"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title						Go Gin Rest API
//	@version					1.0
//	@description				A rest API in Go using Gin framework
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer &lt;token&gt;**

func main() {
	var exitCode int
	router, db := config.RouterSetup(context.Background())
	// err := router.Run("localhost:8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("listen: ", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	if err := definedMetrics.ShutdownOTelMetrics(ctx); err != nil {
		log.Println("Metrics Exporter Shutdown : ", err)
	} else {
		log.Println("Metrics Exporter Closed")
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
		exitCode = 1
	} else {
		log.Println("HTTP server Shutdown")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB from GORM: %v", err)
		exitCode = 1
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error Closing DB : %s", err.Error())
			exitCode = 1
		} else {
			log.Println("DB Connection Closed")
		}
	}
	log.Println("Server exiting")
	os.Exit(exitCode)
}
