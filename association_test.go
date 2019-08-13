package grimoire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		entity         string
		field          string
		data           interface{}
		typ            AssociationType
		target         interface{}
		loaded         bool
		referenceField string
		referenceValue interface{}
		foreignField   string
		foreignValue   interface{}
	}{
		{
			entity:         "Transaction",
			field:          "Buyer",
			data:           transaction,
			typ:            BelongsTo,
			target:         newDocument(&transaction.Buyer),
			loaded:         false,
			referenceField: "user_id",
			referenceValue: transaction.BuyerID,
			foreignField:   "id",
			foreignValue:   transaction.Buyer.ID,
		},
		{
			entity:         "Transaction",
			field:          "Buyer",
			data:           transactionLoaded,
			typ:            BelongsTo,
			target:         newDocument(&transactionLoaded.Buyer),
			loaded:         true,
			referenceField: "user_id",
			referenceValue: transactionLoaded.BuyerID,
			foreignField:   "id",
			foreignValue:   transactionLoaded.Buyer.ID,
		},
		{
			entity:         "User",
			field:          "Transactions",
			data:           user,
			typ:            HasMany,
			target:         newCollection(&user.Transactions),
			loaded:         false,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			entity:         "User",
			field:          "Transactions",
			data:           userLoaded,
			typ:            HasMany,
			target:         newCollection(&userLoaded.Transactions),
			loaded:         true,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			entity:         "User",
			field:          "Address",
			data:           user,
			typ:            HasOne,
			target:         newDocument(&user.Address),
			loaded:         false,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			entity:         "User",
			field:          "Address",
			data:           userLoaded,
			typ:            HasOne,
			target:         newDocument(&userLoaded.Address),
			loaded:         true,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		// {
		// 	entity:         "Address",
		// 	field:          "User",
		// 	data:           address,
		// 	typ:            BelongsTo,
		// 	target:         newDocument(&User{}), // should be initialized to zero struct
		// 	loaded:         false,
		// 	referenceField: "user_id",
		// 	referenceValue: address.UserID,
		// 	foreignField:   "id",
		// 	foreignValue:   0,
		// },
		{
			entity:         "Address",
			field:          "User",
			data:           addressLoaded,
			typ:            BelongsTo,
			target:         newDocument(addressLoaded.User),
			loaded:         true,
			referenceField: "user_id",
			referenceValue: *addressLoaded.UserID,
			foreignField:   "id",
			foreignValue:   addressLoaded.User.ID,
		},
	}

	for _, test := range tests {
		t.Run(test.entity+"."+test.field, func(t *testing.T) {
			var (
				rv             = reflect.ValueOf(test.data)
				assoc          = newAssociation(rv, test.field)
				target, loaded = assoc.Target()
			)

			switch v := test.target.(type) {
			case *document:
				v.reflect()
			case *collection:
				v.reflect()
			}

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.target, target)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.referenceField, assoc.ReferenceField())
			assert.Equal(t, test.referenceValue, assoc.ReferenceValue())
			assert.Equal(t, test.foreignField, assoc.ForeignField())

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
