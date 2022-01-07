package main

import (
	"app/employee/handler"
	"app/employee/storage/cache"

	"github.com/gin-gonic/gin"
)

func main() {
	memoryStorage := cache.NewMemoryStorage()
	handler := handler.NewHandler(memoryStorage)
	router := gin.Default()

	router.POST("/employee/add", handler.CreateEmployee)
	router.GET("/employee/:id", handler.GetEmployee)
	router.PUT("/employee/:id", handler.UpdateEmployee)
	router.DELETE("/employee/:id", handler.DeleteEmployee)
	router.GET("/employees", handler.GetEmployees)

	router.Run()
}
