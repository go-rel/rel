package grimoire

import (
	"strings"
)

type ConstraintType int8

const (
	CheckConstraint = iota
	UniqueConstraint
	ForeignKeyConstraint
)

func (ct ConstraintType) String() string {
	switch ct {
	case CheckConstraint:
		return "CheckConstraint"
	case UniqueConstraint:
		return "UniqueConstraint"
	case ForeignKeyConstraint:
		return "ForeignKeyConstraint"
	default:
		return ""
	}
}

type constraint struct {
	typ     ConstraintType
	key     string
	exact   bool
	field   string
	message string
}

func (c constraint) Build(changes *Changes) {
	changes.constraints = append(changes.constraints, c)
}

func Constraint(typ ConstraintType, key string, exact bool, field string, message string) Changer {
	return constraint{
		typ:     typ,
		key:     key,
		exact:   exact,
		field:   field,
		message: message,
	}
}

type constraints []constraint

func (cs constraints) transform(err error) error {
	cerr, ok := err.(ConstraintError)
	if !ok {
		return err
	}

	for _, c := range cs {
		if c.typ == cerr.Type {
			if c.exact && c.key != cerr.Key {
				continue
			}

			if !c.exact && !strings.Contains(cerr.Key, c.key) {
				continue
			}

			return ValidationError{Field: c.field, Message: c.message}
		}
	}

	return err
}
