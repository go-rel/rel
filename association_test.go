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
		autosave       bool
		autoload       bool
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
			autoload:       true,
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
			autoload:       true,
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
			autosave:       true,
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
			autosave:       true,
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
				assoc       = newAssociation(rv, sf.Index)
				doc, loaded = assoc.Document()
			)

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.doc.rt, doc.rt)
			assert.Equal(t, test.doc.meta, doc.meta)
			assert.Equal(t, test.doc.v, doc.v)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.isZero, assoc.IsZero())
			assert.Equal(t, test.referenceField, assoc.ReferenceField())
			assert.Equal(t, test.referenceValue, assoc.ReferenceValue())
			assert.Equal(t, test.foreignField, assoc.ForeignField())
			assert.Equal(t, test.autoload, assoc.Autoload())
			assert.Equal(t, test.autosave, assoc.Autosave())

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

func TestAssociation_Embedded(t *testing.T) {
	type EmbeddedID struct {
		ID int
	}
	type Purchase struct {
		EmbeddedID
	}
	type PurchaseInfo struct {
		PurchaseID int
		Purchase   *Purchase
	}
	type Item struct {
		EmbeddedID
		PurchaseInfo
	}

	var (
		purchase = &Purchase{EmbeddedID: EmbeddedID{ID: 1}}
		item     = &Item{EmbeddedID: EmbeddedID{ID: 1},
			PurchaseInfo: PurchaseInfo{Purchase: purchase, PurchaseID: purchase.ID}}
		rv        = reflect.ValueOf(item)
		sf, _     = rv.Type().Elem().FieldByName("Purchase")
		assoc     = newAssociation(rv, sf.Index)
		_, loaded = assoc.Document()
	)

	assert.Equal(t, AssociationType(BelongsTo), assoc.Type())
	assert.True(t, loaded)
	assert.Equal(t, "purchase_id", assoc.ReferenceField())
	assert.Equal(t, item.ID, assoc.ReferenceValue())
	assert.Equal(t, "id", assoc.ForeignField())
	assert.Equal(t, purchase.ID, assoc.ForeignValue())
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
		autoload         bool
		autosave         bool
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
			record:         "User",
			field:          "Roles",
			data:           user,
			typ:            HasMany,
			col:            NewCollection(&user.Roles),
			loaded:         false,
			isZero:         true,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "id",
			foreignValue:   nil,
			through:        "user_roles",
		},
		{
			record:         "Role",
			field:          "Users",
			data:           role,
			typ:            HasMany,
			col:            NewCollection(&role.Users),
			loaded:         false,
			isZero:         true,
			referenceField: "id",
			referenceValue: role.ID,
			foreignField:   "id",
			foreignValue:   nil,
			through:        "user_roles",
		},
		{
			record:         "User",
			field:          "Followers",
			data:           user,
			typ:            HasMany,
			col:            NewCollection(&user.Followers),
			loaded:         false,
			isZero:         true,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "id",
			foreignValue:   nil,
			through:        "followeds",
		},
		{
			record:         "User",
			field:          "Followings",
			data:           user,
			typ:            HasMany,
			col:            NewCollection(&user.Followings),
			loaded:         false,
			isZero:         true,
			referenceField: "id",
			referenceValue: user.ID,
			foreignField:   "id",
			foreignValue:   nil,
			through:        "follows",
		},
	}

	for _, test := range tests {
		t.Run(test.record+"."+test.field, func(t *testing.T) {
			var (
				rv          = reflect.ValueOf(test.data)
				sf, _       = rv.Type().Elem().FieldByName(test.field)
				assoc       = newAssociation(rv, sf.Index)
				col, loaded = assoc.Collection()
			)

			assert.Equal(t, test.typ, assoc.Type())
			assert.Equal(t, test.col, col)
			assert.Equal(t, test.loaded, loaded)
			assert.Equal(t, test.isZero, assoc.IsZero())
			assert.Equal(t, test.referenceField, assoc.ReferenceField())
			assert.Equal(t, test.referenceValue, assoc.ReferenceValue())
			assert.Equal(t, test.foreignField, assoc.ForeignField())
			assert.Equal(t, test.through, assoc.Through())
			assert.Equal(t, test.autoload, assoc.Autoload())
			assert.Equal(t, test.autosave, assoc.Autosave())

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

func TestAssociation_autosaveWithThrough(t *testing.T) {
	type Alpha struct {
		ID int
	}

	type Beta struct {
		ID    int
		Alpha Alpha `through:"other" autosave:"true"`
	}

	assert.Panics(t, func() {
		NewDocument(&Beta{})
	})
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
