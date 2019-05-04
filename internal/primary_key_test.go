package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferPrimaryKey(t *testing.T) {
	record := struct {
		ID   uint
		Name string
	}{
		ID: 1,
	}

	// should not be cached yet
	typ := reflectInternalType(record)
	_, cached := primaryKeysCache.Load(typ)
	assert.False(t, cached)

	// infer primary key
	field, value := InferPrimaryKey(record, true)
	assert.Equal(t, "id", field)
	assert.Equal(t, uint(1), value)

	// cached
	_, cached = primaryKeysCache.Load(typ)
	assert.True(t, cached)

	record.ID = 2

	// infer primary key using cache
	field, value = InferPrimaryKey(record, true)
	assert.Equal(t, "id", field)
	assert.Equal(t, uint(2), value)
}

func TestInferPrimaryKey_usingInterface(t *testing.T) {
	record := CustomSchema{
		UUID: "abc123",
	}

	// should not be cached yet
	typ := reflectInternalType(record)
	_, cached := primaryKeysCache.Load(typ)
	assert.False(t, cached)

	// infer primary key
	field, value := InferPrimaryKey(record, true)
	assert.Equal(t, "uuid", field)
	assert.Equal(t, record.UUID, value)

	// never cache
	_, cached = primaryKeysCache.Load(typ)
	assert.False(t, cached)
}

func TestInferPrimaryKey_usingTag(t *testing.T) {
	record := struct {
		ID         uint
		ExternalID int `db:",primary"`
		Name       string
	}{
		ExternalID: 12345,
	}

	// infer primary key
	field, value := InferPrimaryKey(record, true)
	assert.Equal(t, "external_id", field)
	assert.Equal(t, 12345, value)
}

func TestInferPrimaryKey_usingTagAmdCustomName(t *testing.T) {
	record := struct {
		ID         uint
		ExternalID int `db:"partner_id,primary"`
		Name       string
	}{
		ExternalID: 1111,
	}

	// infer primary key
	field, value := InferPrimaryKey(record, true)
	assert.Equal(t, "partner_id", field)
	assert.Equal(t, 1111, value)
}

func TestInferPrimaryKey_notFound(t *testing.T) {
	record := struct {
		ExternalID int
		Name       string
	}{}

	assert.Panics(t, func() {
		InferPrimaryKey(record, true)
	})
}
