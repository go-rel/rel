package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID           int
	Transactions []Transaction `references:"ID" foreign_key:"BuyerID"`
	Addresses    []Address
}

type Transaction struct {
	ID      int
	BuyerID int  `db:"user_id"`
	Buyer   User `references:"BuyerID" foreign_key:"ID"`
}

type Address struct {
	ID     int
	UserID *int
	User   *User
}

func TestAssociationField(t *testing.T) {
	var (
		usert        = reflect.TypeOf(User{})
		transactiont = reflect.TypeOf(Transaction{})
		addresst     = reflect.TypeOf(Address{})
	)

	_, cached := associationFieldCache.Load(associationFieldKey{rt: usert, field: "Transactions"})
	assert.False(t, cached)

	assoc := InferAssociationField(usert, "Transactions")
	assert.Equal(t, "user_id", assoc.ForeignColumn)

	_, cached = associationFieldCache.Load(associationFieldKey{rt: usert, field: "Transactions"})
	assert.True(t, true)

	// with cache
	assoc = InferAssociationField(usert, "Transactions")
	assert.Equal(t, "user_id", assoc.ForeignColumn)

	assoc = InferAssociationField(transactiont, "Buyer")
	assert.Equal(t, "id", assoc.ForeignColumn)

	// without struct tags
	assoc = InferAssociationField(addresst, "User")
	assert.Equal(t, "id", assoc.ForeignColumn)

	assoc = InferAssociationField(usert, "Addresses")
	assert.Equal(t, "user_id", assoc.ForeignColumn)
}

func TestAssociation_fieldNotFound(t *testing.T) {
	assert.Panics(t, func() {
		InferAssociationField(reflect.TypeOf(User{}), "Unknown")
	})
}

func TestAssociation_refFieldNotFound(t *testing.T) {
	type Invoice struct {
		User User `references:"UserID" foreign_key:"ID"`
	}

	assert.Panics(t, func() {
		InferAssociationField(reflect.TypeOf(Invoice{}), "User")
	})
}

func TestAssociation_fkFieldNotFound(t *testing.T) {
	type Invoice struct {
		UserID int
		User   User `references:"UserID" foreign_key:"UnknowFieldID"`
	}

	assert.Panics(t, func() {
		InferAssociationField(reflect.TypeOf(Invoice{}), "User")
	})
}
