package changeset

import (
	"regexp"
	"strings"
)

var ValidateRegexpErrorMessage = "{field}'s format is invalid"

func ValidateRegexp(ch *Changeset, field string, exp *regexp.Regexp, opts ...Option) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	options := Options{
		Message: ValidateRegexpErrorMessage,
	}
	options.Apply(opts)

	if str, ok := val.(string); ok {
		match := exp.MatchString(str)
		if !match {
			msg := strings.Replace(options.Message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
