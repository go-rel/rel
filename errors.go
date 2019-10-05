package rel

type NoResultError struct{}

func (nre NoResultError) Error() string {
	return "No result found"
}

type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return "ValidationError: " + ve.Message
}

type ConstraintError struct {
	Key  string
	Type ConstraintType
	Err  error
}

func (ce ConstraintError) Error() string {
	return ce.Type.String() + "Error: " + ce.Err.Error()
}
