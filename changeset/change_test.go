package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChange(t *testing.T) {
	var data struct {
		FIELD1 int `db:"field1"`
		Field2 string
		Field3 bool
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": false,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(false),
	}

	ch := Change(data)
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedValues, ch.Values())
	assert.Equal(t, expectedTypes, ch.Types())
	assert.Equal(t, map[string]interface{}{}, ch.Changes())

	ch = Change(data, map[string]interface{}{"field1": 2})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedValues, ch.Values())
	assert.Equal(t, expectedTypes, ch.Types())
	assert.Equal(t, map[string]interface{}{"field1": 2}, ch.Changes())
}
