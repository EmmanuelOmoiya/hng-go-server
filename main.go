package main

import (
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"time"
	"context"
	"hng.tech/backend-track/stage-1/config"
	"hng.tech/backend-track/stage-1/routes"
	"syscall"
)

func main(){
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config.LoadEnv()
	// Get and validate the configuration
	_, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config validation error: %s", err)
	}

	routes.InitGin()
	r := routes.New()

	srv := &http.Server{
		Addr:    ":8080",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()

	stop()
	fmt.Printf("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// logger.Error("Server forced to shutdown: ", err)
		log.Fatalf("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

