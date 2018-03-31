package changeset

import (
	"strings"
)

var ValidateRequiredErrorMessage = "{field} is required"

func ValidateRequired(ch *Changeset, fields []string, opts ...Option) {
	options := Options{
		message: ValidateRequiredErrorMessage,
	}
	options.Apply(opts)

	for _, f := range fields {
		_, exist := ch.changes[f]
		if !exist {
			msg := strings.Replace(options.message, "{field}", f, 1)
			AddError(ch, f, msg)
		}
	}
}
