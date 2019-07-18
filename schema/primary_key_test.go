package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferPrimaryKey(t *testing.T) {
	var (
		record = struct {
			ID   uint
			Name string
		}{
			ID: 1,
		}
		rt            = reflectInternalType(record)
		expectedField = "id"
		expectedValue = uint(1)
	)

	// should not be cached yet
	_, cached := primaryKeysCache.Load(rt)
	assert.False(t, cached)

	// infer primary key
	field, values := InferPrimaryKey(record, true)
	assert.Equal(t, expectedField, field)
	assert.Equal(t, expectedValue, values[0])

	// cached
	_, cached = primaryKeysCache.Load(rt)
	assert.True(t, cached)

	record.ID = 2

	// infer primary key using cache
	field, values = InferPrimaryKey(record, true)
	assert.Equal(t, expectedField, field)
	assert.Equal(t, uint(2), values[0])
}

func TestInferPrimaryKey_usingInterface(t *testing.T) {
	var (
		record = CustomSchema{
			UUID: "abc123",
		}
		rt             = reflectInternalType(record)
		expectedField  = "_uuid"
		expectedValues = "abc123"
	)

	// should not be cached yet
	_, cached := primaryKeysCache.Load(rt)
	assert.False(t, cached)

	// infer primary key
	field, values := InferPrimaryKey(record, true)
	assert.Equal(t, expectedField, field)
	assert.Equal(t, expectedValues, values[0])

	// never cache
	_, cached = primaryKeysCache.Load(rt)
	assert.False(t, cached)
}

func TestInferPrimaryKey_usingTag(t *testing.T) {
	var (
		record = struct {
			ID         uint
			ExternalID int `db:",primary"`
			Name       string
		}{
			ExternalID: 12345,
		}
	)

	// infer primary key
	field, values := InferPrimaryKey(record, true)
	assert.Equal(t, "external_id", field)
	assert.Equal(t, 12345, values[0])
}

func TestInferPrimaryKey_usingTagAmdCustomName(t *testing.T) {
	var (
		record = struct {
			ID         uint
			ExternalID int `db:"partner_id,primary"`
			Name       string
		}{
			ExternalID: 1111,
		}
	)

	// infer primary key
	field, values := InferPrimaryKey(record, true)
	assert.Equal(t, "partner_id", field)
	assert.Equal(t, 1111, values[0])
}

func TestInferPrimaryKey_notFound(t *testing.T) {
	var (
		record = struct {
			ExternalID int
			Name       string
		}{}
	)

	assert.Panics(t, func() {
		InferPrimaryKey(record, true)
	})
}

func TestInferPrimaryKey_slice(t *testing.T) {
	var (
		records = []struct {
			ID   int
			Name string
		}{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		}
		expectedField  = "id"
		expectedValues = []interface{}{1, 2}
	)

	// infer primary keys
	field, values := InferPrimaryKey(records, true)
	assert.Equal(t, expectedField, field)
	assert.Equal(t, expectedValues, values)
}
