package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := Options{}
	opts.apply([]Option{
		Message("message"),
		Code(1000),
		Name("name_fk"),
		Exact(true),
		ChangeOnly(true),
	})

	assert.Equal(t, "message", opts.message)
	assert.Equal(t, 1000, opts.code)
	assert.Equal(t, "name_fk", opts.name)
	assert.Equal(t, true, opts.exact)
	assert.Equal(t, true, opts.changeOnly)
}
