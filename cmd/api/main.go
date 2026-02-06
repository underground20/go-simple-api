package main

import (
	"app/internal/config"
	"app/internal/department"
	"app/internal/employee"
	loggingMiddleware "app/internal/http/middleware/logger"
	"app/internal/http/middleware/metrics"
	formatErr "app/lib/logger"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := setupLogger()
	cfg := config.MustLoad()
	client := clientConnect(logger, cfg, context.Background())
	db := client.Database(cfg.Dbname)
	router := setupRouter(logger)
	employee.Setup(router, db, logger)
	department.Setup(router, db, logger)

	srv := &http.Server{
		Addr:    ":" + cfg.HostPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", formatErr.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", formatErr.Err(err))
	}

	logger.Info("Server exiting")
}

func setupRouter(logger *slog.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(loggingMiddleware.SlogMiddleware(logger))
	router.Use(metrics.PrometheusMiddleware())
	router.GET("/metrics", metrics.PrometheusHandler())
	return router
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func clientConnect(logger *slog.Logger, cfg *config.Config, ctx context.Context) *mongo.Client {
	clientOpts := options.Client().ApplyURI(cfg.Dsn)
	clientOpts.SetAuth(options.Credential{Username: cfg.DbUser, Password: cfg.DbPass, AuthMechanism: "SCRAM-SHA-256"})
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB:", err)
		os.Exit(1)
	}

	return client
}
