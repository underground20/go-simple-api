package cache

import (
	"app/employee/models"
	"errors"
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

func (s *MemoryStorage) Insert(e *models.Employee) {
	s.Lock()
	e.Id = s.counter
	s.data[e.Id] = *e

	s.counter++
	s.Unlock()
}

func (s *MemoryStorage) Delete(id int) {
	s.Lock()
	delete(s.data, id)
	s.Unlock()
}

func (s *MemoryStorage) Get(id int) (models.Employee, error) {
	s.Lock()
	defer s.Unlock()

	employee, exists := s.data[id]
	if !exists {
		return models.Employee{}, errors.New("employee not found")
	}

	return employee, nil
}

func (s *MemoryStorage) GetAll() []models.Employee {
	employees := make([]models.Employee, 0, len(s.data))
	for _, value := range s.data {
		employees = append(employees, value)
	}

	return employees
}

func (s *MemoryStorage) Update(id int, e models.Employee) {
	s.Lock()
	s.data[id] = e
	s.Unlock()
}
