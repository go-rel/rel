package changeset

func AddError(ch *Changeset, field string, message string) {
	ch.errors = append(ch.errors, Error{
		Field:   field,
		Message: message,
	})
}
