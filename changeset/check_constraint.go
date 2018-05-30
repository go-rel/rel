package changeset

import (
	"strings"

	"github.com/Fs02/grimoire/errors"
)

// CheckConstraintMessage is the default error message for CheckConstraint.
var CheckConstraintMessage = "{field} is invalid"

// CheckConstraint adds an unique constraint to changeset.
func CheckConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: CheckConstraintMessage,
		name:    field,
		exact:   false,
	}
	options.apply(opts)

	ch.constraints = append(ch.constraints, Constraint{
		Field:   field,
		Message: strings.Replace(options.message, "{field}", field, 1),
		Code:    options.code,
		Name:    options.name,
		Exact:   options.exact,
		Kind:    errors.CheckConstraint,
	})
}
