package storage

import "app/employee/models"

type Storage interface {
	Insert(e *models.Employee)
	Get(id int) (models.Employee, error)
	Update(id int, e models.Employee)
	Delete(id int)
	GetAll() []models.Employee
}
