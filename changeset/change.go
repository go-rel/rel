package changeset

// Change make a new changeset without changes and build from given schema. Returns new Changeset.
func Change(schema interface{}, changes ...map[string]interface{}) *Changeset {
	ch := &Changeset{}
	ch.changes = make(map[string]interface{})
	ch.values, ch.types, _ = mapSchema(schema, false)

	if len(changes) > 0 {
		ch.changes = changes[0]
	} else {
		ch.changes = make(map[string]interface{})
	}

	return ch
}
