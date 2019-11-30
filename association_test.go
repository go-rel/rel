package rel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssociation_Document(t *testing.T) {
	var (
		transaction       = &Transaction{ID: 1}
		user              = &User{ID: 2}
		address           = &Address{ID: 3}
		transactionLoaded = &Transaction{ID: 1, BuyerID: user.ID, Buyer: *user}
		userLoaded        = &User{ID: 2, Address: *address, Transactions: []Transaction{*transaction}}
		addressLoaded     = &Address{ID: 3, UserID: &user.ID, User: user}
	)

	tests := []struct {
		record         string
		field          string
		data           interface{}
		typ            AssociationType
		doc            *Document
		loaded         bool
		referenceField string
		referenceValue interface{}
		foreignField   string
		foreignValue   interface{}
	}{
		{
			record:         "Transaction",
			field:          "Buyer",
			data:           transaction,
			typ:            BelongsTo,
			doc:            NewDocument(&transaction.Buyer),
			loaded:         false,
			referenceField: "user_id",
			referenceValue: transaction.BuyerID,
			foreignField:   "id",
			foreignValue:   transaction.Buyer.ID,
		},
		{
			record:         "Transaction",
			field:          "Buyer",
			data:           transactionLoaded,
			typ:            BelongsTo,
			doc:            NewDocument(&transactionLoaded.Buyer),
			loaded:         true,
			referenceField: "user_id",
			referenceValue: transactionLoaded.BuyerID,
			foreignField:   "id",
			foreignValue:   transactionLoaded.Buyer.ID,
		},
		{
			record:         "User",
			field:          "Address",
			data:           user,
			typ:            HasOne,
			doc:            NewDocument(&user.Address),
			loaded:         false,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			record:         "User",
			field:          "Address",
			data:           userLoaded,
			typ:            HasOne,
			doc:            NewDocument(&userLoaded.Address),
			loaded:         true,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		// {
		// 	record:         "Address",
		// 	field:          "User",
		// 	data:           address,
		// 	typ:            BelongsTo,
		// 	doc:         NewDocument(&User{}), // should be initialized to zero struct
		// 	loaded:         false,
		// 	referenceField: "user_id",
		// 	referenceValue: address.UserID,
		// 	foreignField:   "id",
		// 	foreignValue:   0,
		// },
		{
			record:         "Address",
			field:          "User",
			data:           addressLoaded,
			typ:            BelongsTo,
			doc:            NewDocument(addressLoaded.User),
			loaded:         true,
			referenceField: "user_id",
			referenceValue: *addressLoaded.UserID,
			foreignField:   "id",
			foreignValue:   addressLoaded.User.ID,
		},
	}

	for _, test := range tests {
		t.Run(test.record+"."+test.field, func(t *testing.T) {
			var (
				rv          = reflect.ValueOf(test.data)
				sf, _       = rv.Type().Elem().FieldByName(test.field)
				assoc       = newAssociation(rv, sf.Index[0])
				doc, loaded = assoc.Document()
			)

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.doc, doc)
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

func TestAssociation_Collection(t *testing.T) {
	var (
		transaction = &Transaction{ID: 1}
		user        = &User{ID: 2}
		address     = &Address{ID: 3}
		userLoaded  = &User{ID: 2, Address: *address, Transactions: []Transaction{*transaction}}
	)

	tests := []struct {
		record         string
		field          string
		data           interface{}
		typ            AssociationType
		col            *Collection
		loaded         bool
		referenceField string
		referenceValue interface{}
		foreignField   string
		foreignValue   interface{}
	}{
		{
			record:         "User",
			field:          "Transactions",
			data:           user,
			typ:            HasMany,
			col:            NewCollection(&user.Transactions),
			loaded:         false,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			record:         "User",
			field:          "Transactions",
			data:           userLoaded,
			typ:            HasMany,
			col:            NewCollection(&userLoaded.Transactions),
			loaded:         true,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.record+"."+test.field, func(t *testing.T) {
			var (
				rv          = reflect.ValueOf(test.data)
				sf, _       = rv.Type().Elem().FieldByName(test.field)
				assoc       = newAssociation(rv, sf.Index[0])
				col, loaded = assoc.Collection()
			)

			test.col.reflect()

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.col, col)
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
