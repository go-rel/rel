package changeset

// Change struct as changeset, every field's value will be treated as changes. Returns a new changeset.
func Change(entity interface{}) *Changeset {
	ch := &Changeset{}
	ch.entity = entity
	ch.values = make(map[string]interface{})
	ch.changes, ch.types = mapSchema(ch.entity)

	return ch
}
