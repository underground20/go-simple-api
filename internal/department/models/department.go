package models

type Department struct {
	Id          int    `json:"id" validate:"required"`
	RootId      int    `json:"root_id" bson:"root_id"`
	Name        string `json:"name" validate:"required,max=255"`
	EmployeeIds []int  `json:"employee_ids" bson:"employee_ids"`
}
