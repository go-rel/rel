package changeset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateExclusion(t *testing.T) {
	tests := []interface{}{
		2,
		3.0,
		"d",
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateExclusion(ch, "field", []interface{}{1, 2.0, "c"})
			assert.Nil(t, ch.Errors())
		})
	}
}

func TestValidateExclusion_error(t *testing.T) {
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

			ValidateExclusion(ch, "field", []interface{}{1, 2.0, "c"})
			assert.NotNil(t, ch.Errors())
			assert.Equal(t, "field must not be any of [1 2 c]", ch.Error().Error())
		})
	}
}

func TestValidateExclusion_missing(t *testing.T) {
	ch := &Changeset{}
	ValidateExclusion(ch, "field", []interface{}{5})
	assert.Nil(t, ch.Errors())
}
