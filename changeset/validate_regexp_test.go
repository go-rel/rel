package changeset

import (
	"fmt"
	"regexp"
	"testing"
)

func TestValidateRegexp(t *testing.T) {
	exp := regexp.MustCompile(`foo.*`)

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

			ValidateRegexp(ch, "field", exp)

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidateRegexpError(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": "seafood",
		},
	}

	ValidateRegexp(ch, "field", regexp.MustCompile(`boo.*`))

	if ch.Errors().Error() != "field is invalid" {
		t.Error(`Expected "field is invalid" but got`, ch.Errors().Error())
	}
}

func TestValidateRegexpMissing(t *testing.T) {
	ch := &Changeset{}
	ValidateRegexp(ch, "field", regexp.MustCompile(`foo.*`))

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
