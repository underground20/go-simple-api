package handler

import (
	"app/internal/department/models"
	depStorage "app/internal/department/storage"
	emp "app/internal/employee/models"
	empStorage "app/internal/employee/storage"
	"app/internal/http/response"
	"app/lib/logger"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DepartmentResponse struct {
	Id        int            `json:"id"`
	RootId    int            `json:"root_id"`
	Name      string         `json:"name"`
	Employees []emp.Employee `json:"employees"`
}

type AddEmployeeResponse struct {
	DepartmentId int `json:"department_id"`
	EmployeeId   int `json:"employee_id"`
}

type Handler struct {
	departmentStorage depStorage.Storage
	employeeStorage   empStorage.Storage
	logger            *slog.Logger
}

func NewHandler(departmentStorage depStorage.Storage, empStorage empStorage.Storage, logger *slog.Logger) *Handler {
	return &Handler{departmentStorage: departmentStorage, employeeStorage: empStorage, logger: logger}
}

func (h *Handler) CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.BindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Message: "Invalid json format",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(department); err != nil {
		resp := response.ValidationError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := h.departmentStorage.Insert(&department)
	if err != nil {
		h.logger.Error("Failed to insert department", logger.Err(err))
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Message: "Department created successfully",
	})
}

func (h *Handler) AddEmployee(c *gin.Context) {
	var resp AddEmployeeResponse
	if err := c.BindJSON(&resp); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Message: "Invalid json format",
		})
		return
	}

	err := h.departmentStorage.Update(resp.DepartmentId, resp.EmployeeId)
	if err != nil {
		h.logger.Error("Failed to add employee to department", logger.Err(err))
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "Employee added successfully"})
}

func (h *Handler) GetDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to convert id param to int", logger.Err(err))
		c.JSON(http.StatusBadRequest, response.Response{
			Message: err.Error(),
		})
		return
	}

	department, err := h.departmentStorage.Get(id)
	if err != nil {
		if depStorage.IsDepartmentNotFound(err) {
			c.JSON(http.StatusNotFound, response.Response{Message: err.Error()})
			return
		}

		h.logger.Error("Failed to get department", logger.Err(err))
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	employees := h.employeeStorage.GetAllByIds(department.EmployeeIds)
	departmentResponse := DepartmentResponse{
		Id:        department.Id,
		RootId:    department.RootId,
		Name:      department.Name,
		Employees: employees,
	}

	c.JSON(http.StatusOK, departmentResponse)
}

func (h *Handler) GetTree(c *gin.Context) {
	departments, err := h.departmentStorage.GetAll()
	if err != nil {
		h.logger.Error("Failed to get all departments", logger.Err(err))
		c.JSON(http.StatusInternalServerError, response.UnhandledError())
		return
	}

	tree := BuildTree(departments)
	c.JSON(http.StatusOK, tree)
}
