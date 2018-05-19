package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	var data struct {
		FIELD1 int `db:"field1"`
		Field2 string
		Field3 bool
	}

	expectedChanges := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": false,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(false),
	}

	ch := Convert(data)
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedTypes, ch.types)
}
