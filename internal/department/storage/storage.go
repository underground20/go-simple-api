package storage

import (
	"app/internal/department/models"
	"errors"
	"fmt"
)

type Storage interface {
	Insert(e *models.Department) error
	Get(id int) (models.Department, error)
	Update(departmentId int, employeeId int) error
	GetAll() ([]models.Department, error)
}

type DepartmentNotFoundErr struct {
	Id int
}

func (e *DepartmentNotFoundErr) Error() string {
	return fmt.Sprintf("department with id=%d not found", e.Id)
}

func IsDepartmentNotFound(err error) bool {
	var departmentNotFoundErr *DepartmentNotFoundErr
	ok := errors.As(err, &departmentNotFoundErr)
	return ok
}
