package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferValues(t *testing.T) {
	var (
		address = "address"
		record  = struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}{
			ID:      1,
			Name:    "name",
			Number:  10.5,
			Address: &address,
			Data:    []byte("data"),
		}
		expected = []interface{}{1, "name", 10.5, address, []byte("data")}
	)

	assert.Equal(t, expected, InferValues(record))
	assert.Equal(t, expected, InferValues(&record))
}

func TestInferValues_usingInterface(t *testing.T) {
	var (
		record = CustomSchema{
			UUID:  "abc123",
			Price: 100,
		}
		expected = []interface{}{"abc123", 100}
	)

	assert.Equal(t, expected, InferValues(record))
	assert.Equal(t, expected, InferValues(&record))
}
