package changeset

import (
	"strconv"
	"strings"
)

// ValidateMinErrorMessage is the default error message for ValidateMin.
var ValidateMinErrorMessage = "{field} must be more than {min}"

// ValidateMin validates the value of given field is not smaller than min.
// Validation can be performed against string, slice and numbers.
func ValidateMin(ch *Changeset, field string, min int, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidateMinErrorMessage,
	}
	options.apply(opts)

	invalid := false

	switch v := val.(type) {
	case string:
		invalid = len(v) < min
	case []interface{}:
		invalid = len(v) < min
	case []*Changeset:
		invalid = len(v) < min
	case int:
		invalid = v < min
	case int8:
		invalid = v < int8(min)
	case int16:
		invalid = v < int16(min)
	case int32:
		invalid = v < int32(min)
	case int64:
		invalid = v < int64(min)
	case uint:
		invalid = v < uint(min)
	case uint8:
		invalid = v < uint8(min)
	case uint16:
		invalid = v < uint16(min)
	case uint32:
		invalid = v < uint32(min)
	case uint64:
		invalid = v < uint64(min)
	case uintptr:
		invalid = v < uintptr(min)
	case float32:
		invalid = v < float32(min)
	case float64:
		invalid = v < float64(min)
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{min}", strconv.Itoa(min))
		AddError(ch, field, r.Replace(options.message))
	}
}
