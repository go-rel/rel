package changeset

import (
	"testing"
)

func TestAddError(t *testing.T) {
	ch := &Changeset{}

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}

	AddError(ch, "field1", "field1 is required")

	if ch.Errors().Error() != "field1 is required" {
		t.Error(`Expected "field1 is required" but got`, ch.Errors().Error())
	}

	AddError(ch, "field2", "field2 is not valid")

	if ch.Errors().Error() != "field1 is required, field2 is not valid" {
		t.Error(`Expected "field1 is required, field2 is not valid" but got`, ch.Errors().Error())
	}
}
