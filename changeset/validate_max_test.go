package changeset

import (
	"fmt"
	"testing"
)

func TestValidateRange(t *testing.T) {
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

			ValidateRange(ch, "field", 5, 15)

			if ch.Errors() != nil {
				t.Error(`Expected nil but got`, ch.Errors())
			}
		})
	}
}

func TestValidateRangeError(t *testing.T) {
	tests := []interface{}{
		"long text",
		[]interface{}{"a", "b", "c", "d", "e", "f"},
		10,
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

			ValidateRange(ch, "field", 15, 20)

			if ch.Errors().Error() != "field must be between 15 and 20" {
				t.Error(`Expected "field must be between 15 and 20" but got`, ch.Errors().Error())
			}
		})
	}
}

func TestValidateRangeMissing(t *testing.T) {
	ch := &Changeset{}
	ValidateRange(ch, "field", 5, 15)

	if ch.Errors() != nil {
		t.Error(`Expected nil but got`, ch.Errors())
	}
}
