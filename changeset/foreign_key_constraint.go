package changeset

import (
	"strings"

	"github.com/Fs02/grimoire/errors"
)

// ForeignKeyConstraintMessage is the default error message for ForeignKeyConstraint.
var ForeignKeyConstraintMessage = "does not exist"

// ForeignKeyConstraint adds an unique constraint to changeset.
func ForeignKeyConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: ForeignKeyConstraintMessage,
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
		Kind:    errors.ForeignKeyConstraint,
	})
}
