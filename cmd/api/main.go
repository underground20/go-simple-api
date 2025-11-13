package main

import (
	"app/internal/config"
	"app/internal/employee/handler"
	db "app/internal/employee/storage/mongo"
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
	ctx := context.Background()
	cfg := config.MustLoad()
	client := clientConnect(logger, cfg, ctx)
	collection := client.Database(cfg.Dbname).Collection("employees")
	storage := db.NewStorage(collection, ctx)
	employeeHandler := handler.NewHandler(storage, logger)
	router := setupRouter(employeeHandler, logger)
	router.Run()
}

func setupRouter(handler *handler.Handler, logger *slog.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(loggingMiddleware.SlogMiddleware(logger))
	router.Use(metrics.PrometheusMiddleware())

	router.GET("/metrics", metrics.PrometheusHandler())

	router.POST("/employee/add", handler.CreateEmployee)
	router.GET("/employee/:id", handler.GetEmployee)
	router.PUT("/employee/:id", handler.UpdateEmployee)
	router.DELETE("/employee/:id", handler.DeleteEmployee)
	router.GET("/employees", handler.GetEmployees)

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
