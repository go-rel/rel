package changeset

// Change make a new changeset without changes and build from given schema. Returns new Changeset.
func Change(schema interface{}) *Changeset {
	ch := &Changeset{}
	ch.changes = make(map[string]interface{})
	ch.values, ch.types = mapSchema(schema)

	return ch
}
