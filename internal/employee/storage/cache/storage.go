package cache

import (
	"app/internal/employee/models"
	"app/internal/employee/storage"
	"context"
	"sync"
)

type MemoryStorage struct {
	counter int
	data    map[int]models.Employee
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[int]models.Employee),
		counter: 1,
	}
}

func (s *MemoryStorage) Insert(_ context.Context, e *models.Employee) error {
	s.Lock()
	e.Id = s.counter
	s.data[e.Id] = *e

	s.counter++
	s.Unlock()

	return nil
}

func (s *MemoryStorage) Delete(_ context.Context, id int) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.data[id]; !ok {
		return &storage.EmployeeNotFoundErr{Id: id}
	}

	delete(s.data, id)

	return nil
}

func (s *MemoryStorage) Get(_ context.Context, id int) (models.Employee, error) {
	s.Lock()
	defer s.Unlock()

	employee, exists := s.data[id]
	if !exists {
		return models.Employee{}, &storage.EmployeeNotFoundErr{Id: id}
	}

	return employee, nil
}

func (s *MemoryStorage) GetAll(_ context.Context) []models.Employee {
	employees := make([]models.Employee, 0, len(s.data))
	for _, value := range s.data {
		employees = append(employees, value)
	}

	return employees
}

func (s *MemoryStorage) GetAllByIds(_ context.Context, _ []int) []models.Employee {
	employees := make([]models.Employee, 0, len(s.data))
	for _, value := range s.data {
		employees = append(employees, value)
	}

	return employees
}

func (s *MemoryStorage) Update(_ context.Context, id int, e models.Employee) error {
	s.Lock()
	s.data[id] = e
	s.Unlock()

	return nil
}
