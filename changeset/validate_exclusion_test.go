package changeset

import (
	"fmt"
	"testing"
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

			ValidateExclusion(ch, "field", 1, 2.0, "c")

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidateExclusionError(t *testing.T) {
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

			ValidateExclusion(ch, "field", 1, 2.0, "c")

			if ch.Errors().Error() != "field must not be any of [1 2 c]" {
				t.Error(`Expected "field must not be any of [1 2 c]" but got`, ch.Errors().Error())
			}
		})
	}
}

func TestValidateExclusionMissing(t *testing.T) {
	ch := &Changeset{}
	ValidateExclusion(ch, "field", 5)

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
