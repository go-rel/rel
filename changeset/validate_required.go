package changeset

import (
	"strings"
)

var ValidateRequiredErrorMessage = "{field} is required"

func ValidateRequired(ch *Changeset, fields ...string) {
	for _, f := range fields {
		_, exist := ch.changes[f]
		if !exist {
			msg := strings.Replace(ValidateRequiredErrorMessage, "{field}", f, 1)
			AddError(ch, f, msg)
		}
	}
}
