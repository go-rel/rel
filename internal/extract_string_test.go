package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractString(t *testing.T) {
	s := "Duplicate entry '1' for key 'slug'"
	assert.Equal(t, "slug", ExtractString(s, "key '", "'"))
}
