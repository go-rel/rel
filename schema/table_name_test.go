package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferTableName(t *testing.T) {
	type User struct{}

	var (
		record       = User{}
		rt           = reflectInternalType(record)
		expectedName = "users"
	)

	// should not be cached yet
	_, cached := tableNamesCache.Load(rt)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(record)
	assert.Equal(t, expectedName, name)

	// cached
	_, cached = tableNamesCache.Load(rt)
	assert.True(t, cached)

	// infer table name using cache
	name = InferTableName(record)
	assert.Equal(t, expectedName, name)
}

func TestInferTableName_usingInterface(t *testing.T) {
	var (
		record       = CustomSchema{}
		rt           = reflectInternalType(record)
		expectedName = "_users"
	)

	// should not be cached yet
	_, cached := tableNamesCache.Load(rt)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(record)
	assert.Equal(t, expectedName, name)

	// never cache
	_, cached = tableNamesCache.Load(rt)
	assert.False(t, cached)
}
