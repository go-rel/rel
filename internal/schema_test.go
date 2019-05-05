package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	var (
		record = struct {
			A string
			B *int
			C interface{} `db:"-"`
			D bool        `db:"D"`
			E []*float64  `db:"Ee,primary"`
		}{}
		rt             = reflectTypePtr(record)
		expectedFields = map[string]int{
			"a":  0,
			"b":  1,
			"D":  2,
			"Ee": 3,
		}
		expectedTypes = []reflect.Type{
			reflect.TypeOf(""),
			reflect.TypeOf(0),
			reflect.TypeOf(false),
			reflect.TypeOf([]float64{}),
		}
	)

	_, cached := schemasCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expectedFields, InferFields(record))
	assert.Equal(t, expectedTypes, InferTypes(record))

	_, cached = schemasCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, expectedFields, InferFields(&record))
	assert.Equal(t, expectedTypes, InferTypes(&record))
}

func TestSchema_usingInterface(t *testing.T) {
	var (
		record         = CustomSchema{}
		rt             = reflectTypePtr(record)
		expectedFields = record.Fields()
		expectedTypes  = record.Types()
	)

	_, cached := schemasCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expectedFields, InferFields(record))
	assert.Equal(t, expectedTypes, InferTypes(record))

	_, cached = schemasCache.Load(rt)
	assert.False(t, cached)
}
