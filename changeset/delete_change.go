package changeset

// DeleteChange from changeset.
func DeleteChange(ch *Changeset, field string) {
	delete(ch.changes, field)
}
