package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequired(t *testing.T) {
	ch := &Changeset{
		values: map[string]interface{}{
			"field3": true,
		},
		changes: map[string]interface{}{
			"field1": 1,
			"field2": " 1 ",
		},
	}

	ValidateRequired(ch, []string{"field1", "field2", "field3"})
	assert.Nil(t, ch.Errors())
}

func TestValidateRequiredError(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field1": nil,
			"field2": "  ",
		},
	}

	ValidateRequired(ch, []string{"field1", "field2", "field3"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 3, len(ch.Errors()))
	assert.Equal(t, "field1 is required", ch.Errors()[0].Error())
	assert.Equal(t, "field2 is required", ch.Errors()[1].Error())
}
