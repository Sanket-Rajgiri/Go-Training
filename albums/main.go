package main

import (
	"albums/config"
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

	router, db := config.RouterSetup()
	// err := router.Run("localhost:8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB from GORM: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error Closing DB : %s", err.Error())
		} else {
			log.Println("DB Connection Closed")
		}
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("Server exiting")
}
