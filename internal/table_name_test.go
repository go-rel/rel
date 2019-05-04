package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct{}

type UserWithTableName struct{}

func (u UserWithTableName) TableName() string {
	return "users"
}

func TestInferTableName(t *testing.T) {
	user := User{}

	// should not be cached yet
	typ := reflectInternalType(user)
	_, cached := tableNamesCache.Load(typ)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(user)
	assert.Equal(t, "users", name)

	// cached
	_, cached = tableNamesCache.Load(typ)
	assert.True(t, cached)

	// infer table name using cache
	name = InferTableName(user)
	assert.Equal(t, "users", name)
}

func TestInferTableName_withInterface(t *testing.T) {
	user := UserWithTableName{}

	// should not be cached yet
	typ := reflectInternalType(user)
	_, cached := tableNamesCache.Load(typ)
	assert.False(t, cached)

	// infer table name
	name := InferTableName(user)
	assert.Equal(t, "users", name)

	// never cache
	_, cached = tableNamesCache.Load(typ)
	assert.False(t, cached)
}
