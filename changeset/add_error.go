package changeset

import (
	"github.com/Fs02/grimoire/errors"
)

// AddError adds an error to changeset.
//	ch := changeset.Cast(user, params, fields)
//	changeset.AddError(ch, "field", "error")
//	ch.Errors() // []errors.Error{{Field: "field", Message: "error"}}
func AddError(ch *Changeset, field string, message string) {
	ch.errors = append(ch.errors, errors.New(message, field, errors.Changeset))
}
