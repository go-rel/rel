package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	ch := Cast(params, []string{"field1", "field2", "field3"})
	assert.Equal(t, ch.Changes(), expected)
}
