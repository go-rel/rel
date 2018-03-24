package changeset

type Changeset struct {
	errors  Errors
	changes map[string]interface{}
}

func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}

func (changeset *Changeset) Errors() error {
	if len(changeset.errors) > 0 {
		return changeset.errors
	} else {
		return nil
	}
}
