package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutChange(t *testing.T) {
	ch := &Changeset{
		changes: make(map[string]interface{}),
		values: map[string]interface{}{
			"field1": 0,
		},
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
		},
	}

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 0, len(ch.Changes()))

	PutChange(ch, "field1", 10)
	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])

	PutChange(ch, "field1", "10")
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Errors()))
	assert.Equal(t, "field1 is invalid", ch.Error().Error())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])
}
