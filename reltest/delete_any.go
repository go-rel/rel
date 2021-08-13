package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type deleteAny []*MockDeleteAny

func (da *deleteAny) register(ctxData ctxData, query rel.Query) *MockDeleteAny {
	mda := &MockDeleteAny{ctxData: ctxData, argQuery: query}
	*da = append(*da, mda)
	return mda
}

func (da deleteAny) execute(ctx context.Context, query rel.Query) (int, error) {
	for _, mda := range da {
		if fetchContext(ctx) == mda.ctxData && matchQuery(mda.argQuery, query) {
			if query.Table == "" {
				panic("reltest: Cannot call DeleteAny without table. use rel.From(tableName)")
			}

			if !mda.unsafe && query.WhereQuery.None() {
				panic("reltest: unsafe DeleteAny detected. if you want to mutate all records without filter, please use call .Unsafe()")
			}

			return mda.retDeletedCount, mda.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockDeleteAny asserts and simulate DeleteAny function for test.
type MockDeleteAny struct {
	unsafe          bool
	ctxData         ctxData
	argQuery        rel.Query
	retDeletedCount int
	retError        error
}

// Unsafe allows for unsafe operation to delete records without where condition.
func (mda *MockDeleteAny) Unsafe() *MockDeleteAny {
	mda.unsafe = true
	return mda
}

// DeletedCount set the returned deleted count of this function.
func (mda *MockDeleteAny) DeletedCount(deletedCount int) {
	mda.retDeletedCount = deletedCount
}

// Error sets error to be returned.
func (mda *MockDeleteAny) Error(err error) {
	mda.retDeletedCount = 0
	mda.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mda *MockDeleteAny) ConnectionClosed() {
	mda.Error(ErrConnectionClosed)
}
