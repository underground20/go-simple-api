package handler_test

import (
	"app/internal/employee/handler"
	"app/internal/employee/models"
	"app/internal/employee/storage/cache"
	"app/internal/http/response"
	"app/lib/kafka"
	"app/lib/logger"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmployee_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	producerMock := &kafka.ProducerMock{}
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger(), producerMock)

	router.POST("/employees", createHandler.CreateEmployee)

	employee := models.Employee{
		Id:       1,
		Name:     "John Doe",
		Sex:      "male",
		Age:      20,
		Position: "manager",
		Salary:   20000,
	}
	jsonData, _ := json.Marshal(employee)

	req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.Response
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Employee created successfully", resp.Message)
	emp, _ := memoryStorage.Get(req.Context(), 1)
	assert.Equal(t, "John Doe", emp.Name)
}

func TestCreateEmployee_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	producerMock := &kafka.ProducerMock{}
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger(), producerMock)

	router.POST("/employees", createHandler.CreateEmployee)

	invalidJSON := []byte(`{"name": "John Doe`)

	req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp response.Response
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Invalid json format", resp.Message)
	assert.Empty(t, memoryStorage.GetAll(req.Context()))
}

func TestCreateEmployee_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	memoryStorage := cache.NewMemoryStorage()
	producerMock := &kafka.ProducerMock{}
	createHandler := handler.NewHandler(memoryStorage, logger.NewDiscardLogger(), producerMock)
	router.POST("/employees", createHandler.CreateEmployee)

	tests := []struct {
		name       string
		input      string
		wantStatus int
		wantMsg    string
	}{
		{
			"Without salary", `{"id": 1, "name":"John", "sex": "male", "age": 20, "position": "manager"}`,
			http.StatusBadRequest,
			"Validation error: field 'salary' - required",
		},
		{
			"Incorrect sex", `{"id": 1, "name":"John", "sex": "t", "age": 20, "position": "manager", "salary": 10000}`,
			http.StatusBadRequest,
			"Validation error: field 'sex' - oneof",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer([]byte(tt.input)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var resp response.Response
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, tt.wantMsg, resp.Message)
			assert.Empty(t, memoryStorage.GetAll(req.Context()))
		})
	}
}
