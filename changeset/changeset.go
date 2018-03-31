package changeset

import (
	"reflect"
)

type Changeset struct {
	errors  []error
	entity  interface{}
	params  map[string]interface{}
	changes map[string]interface{}
	values  map[string]interface{}
	types   map[string]reflect.Type
}

func (changeset *Changeset) Errors() []error {
	return changeset.errors
}

func (changeset *Changeset) Error() error {
	if changeset.errors != nil {
		return changeset.errors[0]
	}
	return nil
}

func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}

func (changeset *Changeset) Values() map[string]interface{} {
	return changeset.values
}

func (changeset *Changeset) Types() map[string]reflect.Type {
	return changeset.types
}
