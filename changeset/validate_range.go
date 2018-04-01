package changeset

import (
	"strconv"
	"strings"
)

// ValidateRangeErrorMessage is the default error message for ValidateRange.
var ValidateRangeErrorMessage = "{field} must be between {min} and {max}"

// ValidateRange validates the value of given field is not larger than max and not smaller than min.
// Validation can be performed against string, slice and numbers.
func ValidateRange(ch *Changeset, field string, min int, max int, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidateRangeErrorMessage,
	}
	options.apply(opts)

	invalid := false

	switch v := val.(type) {
	case string:
		invalid = len(v) < min || len(v) > max
	case []interface{}:
		invalid = len(v) < min || len(v) > max
	case int:
		invalid = v < min || v > max
	case int8:
		invalid = v < int8(min) || v > int8(max)
	case int16:
		invalid = v < int16(min) || v > int16(max)
	case int32:
		invalid = v < int32(min) || v > int32(max)
	case int64:
		invalid = v < int64(min) || v > int64(max)
	case uint:
		invalid = v < uint(min) || v > uint(max)
	case uint8:
		invalid = v < uint8(min) || v > uint8(max)
	case uint16:
		invalid = v < uint16(min) || v > uint16(max)
	case uint32:
		invalid = v < uint32(min) || v > uint32(max)
	case uint64:
		invalid = v < uint64(min) || v > uint64(max)
	case uintptr:
		invalid = v < uintptr(min) || v > uintptr(max)
	case float32:
		invalid = v < float32(min) || v > float32(max)
	case float64:
		invalid = v < float64(min) || v > float64(max)
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{min}", strconv.Itoa(min), "{max}", strconv.Itoa(max))
		AddError(ch, field, r.Replace(options.message))
	}
}
