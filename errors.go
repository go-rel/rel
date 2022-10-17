package rel

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound returned when entities not found.
	ErrNotFound = NotFoundError{}

	// ErrCheckConstraint is an auxiliary variable for error handling.
	// This is only to be used when checking error with errors.Is(err, ErrCheckConstraint).
	ErrCheckConstraint = ConstraintError{Type: CheckConstraint}

	// ErrNotNullConstraint is an auxiliary variable for error handling.
	// This is only to be used when checking error with errors.Is(err, ErrNotNullConstraint).
	ErrNotNullConstraint = ConstraintError{Type: NotNullConstraint}

	// ErrUniqueConstraint is an auxiliary variable for error handling.
	// This is only to be used when checking error with errors.Is(err, ErrUniqueConstraint).
	ErrUniqueConstraint = ConstraintError{Type: UniqueConstraint}

	// ErrPrimaryKeyConstraint is an auxiliary variable for error handling.
	// This is only to be used when checking error with errors.Is(err, ErrPrimaryKeyConstraint).
	ErrPrimaryKeyConstraint = ConstraintError{Type: PrimaryKeyConstraint}

	// ErrForeignKeyConstraint is an auxiliary variable for error handling.
	// This is only to be used when checking error with errors.Is(err, ErrForeignKeyConstraint).
	ErrForeignKeyConstraint = ConstraintError{Type: ForeignKeyConstraint}
)

// NotFoundError returned whenever Find returns no result.
type NotFoundError struct{}

// Error message.
func (nfe NotFoundError) Error() string {
	return "entity not found"
}

// Is returns true when target error is sql.ErrNoRows.
func (nfe NotFoundError) Is(target error) bool {
	return errors.Is(target, sql.ErrNoRows)
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

// Is returns true when target error have the same type and key if defined.
func (ce ConstraintError) Is(target error) bool {
	if err, ok := target.(ConstraintError); ok {
		return ce.Type == err.Type && (ce.Key == "" || err.Key == "" || ce.Key == err.Key)
	}

	return false
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
