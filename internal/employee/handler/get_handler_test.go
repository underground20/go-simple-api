package handler_test

import (
	"app/internal/employee/handler"
	"app/internal/employee/models"
	"app/internal/employee/storage/cache"
	"app/internal/http/response"
	"app/lib/kafka"
	"app/lib/logger"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployee_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	producerMock := &kafka.ProducerMock{}
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger(), producerMock)

	router.GET("/employee/:id", createHandler.GetEmployee)

	employee := models.Employee{
		Id:       1,
		Name:     "John Doe",
		Sex:      "male",
		Age:      20,
		Position: "manager",
	}
	memoryStorage.Insert(context.Background(), &employee)

	req, _ := http.NewRequest(http.MethodGet, "/employee/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Employee
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, employee, resp)
}

func TestGetEmployee_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	producerMock := &kafka.ProducerMock{}
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger(), producerMock)
	router.GET("/employee/:id", createHandler.GetEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/employee/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "employee with id=2 not found", resp.Message)
}
