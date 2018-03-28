package changeset

import (
	"fmt"
	"strings"
)

var ValidateExclusionErrorMessage = "{field} must not be any of {values}"

func ValidateExclusion(ch *Changeset, field string, values []interface{}, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		Message: ValidateExclusionErrorMessage,
	}
	options.Apply(opts)

	invalid := false
	for _, inval := range values {
		if val == inval {
			invalid = true
			break
		}
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{values}", fmt.Sprintf("%v", values))
		AddError(ch, field, r.Replace(options.Message))
	}
}
