package changeset

import (
	"github.com/Fs02/grimoire/errors"
)

type Changeset struct {
	errors  errors.Errors
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
