package rel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentMeta_Association(t *testing.T) {
	var (
		docMeta   = getDocumentMeta(reflect.TypeOf(User{}), false)
		assocMeta = docMeta.Association("address")
	)

	assert.Equal(t, getDocumentMeta(reflect.TypeOf(Address{}), false), assocMeta.DocumentMeta())
}

func TestDocumentMeta_Association_ptr(t *testing.T) {
	var (
		docMeta   = getDocumentMeta(reflect.TypeOf(User{}), false)
		assocMeta = docMeta.Association("work_address")
	)

	assert.Equal(t, getDocumentMeta(reflect.TypeOf(Address{}), false), assocMeta.DocumentMeta())
}

func TestDocumentMeta_Association_slice(t *testing.T) {
	var (
		docMeta   = getDocumentMeta(reflect.TypeOf(User{}), false)
		assocMeta = docMeta.Association("emails")
	)

	assert.Equal(t, getDocumentMeta(reflect.TypeOf(Email{}), false), assocMeta.DocumentMeta())
}

func TestDocumentMeta_Association_notFound(t *testing.T) {
	var (
		docMeta = getDocumentMeta(reflect.TypeOf(User{}), false)
	)

	assert.Panics(t, func() {
		docMeta.Association("invalid")
	})
}
