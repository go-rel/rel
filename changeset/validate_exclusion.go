package changeset

import (
	"fmt"
	"strings"
)

// ValidateExclusionErrorMessage is the default error message for ValidateExclusion.
var ValidateExclusionErrorMessage = "{field} must not be any of {values}"

// ValidateExclusion validates a change is not included in the given values.
func ValidateExclusion(ch *Changeset, field string, values []interface{}, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidateExclusionErrorMessage,
	}
	options.apply(opts)

	invalid := false
	for _, inval := range values {
		if val == inval {
			invalid = true
			break
		}
	}

	if invalid {
		r := strings.NewReplacer("{field}", field, "{values}", fmt.Sprintf("%v", values))
		AddError(ch, field, r.Replace(options.message))
	}
}
