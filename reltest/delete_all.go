package reltest

import (
	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// DeleteAll asserts and simulate delete all function for test.
type DeleteAll struct {
	*Expect
}

// Unsafe allows for unsafe delete that doesn't contains where clause.
func (eda *DeleteAll) Unsafe() {
	eda.RunFn = nil // clear validation
}

// ExpectDeleteAll to be called with given field and queries.
func ExpectDeleteAll(r *Repository, query rel.Query) *DeleteAll {
	eda := &DeleteAll{
		Expect: newExpect(r, "DeleteAll",
			[]interface{}{query},
			[]interface{}{nil},
		),
	}

	// validation
	eda.Run(func(args mock.Arguments) {
		query := args[0].(rel.Query)

		if query.Table == "" {
			panic("reltest: cannot call DeleteAll without specifying table name. use rel.From(tableName)")
		}

		if query.WhereQuery.None() {
			panic("reltest: unsafe delete all detected. if you want to delete all records without filter, please use DeleteAll().Unsafe()")
		}
	})

	return eda
}
