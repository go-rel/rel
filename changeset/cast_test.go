package changeset

import (
	"reflect"
	"testing"
)

func TestCast(t *testing.T) {
	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
		"field4": "ignore please",
	}

	expected := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
	}

	ch := Cast(params, "field1", "field2", "field3")

	if !reflect.DeepEqual(ch.Changes(), expected) {
		t.Error("Expected", expected, "but got", ch.Changes())
	}
}
