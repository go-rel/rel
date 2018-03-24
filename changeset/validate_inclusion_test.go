package changeset

import (
	"fmt"
	"testing"
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

			ValidateInclusion(ch, "field", 1, 2.0, "c")

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidateInclusionError(t *testing.T) {
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

			ValidateInclusion(ch, "field", 2, 3.0, "d")

			if ch.Errors().Error() != "field must be one of [2 3 d]" {
				t.Error(`Expected "field must be one of [2 3 d]" but got`, ch.Errors().Error())
			}
		})
	}
}

func TestValidateInclusionMissing(t *testing.T) {
	ch := &Changeset{}
	ValidateInclusion(ch, "field", 5)

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
