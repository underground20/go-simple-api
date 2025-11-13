package department

import (
	"app/internal/department/handler"
	depStorage "app/internal/department/storage/mongo"
	empStorage "app/internal/employee/storage/mongo"
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(router *gin.Engine, db *mongo.Database, logger *slog.Logger, ctx context.Context) {
	employeeStorage := empStorage.NewStorage(db.Collection("employees"), ctx)
	departmentStorage := depStorage.NewStorage(db.Collection("departments"), ctx)
	newHandler := handler.NewHandler(departmentStorage, employeeStorage, logger)

	router.GET("/department/:id", newHandler.GetDepartment)
	router.POST("/department/add", newHandler.CreateDepartment)
	router.POST("department/add-employee", newHandler.AddEmployee)
}
