package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsMessage(t *testing.T) {
	opts := Options{}
	opts.Apply([]Option{Message("message")})
	assert.Equal(t, opts.Message, "message")
}
