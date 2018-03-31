package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsMessage(t *testing.T) {
	opts := Options{}
	opts.Apply([]Option{
		Message("message"),
		AllowError(true),
	})

	assert.Equal(t, "message", opts.message)
	assert.Equal(t, true, opts.allowError)
}
