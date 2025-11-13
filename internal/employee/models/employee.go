package models

type Employee struct {
	Id       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,max=255"`
	Sex      string `json:"sex" validate:"required,oneof=male female"`
	Age      int    `json:"age" validate:"required"`
	Position string `json:"position" validate:"required,max=255"`
	Salary   int    `json:"salary" validate:"required"`
}
