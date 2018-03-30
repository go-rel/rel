package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequired(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field1": 1,
		},
	}

	ValidateRequired(ch, []string{"field1"})
	assert.Nil(t, ch.Errors())
}

func TestValidateRequiredError(t *testing.T) {
	ch := &Changeset{}
	ValidateRequired(ch, []string{"field1"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field1 is required", ch.Error().Error())
}
