package changeset

import (
	"regexp"
	"strings"
)

var ValidateRegexpErrorMessage = "{field} is invalid"

func ValidateRegexp(ch *Changeset, field string, exp *regexp.Regexp) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	if str, ok := val.(string); ok {
		match := exp.MatchString(str)
		if !match {
			msg := strings.Replace(ValidateRegexpErrorMessage, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
