package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type updateAny []*MockUpdateAny

func (ua *updateAny) register(ctxData ctxData, query rel.Query, mutates ...rel.Mutate) *MockUpdateAny {
	mua := &MockUpdateAny{ctxData: ctxData, argQuery: query, argMutates: mutates}
	*ua = append(*ua, mua)
	return mua
}

func (ua updateAny) execute(ctx context.Context, query rel.Query, mutates ...rel.Mutate) (int, error) {
	for _, mua := range ua {
		if fetchContext(ctx) == mua.ctxData && matchQuery(mua.argQuery, query) && matchMutates(mua.argMutates, mutates) {
			if query.Table == "" {
				panic("reltest: Cannot call UpdateAny without table. use rel.From(tableName)")
			}

			if !mua.unsafe && query.WhereQuery.None() {
				panic("reltest: unsafe UpdateAny detected. if you want to mutate all records without filter, please use call .Unsafe()")
			}

			return mua.retUpdatedCount, mua.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockUpdateAny asserts and simulate UpdateAny function for test.
type MockUpdateAny struct {
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
func (mua *MockUpdateAny) UpdatedCount(updatedCount int) {
	mua.retUpdatedCount = updatedCount
}

// Error sets error to be returned.
func (mua *MockUpdateAny) Error(err error) {
	mua.retUpdatedCount = 0
	mua.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mua *MockUpdateAny) ConnectionClosed() {
	mua.Error(ErrConnectionClosed)
}
