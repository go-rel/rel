package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := Options{}
	opts.apply([]Option{
		Message("message"),
		Key("name_fk"),
		Exact(true),
		ChangeOnly(true),
		Required(true),
		SourceField("src"),
		EmptyValues("", 0),
	})

	assert.Equal(t, "message", opts.message)
	assert.Equal(t, "name_fk", opts.key)
	assert.Equal(t, true, opts.exact)
	assert.Equal(t, true, opts.changeOnly)
	assert.Equal(t, true, opts.required)
	assert.Equal(t, "src", opts.sourceField)
	assert.Equal(t, []interface{}{"", 0}, opts.emptyValues)
}
