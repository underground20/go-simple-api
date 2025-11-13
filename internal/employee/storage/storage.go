package storage

import (
	"app/internal/employee/models"
	"errors"
	"fmt"
)

type Storage interface {
	Insert(e *models.Employee) error
	Get(id int) (models.Employee, error)
	Update(id int, e models.Employee) error
	Delete(id int) error
	GetAll() []models.Employee
	GetAllByIds(ids []int) []models.Employee
}

type EmployeeNotFoundErr struct {
	Id int
}

func (e *EmployeeNotFoundErr) Error() string {
	return fmt.Sprintf("employee with id=%d not found", e.Id)
}

func IsEmployeeNotFound(err error) bool {
	var employeeNotFoundErr *EmployeeNotFoundErr
	ok := errors.As(err, &employeeNotFoundErr)
	return ok
}
