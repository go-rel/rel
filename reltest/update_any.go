package reltest

import (
	"context"
	"fmt"

	"github.com/go-rel/rel"
)

type updateAny []*MockUpdateAny

func (ua *updateAny) register(ctxData ctxData, query rel.Query, mutates ...rel.Mutate) *MockUpdateAny {
	mua := &MockUpdateAny{
		assert:     &Assert{ctxData: ctxData},
		ctxData:    ctxData,
		argQuery:   query,
		argMutates: mutates,
	}
	*ua = append(*ua, mua)
	return mua
}

func (ua updateAny) execute(ctx context.Context, query rel.Query, mutates ...rel.Mutate) (int, error) {
	for _, mua := range ua {
		if matchQuery(mua.argQuery, query) &&
			matchMutates(mua.argMutates, mutates) &&
			mua.assert.call(ctx) {
			if query.Table == "" {
				panic("reltest: Cannot call UpdateAny without table. use rel.From(tableName)")
			}

			if !mua.unsafe && query.WhereQuery.None() {
				panic("reltest: unsafe UpdateAny detected. if you want to mutate all records without filter, please use call .Unsafe()")
			}

			return mua.retUpdatedCount, mua.retError
		}
	}

	panic(failExecuteMessage(MockUpdateAny{argQuery: query, argMutates: mutates}, ua))
}

func (ua *updateAny) assert(t T) bool {
	for _, mua := range *ua {
		if !mua.assert.assert(t, mua) {
			return false
		}
	}

	*ua = nil
	return true
}

// MockUpdateAny asserts and simulate UpdateAny function for test.
type MockUpdateAny struct {
	assert          *Assert
	unsafe          bool
	ctxData         ctxData
	argQuery        rel.Query
	argMutates      []rel.Mutate
	retUpdatedCount int
	retError        error
}

// Unsafe allows for unsafe operation to delete records without where condition.
func (mua *MockUpdateAny) Unsafe() *MockUpdateAny {
	mua.unsafe = true
	return mua
}

// UpdatedCount set the returned deleted count of this function.
func (mua *MockUpdateAny) UpdatedCount(updatedCount int) *Assert {
	mua.retUpdatedCount = updatedCount
	return mua.assert
}

// Error sets error to be returned.
func (mua *MockUpdateAny) Error(err error) *Assert {
	mua.retUpdatedCount = 0
	mua.retError = err
	return mua.assert
}

// ConnectionClosed sets this error to be returned.
func (mua *MockUpdateAny) ConnectionClosed() *Assert {
	return mua.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mua MockUpdateAny) String() string {
	argMutates := ""
	for i := range mua.argMutates {
		argMutates += fmt.Sprintf(", %v", mua.argMutates[i])
	}

	return fmt.Sprintf("UpdateAny(ctx, %s%s)", mua.argQuery, argMutates)
}

// ExpectString representation of mocked call.
func (mua MockUpdateAny) ExpectString() string {
	argMutates := ""
	for i := range mua.argMutates {
		argMutates += fmt.Sprintf(", %v", mua.argMutates[i])
	}

	return fmt.Sprintf("ExpectUpdateAny(%s%s)", mua.argQuery, argMutates)
}
