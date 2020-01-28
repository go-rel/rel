package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractString(t *testing.T) {
	s := "Duplicate entry '1' for key 'slug'"
	assert.Equal(t, "slug", ExtractString(s, "key '", "'"))
}

func TestExtractString_notFound(t *testing.T) {
	s := "Duplicate entry '1' for field 'slug'"
	assert.Equal(t, "Duplicate entry '1' for field 'slug'", ExtractString(s, "key '", "'"))
}
