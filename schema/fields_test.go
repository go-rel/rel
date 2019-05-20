package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	var (
		record = struct {
			A string
			B *int
			C []byte     `db:",primary"`
			D bool       `db:"D"`
			E []*float64 `db:"-"`
		}{}
		rt       = reflectTypePtr(record)
		expected = map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
			"D": 3,
		}
	)

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, InferFields(record))

	_, cached = fieldsCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, expected, InferFields(&record))
}

func TestFields_usingInterface(t *testing.T) {
	var (
		record   = CustomSchema{}
		rt       = reflectTypePtr(record)
		expected = record.Fields()
	)

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, InferFields(record))

	_, cached = fieldsCache.Load(rt)
	assert.False(t, cached)
}
