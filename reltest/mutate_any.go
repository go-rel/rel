package reltest

import (
	"github.com/go-rel/rel"
	"github.com/stretchr/testify/mock"
)

// MutateAny asserts and simulate mutate all function for test.
type MutateAny struct {
	*Expect
}

// Unsafe allows for unsafe operation that doesn't contains where clause.
func (ma *MutateAny) Unsafe() {
	ma.RunFn = nil // clear validation
}

// Result sets the returned number of deleted/updated counts.
func (ma *MutateAny) Result(count int) {
	ma.Return(count, nil)
}

// Error sets error to be returned.
func (ma *MutateAny) Error(err error) {
	ma.Return(0, err)
}

// ConnectionClosed sets this error to be returned.
func (ma *MutateAny) ConnectionClosed() {
	ma.Error(ErrConnectionClosed)
}

func expectMutateAny(r *Repository, methodName string, args ...interface{}) *MutateAny {
	ma := &MutateAny{
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

// ExpectUpdateAny to be called.
func ExpectUpdateAny(r *Repository, query rel.Query, mutates []rel.Mutate) *MutateAny {
	return expectMutateAny(r, "UpdateAny", r.ctxData, query, mutates)
}

// ExpectDeleteAny to be called.
func ExpectDeleteAny(r *Repository, query rel.Query) *MutateAny {
	return expectMutateAny(r, "DeleteAny", r.ctxData, query)
}
