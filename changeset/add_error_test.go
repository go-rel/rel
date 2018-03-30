package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddError(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 0, len(ch.Errors()))

	AddError(ch, "field1", "field1 is required")
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Errors()))
	assert.Equal(t, "field1 is required", ch.Error().Error())

	AddError(ch, "field2", "field2 is not valid")
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 2, len(ch.Errors()))
	assert.Equal(t, "field2 is not valid", ch.Errors()[1].Error())
}
