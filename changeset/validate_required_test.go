package changeset

import (
	"testing"
)

func TestValidateRequired(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field1": 1,
		},
	}

	ValidateRequired(ch, "field1")

	if ch.Errors() != nil {
		t.Error("Expected nil but got", ch.Errors())
	}
}

func TestValidateRequiredError(t *testing.T) {
	ch := &Changeset{}
	ValidateRequired(ch, "field1")

	if ch.Errors().Error() != "field1 is required" {
		t.Error(`Expected "field1 is required" but got`, ch.Errors().Error())
	}
}
