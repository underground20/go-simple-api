package handler

import (
	"app/internal/employee/models"
	"app/internal/employee/storage"
	"app/internal/http/response"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewHandler(storage storage.Storage, logger *slog.Logger) *Handler {
	return &Handler{storage: storage, logger: logger}
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Message: "Invalid json format",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(employee); err != nil {
		resp := response.ValidationError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := h.storage.Insert(&employee)
	if err != nil {
		h.logger.Error("Failed to insert employee", err)
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Message: "Employee created successfully",
	})
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, response.Response{
			Message: err.Error(),
		})
		return
	}

	var employee models.Employee
	if err := c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Message: "Invalid json format",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(employee); err != nil {
		resp := response.ValidationError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = h.storage.Update(id, employee)
	if err != nil {
		if storage.IsEmployeeNotFound(err) {
			c.JSON(http.StatusNotFound, response.Response{Message: err.Error()})
			return
		}

		h.logger.Error("Failed to update employee", err)
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": employee.Id,
	})
}

func (h *Handler) GetEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, response.Response{
			Message: err.Error(),
		})
		return
	}

	employee, err := h.storage.Get(id)
	if err != nil {
		if storage.IsEmployeeNotFound(err) {
			c.JSON(http.StatusNotFound, response.Response{Message: err.Error()})
			return
		}

		h.logger.Error("Failed to get employee", err)
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) GetEmployees(c *gin.Context) {
	employees := h.storage.GetAll()
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Message: fmt.Sprintf("Failed to convert id to int %s\n", err.Error()),
		})
		return
	}

	err = h.storage.Delete(id)
	if err != nil {
		if storage.IsEmployeeNotFound(err) {
			c.JSON(http.StatusNotFound, response.Response{Message: err.Error()})
			return
		}

		h.logger.Error("Failed to delete employee", err)
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Message: fmt.Sprintf("Employee %d successfully deleted", id),
	})
}
