package mongo

import (
	"app/internal/department/models"
	"app/internal/department/storage"
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

func (s *Storage) Insert(d *models.Department) error {
	result, err := s.collection.InsertOne(s.context, d)
	if err != nil {
		return err
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("could not convert inserted ID")
	}

	return nil
}

func (s *Storage) Update(departmentId int, employeeId int) error {
	filter := bson.M{"id": departmentId}
	update := bson.M{
		"$addToSet": bson.M{
			"employeeids": employeeId,
		},
	}
	result, err := s.collection.UpdateOne(s.context, filter, update)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("department with Id=%d not found", departmentId)
	}

	return nil
}

func (s *Storage) Get(id int) (models.Department, error) {
	var department models.Department
	filter := bson.M{"id": id}
	err := s.collection.FindOne(s.context, filter).Decode(&department)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Department{}, &storage.DepartmentNotFoundErr{Id: id}
		}
		return models.Department{}, err
	}

	return department, nil
}

func (s *Storage) GetAll() ([]models.Department, error) {
	cursor, err := s.collection.Find(s.context, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.context)

	var departments []models.Department
	if err = cursor.All(s.context, &departments); err != nil {
		return nil, err
	}

	return departments, nil
}

func (s *Storage) ChangeRoot(departmentId int, newRootId int) error {
	if _, err := s.Get(newRootId); err != nil {
		return err
	}

	filter := bson.M{"id": departmentId}
	update := bson.M{
		"$set": bson.M{
			"rootid": newRootId,
		},
	}
	result, err := s.collection.UpdateOne(s.context, filter, update)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if result.MatchedCount == 0 {
		return &storage.DepartmentNotFoundErr{Id: departmentId}
	}

	return nil
}
