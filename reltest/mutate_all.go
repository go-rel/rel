package reltest

import (
	"github.com/go-rel/rel"
	"github.com/stretchr/testify/mock"
)

// MutateAll asserts and simulate mutate all function for test.
type MutateAll struct {
	*Expect
}

// Unsafe allows for unsafe operation that doesn't contains where clause.
func (ma *MutateAll) Unsafe() {
	ma.RunFn = nil // clear validation
}

// Result sets the returned number of deleted/updated counts.
func (ma *MutateAll) Result(count int) {
	ma.Return(count, nil)
}

// Error sets error to be returned.
func (ma *MutateAll) Error(err error) {
	ma.Return(0, err)
}

// ConnectionClosed sets this error to be returned.
func (ma *MutateAll) ConnectionClosed() {
	ma.Error(ErrConnectionClosed)
}

func expectMutateAll(r *Repository, methodName string, args ...interface{}) *MutateAll {
	ma := &MutateAll{
		Expect: newExpect(r, methodName,
			args,
			[]interface{}{0, nil},
		),
	}

	// validation
	ma.Run(func(args mock.Arguments) {
		query := args[1].(rel.Query)

		if query.Table == "" {
			panic("reltest: cannot call " + methodName + " without specifying table name. use rel.From(tableName)")
		}

		if query.WhereQuery.None() {
			panic("reltest: unsafe " + methodName + " detected. if you want to mutate all records without filter, please use call .Unsafe()")
		}
	})

	return ma
}

// ExpectUpdateAll to be called.
func ExpectUpdateAll(r *Repository, query rel.Query, mutates []rel.Mutate) *MutateAll {
	return expectMutateAll(r, "UpdateAll", r.ctxData, query, mutates)
}

// ExpectDeleteAll to be called.
func ExpectDeleteAll(r *Repository, query rel.Query) *MutateAll {
	return expectMutateAll(r, "DeleteAll", r.ctxData, query)
}
