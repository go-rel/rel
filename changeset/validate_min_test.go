package changeset

import (
	"fmt"
	"testing"
)

func TestValidateMin(t *testing.T) {
	tests := []interface{}{
		"long text",
		10,
		[]interface{}{"a", "b", "c", "d", "e", "f"},
		int8(10),
		int16(10),
		int32(10),
		int64(10),
		uint(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
		uintptr(10),
		float32(10),
		float64(10),
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateMin(ch, "field", 5)

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidateMinError(t *testing.T) {
	tests := []interface{}{
		"long text",
		10,
		[]interface{}{"a", "b", "c", "d", "e", "f"},
		int8(10),
		int16(10),
		int32(10),
		int64(10),
		uint(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
		uintptr(10),
		float32(10),
		float64(10),
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateMin(ch, "field", 15)

			if ch.Errors().Error() != "field must be more than 15" {
				t.Error(`Expected "field must be more than 15" but got`, ch.Errors().Error())
			}
		})
	}
}

func TestValidateMinMissing(t *testing.T) {
	ch := &Changeset{}
	ValidateMin(ch, "field", 5)

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
