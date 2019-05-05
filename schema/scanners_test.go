package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferScanners(t *testing.T) {
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
		expected = []interface{}{
			Indirect(&record.ID),
			Indirect(&record.Name),
			Indirect(&record.Number),
			&record.Address,
			Indirect(&record.Data),
		}
	)

	assert.Equal(t, expected, InferScanners(&record))
}

func TestInferScanners_usingInterface(t *testing.T) {
	var (
		record = CustomSchema{
			UUID:  "abc123",
			Price: 100,
		}
		expected = []interface{}{Indirect(&record.UUID), Indirect(&record.Price)}
	)

	assert.Equal(t, expected, InferScanners(&record))
}

func TestInferScanners_notPointer(t *testing.T) {
	assert.Panics(t, func() {
		InferScanners(struct{}{})
	})
}
