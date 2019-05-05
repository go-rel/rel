package internal

import (
	"reflect"
)

type CustomSchema struct {
	UUID  string
	Price int
}

func (c CustomSchema) TableName() string {
	return "_users"
}

func (c CustomSchema) PrimaryKey() (string, interface{}) {
	return "_uuid", c.UUID
}

func (c CustomSchema) Fields() map[string]int {
	return map[string]int{
		"_uuid":  0,
		"_price": 1,
	}
}

func (c CustomSchema) Types() []reflect.Type {
	return []reflect.Type{String, Int}
}
