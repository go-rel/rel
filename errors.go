package rel

// NotFoundError returned whenever Find returns no result.
type NotFoundError struct{}

// Error message.
func (nfe NotFoundError) Error() string {
	return "Record Not Found"
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
