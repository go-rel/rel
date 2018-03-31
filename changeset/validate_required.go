package changeset

import (
	"strings"
)

// ValidateRequiredErrorMessage is the default error message for ValidateRequired.
var ValidateRequiredErrorMessage = "{field} is required"

// ValidateRequired validates that one or more fields are present in the changeset.
func ValidateRequired(ch *Changeset, fields []string, opts ...Option) {
	options := Options{
		message: ValidateRequiredErrorMessage,
	}
	options.apply(opts)

	for _, f := range fields {
		_, exist := ch.changes[f]
		if !exist {
			msg := strings.Replace(options.message, "{field}", f, 1)
			AddError(ch, f, msg)
		}
	}
}
