package main

import (
	"app/internal/config"
	"app/internal/department"
	"app/internal/employee"
	loggingMiddleware "app/internal/http/middleware/logger"
	"app/internal/http/middleware/metrics"
	"context"
	"log/slog"
	"os"

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
	router.Run()
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
