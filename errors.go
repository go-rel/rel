package rel

type NoResultError struct{}

func (nre NoResultError) Error() string {
	return "No result found"
}

type ConstraintType int8

const (
	CheckConstraint ConstraintType = iota
	NotNullConstraint
	UniqueConstraint
	PrimaryKeyConstraint
	ForeignKeyConstraint
)

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

type ConstraintError struct {
	Key  string
	Type ConstraintType
	Err  error
}

func (ce ConstraintError) Unwrap() error {
	return ce.Err
}

func (ce ConstraintError) Error() string {
	if ce.Err != nil {
		return ce.Type.String() + "Error: " + ce.Err.Error()
	}

	return ce.Type.String() + "Error"
}
