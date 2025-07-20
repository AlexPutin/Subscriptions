package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	_ "github.com/alexputin/subscriptions/docs"
	"github.com/alexputin/subscriptions/internal/config"
	"github.com/alexputin/subscriptions/internal/db"
	"github.com/alexputin/subscriptions/internal/handlers"
	"github.com/alexputin/subscriptions/internal/repositories"
	"github.com/alexputin/subscriptions/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Инициализация zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application")

	config.MustLoadConfig()
	config := config.Get()

	// Create database connection
	db, err := db.CreatePostgresConnection(config.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database connection", zap.Error(err))
		}
	}()
	logger.Info("Database connection established")

	repo := repositories.NewPostgresUserSubscriptionRepository(db)
	service := services.NewUserSubscriptionService(repo)
	logger.Info("Repository and service initialized")

	app := echo.New()
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)
			logger.Info("request completed",
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", latency),
				zap.String("remote_ip", c.RealIP()),
			)
			return err
		}
	})

	app.Use(middleware.Recover())

	// Register routes
	api := handlers.NewSubscriptionsApiHandler(service, logger)
	api.RegisterRoutes(app)
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	logger.Info("Routes registered")

	// Server graceful shutdown
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", config.ServerAddress))
		if err := app.Start(config.ServerAddress); err != nil && err != http.ErrServerClosed {
			logger.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	shutdownCtx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server shutdown complete")
}
