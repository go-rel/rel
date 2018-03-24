package changeset

import (
	"fmt"
	"testing"
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

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidatePatternError(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": "seafood",
		},
	}

	ValidatePattern(ch, "field", "boo.*")

	if ch.Errors().Error() != "field is invalid" {
		t.Error(`Expected "field is invalid" but got`, ch.Errors().Error())
	}
}

func TestValidatePatternMissing(t *testing.T) {
	ch := &Changeset{}
	ValidatePattern(ch, "field", "foo.*")

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
