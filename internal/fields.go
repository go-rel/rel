package internal

import (
	"reflect"
	"sync"
)

type schema interface {
	Fields() map[string]int
	Types() []reflect.Type
	Values() []interface{}
}

var schemasCache sync.Map

func InferFields(record interface{}) map[string]int {
	if s, ok := record.(schema); ok {
		return s.Fields()
	}

	fields, _ := inferFieldAndTypes(record)
	return fields
}

func InferTypes(record interface{}) []reflect.Type {
	if s, ok := record.(schema); ok {
		return s.Types()
	}

	_, types := inferFieldAndTypes(record)
	return types
}

func inferFieldAndTypes(record interface{}) (map[string]int, []reflect.Type) {
	return nil, nil
}

func InferValues(record interface{}) []interface{} {
	// TODO: handle different use case
	// changeset needs non ptr values
	// scanner needs ptr values
	if s, ok := record.(schema); ok {
		return s.Values()
	}

	return nil
}
