package changeset

import (
	"strings"

	"github.com/Fs02/grimoire"
)

// UniqueConstraintMessage is the default error message for UniqueConstraint.
var UniqueConstraintMessage = "{field} has already been taken"

// UniqueConstraint adds an unique constraint to changeset.
func UniqueConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: UniqueConstraintMessage,
		key:     field,
		exact:   false,
	}
	options.apply(opts)

	var (
		constraint = grimoire.Constraint(
			grimoire.UniqueConstraint,
			options.key,
			options.exact,
			field,
			strings.Replace(options.message, "{field}", field, 1), // todo: defer evaluation
		)
	)

	ch.changers = append(ch.changers, constraint)
}
