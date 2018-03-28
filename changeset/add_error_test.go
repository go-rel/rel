package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddError(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Errors())

	AddError(ch, "field1", "field1 is required")
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field1 is required", ch.Errors().Error())

	AddError(ch, "field2", "field2 is not valid")
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field1 is required, field2 is not valid", ch.Errors().Error())
}
