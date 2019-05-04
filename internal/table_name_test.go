package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferTableName(t *testing.T) {
	type User struct{}
	record := User{}

	// should not be cached yet
	typ := reflectInternalType(record)
	_, cached := tableNamesCache.Load(typ)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(record)
	assert.Equal(t, "users", name)

	// cached
	_, cached = tableNamesCache.Load(typ)
	assert.True(t, cached)

	// infer table name using cache
	name = InferTableName(record)
	assert.Equal(t, "users", name)
}

func TestInferTableName_usingInterface(t *testing.T) {
	record := CustomSchema{}

	// should not be cached yet
	typ := reflectInternalType(record)
	_, cached := tableNamesCache.Load(typ)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(record)
	assert.Equal(t, "users", name)

	// never cache
	_, cached = tableNamesCache.Load(typ)
	assert.False(t, cached)
}
