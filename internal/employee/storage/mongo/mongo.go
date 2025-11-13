package mongo

import (
	"app/internal/employee/models"
	"app/internal/employee/storage"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	collection *mongo.Collection
	context    context.Context
}

func NewStorage(collection *mongo.Collection, context context.Context) *Storage {
	return &Storage{
		collection,
		context,
	}
}

func (s *Storage) Insert(e *models.Employee) error {
	result, err := s.collection.InsertOne(s.context, e)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("could not convert inserted ID")
	}

	return nil
}

func (s *Storage) Delete(id int) error {
	filter := bson.M{"id": id}
	result, err := s.collection.DeleteOne(s.context, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &storage.EmployeeNotFoundErr{Id: id}
	}

	return nil
}

func (s *Storage) Get(id int) (models.Employee, error) {
	var employee models.Employee
	filter := bson.M{"id": id}
	err := s.collection.FindOne(s.context, filter).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Employee{}, &storage.EmployeeNotFoundErr{Id: id}
		}
		return models.Employee{}, err
	}

	return employee, nil
}

func (s *Storage) GetAll() []models.Employee {
	return s.getAllByFilter(bson.M{})
}

func (s *Storage) GetAllByIds(ids []int) []models.Employee {
	filter := bson.M{"id": bson.M{"$in": ids}}
	return s.getAllByFilter(filter)
}

func (s *Storage) Update(id int, e models.Employee) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": e}
	result, err := s.collection.UpdateOne(s.context, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &storage.EmployeeNotFoundErr{Id: id}
	}

	return nil
}

func (s *Storage) getAllByFilter(filter any) []models.Employee {
	var employees []models.Employee
	cursor, err := s.collection.Find(s.context, filter)
	if err != nil {
		return []models.Employee{}
	}

	defer cursor.Close(s.context)

	for cursor.Next(s.context) {
		var employee models.Employee
		if err := cursor.Decode(&employee); err != nil {
			continue
		}

		employees = append(employees, employee)
	}

	if employees == nil {
		return []models.Employee{}
	}

	return employees
}
