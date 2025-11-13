package employee

import (
	"app/internal/employee/handler"
	mongodb "app/internal/employee/storage/mongo"
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(router *gin.Engine, db *mongo.Database, logger *slog.Logger, ctx context.Context) {
	collection := db.Collection("employees")
	storage := mongodb.NewStorage(collection, ctx)
	handler := handler.NewHandler(storage, logger)

	router.POST("/employee/add", handler.CreateEmployee)
	router.GET("/employee/:id", handler.GetEmployee)
	router.PUT("/employee/:id", handler.UpdateEmployee)
	router.DELETE("/employee/:id", handler.DeleteEmployee)
	router.GET("/employees", handler.GetEmployees)
}
