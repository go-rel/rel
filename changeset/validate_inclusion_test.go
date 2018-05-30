package changeset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInclusion(t *testing.T) {
	tests := []interface{}{
		1,
		2.0,
		"c",
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateInclusion(ch, "field", []interface{}{1, 2.0, "c"})
			assert.Nil(t, ch.Errors())
		})
	}
}

func TestValidateInclusion_error(t *testing.T) {
	tests := []interface{}{
		1,
		1.0,
		"c",
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateInclusion(ch, "field", []interface{}{2, 3.0, "d"})
			assert.NotNil(t, ch.Errors())
			assert.Equal(t, "field must be one of [2 3 d]", ch.Error().Error())
		})
	}
}

func TestValidateInclusion_missing(t *testing.T) {
	ch := &Changeset{}
	ValidateInclusion(ch, "field", []interface{}{5})
	assert.Nil(t, ch.Errors())
}
