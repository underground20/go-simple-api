package app

import (
	"app/internal/config"
	"app/internal/department"
	"app/internal/employee"
	loggingMiddleware "app/internal/http/middleware/logger"
	"app/internal/http/middleware/metrics"
	"app/internal/notification"
	kafkaApi "app/lib/kafka"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	cfg      *config.Config
	server   *http.Server
	logger   *slog.Logger
	consumer *kafkaApi.Consumer
}

func New(logger *slog.Logger, cfg *config.Config) *App {
	client := clientConnect(logger, cfg, context.Background())
	db := client.Database(cfg.Dbname)
	router := setupRouter(logger)

	producer := kafkaApi.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	employee.Setup(router, db, logger, producer)
	department.Setup(router, db, logger)

	srv := &http.Server{
		Addr:    ":" + cfg.HostPort,
		Handler: router,
	}

	consumer := notification.ConfigureConsumer(logger, cfg)

	return &App{
		logger:   logger,
		cfg:      cfg,
		server:   srv,
		consumer: consumer,
	}
}

func (a *App) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Info("starting http server", slog.String("addr", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("http server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Info("starting kafka consumer")
		a.consumer.ReadMessages(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.logger.Info("shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.consumer.Close(); err != nil {
		a.logger.Error("failed to close kafka consumer", "error", err)
	}
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("server forced to shutdown", "error", err)
	}

	wg.Wait()
	a.logger.Info("server exiting")
}

func setupRouter(logger *slog.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(loggingMiddleware.SlogMiddleware(logger))
	router.Use(metrics.PrometheusMiddleware())
	router.GET("/metrics", metrics.PrometheusHandler())
	return router
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
