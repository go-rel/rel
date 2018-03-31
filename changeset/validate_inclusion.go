package changeset

import (
	"fmt"
	"strings"
)

// ValidateInclusionErrorMessage is the default error message for ValidateInclusion.
var ValidateInclusionErrorMessage = "{field} must be one of {values}"

// ValidateInclusion validates a change is included in the given values.
func ValidateInclusion(ch *Changeset, field string, values []interface{}, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidateInclusionErrorMessage,
	}
	options.apply(opts)

	invalid := true
	for _, inval := range values {
		if val == inval {
			invalid = false
			break
		}
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{values}", fmt.Sprintf("%v", values))
		AddError(ch, field, r.Replace(options.message))
	}
}
