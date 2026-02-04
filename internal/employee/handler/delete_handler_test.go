package handler_test

import (
	"app/internal/employee/handler"
	"app/internal/employee/models"
	"app/internal/employee/storage/cache"
	"app/internal/http/response"
	"app/lib/logger"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEmployee_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger())

	router.DELETE("/employee/:id", createHandler.DeleteEmployee)

	employee := models.Employee{
		Id:       1,
		Name:     "John Doe",
		Sex:      "male",
		Age:      20,
		Position: "manager",
		Salary:   20000,
	}
	memoryStorage.Insert(&employee)

	req, _ := http.NewRequest(http.MethodDelete, "/employee/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.Response
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Employee 1 successfully deleted", resp.Message)
	_, err := memoryStorage.Get(1)
	assert.NotNil(t, err)
}

func TestDeleteEmployee_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger())
	router.DELETE("/employee/:id", createHandler.DeleteEmployee)

	req, _ := http.NewRequest(http.MethodDelete, "/employee/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp response.Response
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "employee with id=1 not found", resp.Message)
}
