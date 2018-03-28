package changeset

import (
	"fmt"
	"strings"
)

var ValidateInclusionErrorMessage = "{field} must be one of {values}"

func ValidateInclusion(ch *Changeset, field string, values []interface{}, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		Message: ValidateInclusionErrorMessage,
	}
	options.Apply(opts)

	invalid := true
	for _, inval := range values {
		if val == inval {
			invalid = false
			break
		}
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{values}", fmt.Sprintf("%v", values))
		AddError(ch, field, r.Replace(options.Message))
	}
}
