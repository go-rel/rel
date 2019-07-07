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

func TestAssociation(t *testing.T) {
	var (
		column       string
		usert        = reflect.TypeOf(User{})
		transactiont = reflect.TypeOf(Transaction{})
		addresst     = reflect.TypeOf(Address{})
	)

	_, cached := associationFieldCache.Load(associationFieldKey{rt: usert, field: "Transactions"})
	assert.False(t, cached)

	_, _, column = InferAssociationField(usert, "Transactions")
	assert.Equal(t, "user_id", column)

	_, cached = associationFieldCache.Load(associationFieldKey{rt: usert, field: "Transactions"})
	assert.True(t, true)

	// with cache
	_, _, column = InferAssociationField(usert, "Transactions")
	assert.Equal(t, "user_id", column)

	_, _, column = InferAssociationField(transactiont, "Buyer")
	assert.Equal(t, "id", column)

	// without struct tags
	_, _, column = InferAssociationField(addresst, "User")
	assert.Equal(t, "id", column)

	_, _, column = InferAssociationField(usert, "Addresses")
	assert.Equal(t, "user_id", column)
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
