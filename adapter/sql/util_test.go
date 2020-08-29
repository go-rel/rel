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

func TestToInt64(t *testing.T) {
	assert.Equal(t, int64(1), toInt64(int(1)))
	assert.Equal(t, int64(1), toInt64(int64(1)))
	assert.Equal(t, int64(1), toInt64(int32(1)))
	assert.Equal(t, int64(1), toInt64(int16(1)))
	assert.Equal(t, int64(1), toInt64(int8(1)))
	assert.Equal(t, int64(1), toInt64(uint(1)))
	assert.Equal(t, int64(1), toInt64(uint64(1)))
	assert.Equal(t, int64(1), toInt64(uint32(1)))
	assert.Equal(t, int64(1), toInt64(uint16(1)))
	assert.Equal(t, int64(1), toInt64(uint8(1)))
}
