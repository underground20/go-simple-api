package messages

import (
	"encoding/json"
	"strconv"
)

type EmployeeCreated struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
}

func (e EmployeeCreated) Key() []byte {
	return []byte(strconv.Itoa(e.Id))
}

func (e EmployeeCreated) Value() []byte {
	data, _ := json.Marshal(e)
	return data
}
