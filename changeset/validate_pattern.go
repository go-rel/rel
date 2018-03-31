package changeset

import (
	"regexp"
	"strings"
)

var ValidatePatternErrorMessage = "{field}'s format is invalid"

func ValidatePattern(ch *Changeset, field string, pattern string, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		message: ValidatePatternErrorMessage,
	}
	options.Apply(opts)

	if str, ok := val.(string); ok {
		match, _ := regexp.MatchString(pattern, str)
		if !match {
			msg := strings.Replace(options.message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
