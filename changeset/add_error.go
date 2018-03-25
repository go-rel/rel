package changeset

import (
	"github.com/Fs02/grimoire/errors"
)

func AddError(ch *Changeset, field string, message string) {
	ch.errors = append(ch.errors, errors.ChangesetError(message, field))
}
