package changeset

import (
	"strconv"
	"strings"
)

// ValidateMaxErrorMessage is the default error message for ValidateMax.
var ValidateMaxErrorMessage = "{field} must be less than {max}"

// ValidateMax validates the value of given field is not larger than max.
// Validation can be performed against string, slice and numbers.
func ValidateMax(ch *Changeset, field string, max int, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidateMaxErrorMessage,
	}
	options.apply(opts)

	invalid := false

	switch v := val.(type) {
	case string:
		invalid = len(v) > max
	case []interface{}:
		invalid = len(v) > max
	case []*Changeset:
		invalid = len(v) > max
	case int:
		invalid = v > max
	case int8:
		invalid = v > int8(max)
	case int16:
		invalid = v > int16(max)
	case int32:
		invalid = v > int32(max)
	case int64:
		invalid = v > int64(max)
	case uint:
		invalid = v > uint(max)
	case uint8:
		invalid = v > uint8(max)
	case uint16:
		invalid = v > uint16(max)
	case uint32:
		invalid = v > uint32(max)
	case uint64:
		invalid = v > uint64(max)
	case uintptr:
		invalid = v > uintptr(max)
	case float32:
		invalid = v > float32(max)
	case float64:
		invalid = v > float64(max)
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{max}", strconv.Itoa(max))
		AddError(ch, field, r.Replace(options.message))
	}
}
