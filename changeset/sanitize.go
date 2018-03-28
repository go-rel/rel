package changeset

import (
	"github.com/microcosm-cc/bluemonday"
)

var SanitizePolicy = bluemonday.UGCPolicy()

func Sanitize(ch *Changeset, fields []string) {
	for _, f := range fields {
		if val, exist := ch.changes[f]; exist {
			if str, ok := val.(string); ok {
				ch.changes[f] = SanitizePolicy.Sanitize(str)
			}
		}
	}
}
