package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type deleteAny []*MockDeleteAny

func (da *deleteAny) register(ctxData ctxData, query rel.Query) *MockDeleteAny {
	mda := &MockDeleteAny{
		assert:   &Assert{ctxData: ctxData},
		argQuery: query,
	}
	*da = append(*da, mda)
	return mda
}

func (da deleteAny) execute(ctx context.Context, query rel.Query) (int, error) {
	for _, mda := range da {
		if matchQuery(mda.argQuery, query) &&
			mda.assert.call(ctx) {
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
	assert          *Assert
	unsafe          bool
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
func (mda *MockDeleteAny) Error(err error) *Assert {
	mda.retDeletedCount = 0
	mda.retError = err
	return mda.assert
}

// Success sets no error to be returned.
func (mda *MockDeleteAny) Success() *Assert {
	return mda.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mda *MockDeleteAny) ConnectionClosed() *Assert {
	return mda.Error(ErrConnectionClosed)
}
