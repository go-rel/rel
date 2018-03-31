package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsMessage(t *testing.T) {
	opts := Options{}
	opts.apply([]Option{
		Message("message"),
	})

	assert.Equal(t, "message", opts.message)
}
