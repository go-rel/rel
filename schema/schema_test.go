package schema

import (
	"database/sql"
	"reflect"
)

type CustomSchema struct {
	UUID  string
	Price int
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

func (c CustomSchema) Values() []interface{} {
	return []interface{}{c.UUID, c.Price}
}

func (c *CustomSchema) Scanners(fields []string) []interface{} {
	var (
		scanners  = make([]interface{}, len(fields))
		tempValue = sql.RawBytes{}
	)

	for index, field := range fields {
		switch field {
		case "_uuid":
			scanners[index] = Nullable(&c.UUID)
		case "_price":
			scanners[index] = Nullable(&c.Price)
		default:
			scanners[index] = &tempValue
		}
	}

	return scanners
}
