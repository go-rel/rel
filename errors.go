package rel

import (
	"fmt"
)

// NoResultError returned whenever Find returns no result.
type NoResultError struct{}

// Error message.
func (nre NoResultError) Error() string {
	return "No result found"
}

// ConstraintType defines the type of constraint error.
type ConstraintType int8

const (
	// CheckConstraint error type.
	CheckConstraint ConstraintType = iota
	// NotNullConstraint error type.1
	NotNullConstraint
	// UniqueConstraint error type.1
	UniqueConstraint
	// PrimaryKeyConstraint error type.1
	PrimaryKeyConstraint
	// ForeignKeyConstraint error type.1
	ForeignKeyConstraint
)

// String representation of the constraint type.
func (ct ConstraintType) String() string {
	switch ct {
	case CheckConstraint:
		return "CheckConstraint"
	case NotNullConstraint:
		return "NotNullConstraint"
	case UniqueConstraint:
		return "UniqueConstraint"
	case PrimaryKeyConstraint:
		return "PrimaryKeyConstraint"
	case ForeignKeyConstraint:
		return "ForeignKeyConstraint"
	default:
		return ""
	}
}

// ConstraintError returned whenever constraint error encountered.
type ConstraintError struct {
	Key  string
	Type ConstraintType
	Err  error
}

// Unwrap internal error returned by database driver.
func (ce ConstraintError) Unwrap() error {
	return ce.Err
}

// Error message.
func (ce ConstraintError) Error() string {
	if ce.Err != nil {
		return ce.Type.String() + "Error: " + ce.Err.Error()
	}

	return ce.Type.String() + "Error"
}

// ValueError returned whenever a field in a document is not exists or assigned using incorrect type.
type ValueError struct {
	Table string
	Field string
	Value interface{}
}

func (ve ValueError) Error() string {
	return fmt.Sprint("rel: cannot assign", ve.Value, "as", ve.Field, "into", ve.Table, "table.")
}
