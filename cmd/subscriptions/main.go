package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alexputin/subscriptions/internal/config"
	"github.com/alexputin/subscriptions/internal/db"
	"github.com/alexputin/subscriptions/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.MustLoadConfig()
	config := config.Get()

	// Create database connection
	db, err := db.CreatePostgresConnection(config.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// TODO: Create repository instances
	// TODO: Create service instances

	app := echo.New()

	if config.Environment == "dev" {
		app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${time_rfc3339}\t| ${status} |  ${method}\t| ${uri}\n",
		}))
	}
	app.Use(middleware.Recover())

	// Register routes
	api := handlers.NewSubscriptionsApiHandler(nil) // TODO: Replace nil with actual repository implementation
	api.RegisterRoutes(app)

	// Server graceful shutdown
	go func() {
		if err := app.Start(config.ServerAddress); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
