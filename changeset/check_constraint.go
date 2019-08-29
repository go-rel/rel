package changeset

import (
	"strings"

	"github.com/Fs02/grimoire"
)

// CheckConstraintMessage is the default error message for CheckConstraint.
var CheckConstraintMessage = "{field} is invalid"

// CheckConstraint adds an unique constraint to changeset.
func CheckConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: CheckConstraintMessage,
		key:     field,
		exact:   false,
	}
	options.apply(opts)

	var (
		constraint = grimoire.Constraint(
			grimoire.CheckConstraint,
			options.key,
			options.exact,
			field,
			strings.Replace(options.message, "{field}", field, 1), // todo: defer evaluation
		)
	)

	ch.changers = append(ch.changers, constraint)
}
