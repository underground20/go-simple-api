package storage

import (
	"app/internal/department/models"
	"context"
	"errors"
	"fmt"
)

type Storage interface {
	Insert(ctx context.Context, e *models.Department) error
	Get(ctx context.Context, id int) (models.Department, error)
	Update(ctx context.Context, departmentId int, employeeId int) error
	GetAll(ctx context.Context) ([]models.Department, error)
	ChangeRoot(ctx context.Context, departmentId int, newRootId int) error
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
