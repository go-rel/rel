// Package changeset allow filtering, casting and validation when manipulating structs.
package changeset

import (
	"reflect"
)

// Changeset allow filtering, casting and validation when manipulating structs.
type Changeset struct {
	errors  []error
	entity  interface{}
	params  map[string]interface{}
	changes map[string]interface{}
	values  map[string]interface{}
	types   map[string]reflect.Type
}

// Errors of changeset.
func (changeset *Changeset) Errors() []error {
	return changeset.errors
}

// Error of changeset, returns the first error if any.
func (changeset *Changeset) Error() error {
	if changeset.errors != nil {
		return changeset.errors[0]
	}
	return nil
}

// Changes of changeset.
func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}

// Values of changeset.
func (changeset *Changeset) Values() map[string]interface{} {
	return changeset.values
}

// Types of changeset.
func (changeset *Changeset) Types() map[string]reflect.Type {
	return changeset.types
}
