package schema

import (
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
)

type fields interface {
	Fields() map[string]int
}

var fieldsCache sync.Map

// InferFields from a struct.
func InferFields(record interface{}) map[string]int {
	if s, ok := record.(fields); ok {
		return s.Fields()
	}

	rt := reflectInternalType(record)
	// rt := reflectTypePtr(record)

	// check for cache
	if v, cached := fieldsCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	var (
		index  = 0
		fields = make(map[string]int, rt.NumField())
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = inferFieldName(sf)
		)

		if name != "" {
			fields[name] = index
			index++
		}
	}

	fieldsCache.Store(rt, fields)

	return fields
}

func inferFieldName(sf reflect.StructField) string {
	if tag := sf.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]

		if name == "-" {
			return ""
		}

		if name != "" {
			return name
		}
	}

	return snakecase.SnakeCase(sf.Name)
}

var fieldMappingCache sync.Map

func inferFieldMapping(record interface{}) map[string]int {
	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := fieldMappingCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	mapping := make(map[string]int, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = inferFieldName(sf)
		)

		if name != "" {
			mapping[name] = i
		}
	}

	fieldMappingCache.Store(rt, mapping)

	return mapping
}
