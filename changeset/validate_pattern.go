package changeset

import (
	"regexp"
	"strings"
)

var ValidatePatternErrorMessage = "{field} is invalid"

func ValidatePattern(ch *Changeset, field string, pattern string) {
	val, exist := ch.changes[field]
	if !exist {
		return
	}

	if str, ok := val.(string); ok {
		match, _ := regexp.MatchString(pattern, str)
		if !match {
			msg := strings.Replace(ValidatePatternErrorMessage, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
