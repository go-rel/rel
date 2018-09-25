package changeset

import (
	"strings"
)

// ValidateRequiredErrorMessage is the default error message for ValidateRequired.
var ValidateRequiredErrorMessage = "{field} is required"

// ValidateRequired validates that one or more fields are present in the changeset.
// It'll add error to changeset if field in the changes is nil or string made only of whitespace,
func ValidateRequired(ch *Changeset, fields []string, opts ...Option) {
	options := Options{
		message: ValidateRequiredErrorMessage,
	}
	options.apply(opts)

	for _, f := range fields {
		val, exist := ch.changes[f]

		// check values if it's not exist in changeset when changeOnly is false
		if !exist && !options.changeOnly {
			val, exist = ch.values[f]
		}

		str, isStr := val.(string)
		if exist && (isStr && strings.TrimSpace(str) != "") || (!isStr && val != nil) {
			continue
		}

		msg := strings.Replace(options.message, "{field}", f, 1)
		AddError(ch, f, msg)
	}
}
