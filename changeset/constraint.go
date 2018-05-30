// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"strings"

	"github.com/Fs02/grimoire/errors"
)

// Constraint defines information to infer constraint error.
type Constraint struct {
	Field   string
	Message string
	Code    int
	Name    string
	Exact   bool
	Kind    errors.Kind
}

// Constraints is slice of Constraint
type Constraints []Constraint

// GetError converts error based on constraints.
// If the original error is constraint error, and it's defined in the constraint list, then it'll be updated with constraint's message.
// If the original error is constraint error but not defined in the constraint list, it'll be converted to unexpected error.
// else it'll not modify the error.
func (constraints Constraints) GetError(err errors.Error) error {
	if err.Kind() == errors.Unexpected || err.Kind() == errors.Changeset || err.Kind() == errors.NotFound {
		return err
	}

	for _, c := range constraints {
		if c.Kind == err.Kind() {
			if c.Exact && c.Name != err.Field {
				continue
			}

			if !c.Exact && !strings.Contains(err.Field, c.Name) {
				continue
			}

			return errors.NewWithCode(c.Message, c.Field, c.Code, c.Kind)
		}
	}

	return errors.NewUnexpected(err.Message)
}
