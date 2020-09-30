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
		isZero         bool
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
			isZero:         true,
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
			isZero:         false,
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
			isZero:         true,
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
			isZero:         false,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			record:         "Address",
			field:          "User",
			data:           address,
			typ:            BelongsTo,
			doc:            NewDocument(&User{}), // should be initialized to zero struct
			loaded:         false,
			isZero:         true,
			referenceField: "user_id",
			referenceValue: nil,
			foreignField:   "id",
			foreignValue:   0,
		},
		{
			record:         "Address",
			field:          "User",
			data:           addressLoaded,
			typ:            BelongsTo,
			doc:            NewDocument(addressLoaded.User),
			loaded:         true,
			isZero:         false,
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
			assert.Equal(t, test.doc.rt, doc.rt)
			assert.Equal(t, test.doc.data, doc.data)
			assert.Equal(t, test.doc.v, doc.v)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.isZero, assoc.IsZero())
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
		role        = &Role{ID: 4}
		userLoaded  = &User{
			ID: 2, Address: *address,
			Transactions: []Transaction{*transaction},
		}
	)

	tests := []struct {
		record           string
		field            string
		data             interface{}
		typ              AssociationType
		col              *Collection
		loaded           bool
		isZero           bool
		referenceField   string
		referenceValue   interface{}
		referenceThrough string
		foreignField     string
		foreignValue     interface{}
		foreignThrough   string
		through          string
	}{
		{
			record:         "User",
			field:          "Transactions",
			data:           user,
			typ:            HasMany,
			col:            NewCollection(&user.Transactions),
			loaded:         false,
			isZero:         true,
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
			isZero:         false,
			referenceField: "id",
			referenceValue: userLoaded.ID,
			foreignField:   "user_id",
			foreignValue:   nil,
		},
		{
			record:           "User",
			field:            "Roles",
			data:             user,
			typ:              ManyToMany,
			col:              NewCollection(&user.Roles),
			loaded:           false,
			isZero:           true,
			referenceField:   "id",
			referenceThrough: "user_id",
			referenceValue:   user.ID,
			foreignField:     "id",
			foreignThrough:   "role_id",
			foreignValue:     nil,
			through:          "user_roles",
		},
		{
			record:           "Role",
			field:            "Users",
			data:             role,
			typ:              ManyToMany,
			col:              NewCollection(&role.Users),
			loaded:           false,
			isZero:           true,
			referenceField:   "id",
			referenceThrough: "role_id",
			referenceValue:   role.ID,
			foreignField:     "id",
			foreignThrough:   "user_id",
			foreignValue:     nil,
			through:          "user_roles",
		},
		{
			record:           "User",
			field:            "Followers",
			data:             user,
			typ:              ManyToMany,
			col:              NewCollection(&user.Followers),
			loaded:           false,
			isZero:           true,
			referenceField:   "id",
			referenceThrough: "following_id",
			referenceValue:   user.ID,
			foreignField:     "id",
			foreignThrough:   "follower_id",
			foreignValue:     nil,
			through:          "followers",
		},
		{
			record:           "User",
			field:            "Followings",
			data:             user,
			typ:              ManyToMany,
			col:              NewCollection(&user.Followings),
			loaded:           false,
			isZero:           true,
			referenceField:   "id",
			referenceThrough: "follower_id",
			referenceValue:   user.ID,
			foreignField:     "id",
			foreignThrough:   "following_id",
			foreignValue:     nil,
			through:          "followers",
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

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.col, col)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.isZero, assoc.IsZero())
			assert.Equal(t, test.referenceField, assoc.ReferenceField())
			assert.Equal(t, test.referenceValue, assoc.ReferenceValue())
			assert.Equal(t, test.referenceThrough, assoc.ReferenceThrough())
			assert.Equal(t, test.foreignField, assoc.ForeignField())
			assert.Equal(t, test.foreignThrough, assoc.ForeignThrough())
			assert.Equal(t, test.through, assoc.Through())

			if test.typ == HasMany || test.typ == ManyToMany {
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

func TestAssociation_refNotFound(t *testing.T) {
	type Alpha struct {
		ID int
	}

	type Beta struct {
		ID    int
		Alpha Alpha `ref:"alpha_id" fk:"id"`
	}

	assert.Panics(t, func() {
		NewDocument(&Beta{})
	})
}

func TestAssociation_fkNotFound(t *testing.T) {
	type Alpha struct {
		ID int
	}

	type Beta struct {
		ID    int
		Alpha Alpha
	}

	assert.Panics(t, func() {
		NewDocument(&Beta{})
	})
}
