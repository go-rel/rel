package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID           int
	Address      Address
	Transactions []Transaction `references:"ID" foreign_key:"BuyerID"`
	Addresses    []Address     // TODO: remmove
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
		transaction       = &Transaction{ID: 1}
		user              = &User{ID: 2}
		address           = &Address{ID: 3}
		transactionLoaded = &Transaction{ID: 1, BuyerID: user.ID, Buyer: *user}
		userLoaded        = &User{ID: 2, Address: *address, Transactions: []Transaction{*transaction}}
		addressLoaded     = &Address{ID: 3, UserID: &user.ID, User: user}
	)

	tests := []struct {
		entity          string
		field           string
		data            interface{}
		typ             AssociationType
		target          interface{}
		loaded          bool
		referenceColumn string
		referenceValue  interface{}
		foreignColumn   string
		foreignValue    interface{}
	}{
		{
			entity:          "Transaction",
			field:           "Buyer",
			data:            transaction,
			typ:             BelongsTo,
			target:          &transaction.Buyer,
			loaded:          false,
			referenceColumn: "user_id",
			referenceValue:  transaction.BuyerID,
			foreignColumn:   "id",
			foreignValue:    transaction.Buyer.ID,
		},
		{
			entity:          "Transaction",
			field:           "Buyer",
			data:            transactionLoaded,
			typ:             BelongsTo,
			target:          &transactionLoaded.Buyer,
			loaded:          true,
			referenceColumn: "user_id",
			referenceValue:  transactionLoaded.BuyerID,
			foreignColumn:   "id",
			foreignValue:    transactionLoaded.Buyer.ID,
		},
		{
			entity:          "User",
			field:           "Transactions",
			data:            user,
			typ:             HasMany,
			target:          &user.Transactions,
			loaded:          false,
			referenceColumn: "id",
			referenceValue:  user.ID,
			foreignColumn:   "user_id",
			foreignValue:    nil,
		},
		{
			entity:          "User",
			field:           "Transactions",
			data:            userLoaded,
			typ:             HasMany,
			target:          &userLoaded.Transactions,
			loaded:          true,
			referenceColumn: "id",
			referenceValue:  userLoaded.ID,
			foreignColumn:   "user_id",
			foreignValue:    nil,
		},
		{
			entity:          "User",
			field:           "Address",
			data:            user,
			typ:             HasOne,
			target:          &user.Address,
			loaded:          false,
			referenceColumn: "id",
			referenceValue:  user.ID,
			foreignColumn:   "user_id",
			foreignValue:    user.Address.UserID,
		},
		{
			entity:          "User",
			field:           "Address",
			data:            userLoaded,
			typ:             HasOne,
			target:          &userLoaded.Address,
			loaded:          true,
			referenceColumn: "id",
			referenceValue:  userLoaded.ID,
			foreignColumn:   "user_id",
			foreignValue:    userLoaded.Address.UserID,
		},
		{
			entity:          "Address",
			field:           "User",
			data:            address,
			typ:             BelongsTo,
			target:          &User{}, // should be initialized to zero struct
			loaded:          false,
			referenceColumn: "user_id",
			referenceValue:  address.UserID,
			foreignColumn:   "id",
			foreignValue:    0,
		},
		{
			entity:          "Address",
			field:           "User",
			data:            addressLoaded,
			typ:             BelongsTo,
			target:          addressLoaded.User,
			loaded:          true,
			referenceColumn: "user_id",
			referenceValue:  addressLoaded.UserID,
			foreignColumn:   "id",
			foreignValue:    addressLoaded.User.ID,
		},
	}

	for _, test := range tests {
		t.Run(test.entity+"."+test.field, func(t *testing.T) {
			var (
				rv             = reflect.ValueOf(test.data)
				assoc          = InferAssociation(rv, test.field)
				target, loaded = assoc.TargetAddr()
			)

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.target, target)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.referenceColumn, assoc.ReferenceColumn())
			assert.Equal(t, test.referenceValue, assoc.ReferenceValue())
			assert.Equal(t, test.foreignColumn, assoc.ForeignColumn())

			if test.typ == HasMany {
				assert.Panics(t, func() {
					assert.Equal(t, test.foreignValue, assoc.ForeignValue())
				})
			} else {
				assert.NotPanics(t, func() {
					assert.Equal(t, test.foreignValue, assoc.ForeignValue())
				})
			}
		})
	}
}
