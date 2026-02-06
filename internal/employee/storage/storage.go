package storage

import (
	"app/internal/employee/models"
	"context"
	"errors"
	"fmt"
)

type Storage interface {
	Insert(ctx context.Context, e *models.Employee) error
	Get(ctx context.Context, id int) (models.Employee, error)
	Update(ctx context.Context, id int, e models.Employee) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) []models.Employee
	GetAllByIds(ctx context.Context, ids []int) []models.Employee
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
