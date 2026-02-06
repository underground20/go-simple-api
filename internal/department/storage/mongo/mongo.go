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
}

func NewStorage(collection *mongo.Collection) *Storage {
	return &Storage{
		collection,
	}
}

func (s *Storage) Insert(ctx context.Context, d *models.Department) error {
	result, err := s.collection.InsertOne(ctx, d)
	if err != nil {
		return err
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("could not convert inserted ID")
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, departmentId int, employeeId int) error {
	filter := bson.M{"id": departmentId}
	update := bson.M{
		"$addToSet": bson.M{
			"employeeids": employeeId,
		},
	}
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("department with Id=%d not found", departmentId)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, id int) (models.Department, error) {
	var department models.Department
	filter := bson.M{"id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&department)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Department{}, &storage.DepartmentNotFoundErr{Id: id}
		}
		return models.Department{}, err
	}

	return department, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]models.Department, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []models.Department
	if err = cursor.All(ctx, &departments); err != nil {
		return nil, err
	}

	return departments, nil
}

func (s *Storage) ChangeRoot(ctx context.Context, departmentId int, newRootId int) error {
	if _, err := s.Get(ctx, newRootId); err != nil {
		return err
	}

	filter := bson.M{"id": departmentId}
	update := bson.M{
		"$set": bson.M{
			"rootid": newRootId,
		},
	}
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if result.MatchedCount == 0 {
		return &storage.DepartmentNotFoundErr{Id: departmentId}
	}

	return nil
}
