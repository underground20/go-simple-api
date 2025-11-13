package models

type Employee struct {
	Id     int    `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Sex    string `json:"sex" validate:"required"`
	Age    int    `json:"age" validate:"required"`
	Salary int    `json:"salary" validate:"required"`
}
