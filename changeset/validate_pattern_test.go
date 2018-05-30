package changeset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePattern(t *testing.T) {
	tests := []interface{}{
		"seafood",
		1,
		2.0,
		false,
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidatePattern(ch, "field", "foo.*")
			assert.Nil(t, ch.Errors())
		})
	}
}

func TestValidatePattern_error(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": "seafood",
		},
	}

	ValidatePattern(ch, "field", "boo.*")
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field's format is invalid", ch.Error().Error())
}

func TestValidatePattern_missing(t *testing.T) {
	ch := &Changeset{}
	ValidatePattern(ch, "field", "foo.*")
	assert.Nil(t, ch.Errors())
}
