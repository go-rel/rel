package schema

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	var (
		record = struct {
			A string
			B *int
			C []byte
			D bool
			E []*float64
			F userDefined
			G time.Time
		}{}
		rt       = reflectTypePtr(record)
		expected = []reflect.Type{
			String,
			Int,
			Bytes,
			Bool,
			reflect.TypeOf([]float64{}),
			reflect.TypeOf(userDefined(0)),
			Time,
		}
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, InferTypes(record))

	_, cached = typesCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, expected, InferTypes(&record))
}

func TestTypes_usingInterface(t *testing.T) {
	var (
		record   = CustomSchema{}
		rt       = reflectTypePtr(record)
		expected = record.Types()
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, InferTypes(record))

	_, cached = typesCache.Load(rt)
	assert.False(t, cached)
}
