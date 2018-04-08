package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteChange(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field1": 10,
		},
		values: map[string]interface{}{
			"field1": 0,
		},
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
		},
	}

	assert.Nil(t, ch.Error())
	assert.Equal(t, 1, len(ch.Changes()))

	// delete change
	DeleteChange(ch, "field1")
	assert.Nil(t, ch.Error())
	assert.Equal(t, 0, len(ch.Changes()))
}
