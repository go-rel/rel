package changeset

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestPutAssoc_one(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": params.Map{
			"field4": 4,
			"field5": "5",
		},
	}

	chInner := changeInner(inner, input["field3"].(params.Map))

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": chInner,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": inner,
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	PutAssoc(ch, "field3", chInner)
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())

	PutAssoc(ch, "field3", inner)
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Errors()))
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestPutAssoc_many(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []params.Map{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	field3 := input["field3"].([]params.Map)

	chs := []*Changeset{
		changeInner(inner, field3[0]),
		changeInner(inner, field3[1]),
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": chs,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf([]Inner{}),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": []Inner(nil),
	}

	// with map assoc
	ch := Cast(data, input, []string{"field1", "field2"})
	PutAssoc(ch, "field3", chs)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}
