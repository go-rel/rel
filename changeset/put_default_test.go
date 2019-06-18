package changeset

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestPutDefaultEmptyChanges(t *testing.T) {
	ch := &Changeset{
		changes: make(map[string]interface{}),
		values: map[string]interface{}{
			"field1": 0,
			"field2": 5,
		},
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
			"field2": reflect.TypeOf(5),
		},
	}

	PutDefault(ch, "field1", 10)
	PutDefault(ch, "field2", 10)

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])
	assert.Nil(t, ch.Changes()["field2"])
}

func TestPutDefaultConditions(t *testing.T) {
	ch := &Changeset{
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
			"field2": reflect.TypeOf(0),
			"field3": reflect.TypeOf(0),
			"field4": reflect.TypeOf(0),
			"field5": reflect.TypeOf(0),
			"field6": reflect.TypeOf(0),
		},
		values: map[string]interface{}{
			"field1": 0,
			"field2": 1,
			"field3": 0,
			"field4": 1,
			"field5": 0,
			"field6": 1,
		},
		params: params.Map{
			"field1": 0,
			"field2": 0,
			"field3": 1,
			"field4": 1,
		},
		changes: map[string]interface{}{
			"field2": 0,
			"field3": 1,
		},
	}

	// existing | input | changes = changes after put default
	//     0    |   0   |   nil   =  nil
	//     1    |   0   |    0    =  0
	//     0    |   1   |    1    =  1
	//     1    |   1   |   nil   =  nil
	//     0    |  nil  |   nil   =  10
	//     1    |  nil  |   nil   =  nil

	tts := []struct {
		Field  string
		Result interface{}
	}{
		{"field1", nil},
		{"field2", 0},
		{"field3", 1},
		{"field4", nil},
		{"field5", 10},
		{"field6", nil},
	}

	for _, tt := range tts {
		t.Run(tt.Field, func(t *testing.T) {
			PutDefault(ch, tt.Field, 10)
			assert.Equal(t, tt.Result, ch.Changes()[tt.Field])
		})
	}
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

func TestPutDefaultNoValueWithInput(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field2": "input_val",
		},
		values: map[string]interface{}{},
		types: map[string]reflect.Type{
			"field2": reflect.TypeOf(""),
		},
	}

	PutDefault(ch, "field2", "must not be changed")

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, "input_val", ch.Changes()["field2"])
}

func TestPutDefaultNoValueNoChange(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{},
		values:  map[string]interface{}{},
		types: map[string]reflect.Type{
			"field2": reflect.TypeOf(""),
		},
	}

	PutDefault(ch, "field2", "must not be changed")

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, "must not be changed", ch.Changes()["field2"])
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
