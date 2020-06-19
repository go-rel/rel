package reltest

import (
	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// MutateAll asserts and simulate mutate all function for test.
type MutateAll struct {
	*Expect
}

// Unsafe allows for unsafe operation that doesn't contains where clause.
func (ema *MutateAll) Unsafe() {
	ema.RunFn = nil // clear validation
}

func expectMutateAll(r *Repository, methodName string, args ...interface{}) *MutateAll {
	ema := &MutateAll{
		Expect: newExpect(r, methodName,
			args,
			[]interface{}{nil},
		),
	}

	// validation
	ema.Run(func(args mock.Arguments) {
		query := args[1].(rel.Query)

		if query.Table == "" {
			panic("reltest: cannot call " + methodName + " without specifying table name. use rel.From(tableName)")
		}

		if query.WhereQuery.None() {
			panic("reltest: unsafe " + methodName + " detected. if you want to mutate all records without filter, please use call .Unsafe()")
		}
	})

	return ema
}

// ExpectUpdateAll to be called.
func ExpectUpdateAll(r *Repository, query rel.Query, mutates []rel.Mutate) *MutateAll {
	return expectMutateAll(r, "UpdateAll", r.ctxData, query, mutates)
}

// ExpectDeleteAll to be called.
func ExpectDeleteAll(r *Repository, query rel.Query) *MutateAll {
	return expectMutateAll(r, "DeleteAll", r.ctxData, query)
}
