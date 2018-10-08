package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutDefaultEmpty(t *testing.T) {
	ch := &Changeset{
		changes: make(map[string]interface{}),
		values: map[string]interface{}{
			"field1": 0,
		},
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
		},
	}

	PutDefault(ch, "field1", 10)

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])
}

func TestPutDefaultExisting(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field2": "default_val",
		},
		values: map[string]interface{}{
			"field2": "default_val",
		},
		types: map[string]reflect.Type{
			"field2": reflect.TypeOf(""),
		},
	}

	PutDefault(ch, "field2", "must not be changed")

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, "default_val", ch.Changes()["field2"])
}

func TestPutDefaultInvalid(t *testing.T) {
	ch := &Changeset{
		changes: make(map[string]interface{}),
		values: map[string]interface{}{
			"field3": 0,
		},
		types: map[string]reflect.Type{
			"field3": reflect.TypeOf(0),
		},
	}

	PutDefault(ch, "field3", "invalid value")
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
}
