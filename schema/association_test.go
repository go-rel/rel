package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssociation(t *testing.T) {
	var (
		trx  = &Transaction{ID: 1}
		user = &User{ID: 1}
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
			data:            trx,
			typ:             BelongsTo,
			target:          &trx.Buyer,
			loaded:          false,
			referenceColumn: "user_id",
			referenceValue:  trx.BuyerID,
			foreignColumn:   "id",
			foreignValue:    trx.Buyer.ID,
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
