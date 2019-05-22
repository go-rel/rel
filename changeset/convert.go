package changeset

// Convert a struct as changeset, every field's value will be treated as changes. Returns a new changeset.
func Convert(data interface{}) *Changeset {
	ch := &Changeset{}
	ch.values = make(map[string]interface{})
	ch.changes, ch.types, _ = mapSchema(data, false)

	return ch
}
