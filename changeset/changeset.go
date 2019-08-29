// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"reflect"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/params"
)

// Changeset used to cast and validate data before saving it to the database.
// TODO: use arrays (make use of schema inferer)
// TODO: impement grimoire.Builder
type Changeset struct {
	errors   []error
	params   params.Params
	changes  map[string]interface{}
	values   map[string]interface{}
	types    map[string]reflect.Type
	changers []grimoire.Changer
	zero     bool
}

func (c *Changeset) Build(changes *grimoire.Changes) {

}

// Errors of changeset.
func (c *Changeset) Errors() []error {
	return c.errors
}

// Error of changeset, returns the first error if any.
func (c *Changeset) Error() error {
	if c.errors != nil {
		return c.errors[0]
	}
	return nil
}

// Get a change from changeset.
func (c *Changeset) Get(field string) interface{} {
	return c.changes[field]
}

// Fetch a change or value from changeset.
func (c *Changeset) Fetch(field string) interface{} {
	if change, ok := c.changes[field]; ok {
		return change
	}

	return c.values[field]
}
