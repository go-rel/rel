package changeset

import (
	"strings"

	"github.com/Fs02/grimoire"
)

// ForeignKeyConstraintMessage is the default error message for ForeignKeyConstraint.
var ForeignKeyConstraintMessage = "does not exist"

// ForeignKeyConstraint adds an unique constraint to changeset.
func ForeignKeyConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: ForeignKeyConstraintMessage,
		key:     field,
		exact:   false,
	}
	options.apply(opts)

	var (
		constraint = grimoire.Constraint(
			grimoire.ForeignKeyConstraint,
			options.key,
			options.exact,
			field,
			strings.Replace(options.message, "{field}", field, 1), // todo: defer evaluation
		)
	)

	ch.changers = append(ch.changers, constraint)
}
