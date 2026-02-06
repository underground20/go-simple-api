package department

import (
	"app/internal/department/handler"
	depStorage "app/internal/department/storage/mongo"
	empStorage "app/internal/employee/storage/mongo"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(router *gin.Engine, db *mongo.Database, logger *slog.Logger) {
	employeeStorage := empStorage.NewStorage(db.Collection("employees"))
	departmentStorage := depStorage.NewStorage(db.Collection("departments"))
	newHandler := handler.NewHandler(departmentStorage, employeeStorage, logger, validator.New())

	router.GET("/department/:id", newHandler.GetDepartment)
	router.POST("/department/add", newHandler.CreateDepartment)
	router.POST("department/add-employee", newHandler.AddEmployee)
	router.GET("/departments/tree", newHandler.GetTree)
	router.POST("/department/:id/change-root", newHandler.ChangeRoot)
}
