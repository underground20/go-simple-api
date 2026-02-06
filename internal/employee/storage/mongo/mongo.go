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
}

func NewStorage(collection *mongo.Collection) *Storage {
	return &Storage{
		collection,
	}
}

func (s *Storage) Insert(ctx context.Context, e *models.Employee) error {
	result, err := s.collection.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("could not convert inserted ID")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	filter := bson.M{"id": id}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &storage.EmployeeNotFoundErr{Id: id}
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, id int) (models.Employee, error) {
	var employee models.Employee
	filter := bson.M{"id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Employee{}, &storage.EmployeeNotFoundErr{Id: id}
		}
		return models.Employee{}, err
	}

	return employee, nil
}

func (s *Storage) GetAll(ctx context.Context) []models.Employee {
	return s.getAllByFilter(ctx, bson.M{})
}

func (s *Storage) GetAllByIds(ctx context.Context, ids []int) []models.Employee {
	filter := bson.M{"id": bson.M{"$in": ids}}
	return s.getAllByFilter(ctx, filter)
}

func (s *Storage) Update(ctx context.Context, id int, e models.Employee) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": e}
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &storage.EmployeeNotFoundErr{Id: id}
	}

	return nil
}

func (s *Storage) getAllByFilter(ctx context.Context, filter any) []models.Employee {
	var employees []models.Employee
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return []models.Employee{}
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
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
