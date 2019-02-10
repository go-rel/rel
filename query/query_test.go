package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Select(t *testing.T) {
	assert.Equal(t, query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"*"},
		},
	}, query.From("users").Select("*"))

	assert.Equal(t, query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"id", "name", "email"},
		},
	}, query.From("users").Select("id", "name", "email"))
}

func TestQuery_Distinct(t *testing.T) {
	assert.Equal(t, query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields:       []string{"users.*"},
			OnlyDistinct: true,
		},
	}, query.From("users").Distinct())
}

func TestQuery_Join(t *testing.T) {
	result := query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		JoinClause: []query.JoinClause{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				From:       "users.transaction_id",
				To:         "transactions.id",
			},
		},
	}

	assert.Equal(t, result, query.From("users").Join("transactions"))
	assert.Equal(t, result, query.Join("transactions").From("users"))
}

func TestQuery_JoinOn(t *testing.T) {
	result := query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		JoinClause: []query.JoinClause{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				From:       "users.transaction_id",
				To:         "transactions.id",
			},
		},
	}

	assert.Equal(t, result, query.From("users").JoinOn("transactions", "users.transaction_id", "transactions.id"))
	assert.Equal(t, result, query.JoinOn("transactions", "users.transaction_id", "transactions.id").From("users"))
}

func TestQuery_JoinFragment(t *testing.T) {
	result := query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		JoinClause: []query.JoinClause{
			{
				Mode:      "JOIN transactions ON transacations.id=?",
				Arguments: []interface{}{1},
			},
		},
	}

	assert.Equal(t, result, query.From("users").JoinFragment("JOIN transactions ON transacations.id=?", 1))
	assert.Equal(t, result, query.JoinFragment("JOIN transactions ON transacations.id=?", 1).From("users"))
}

func TestQuery_Where(t *testing.T) {
	tests := []struct {
		Case     string
		Build    query.Query
		Expected query.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			query.From("users").Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`where id=1 AND deleted_at IS NIL`,
			query.Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			query.From("users").Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")).Where(query.FilterNe("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at"), query.FilterNe("active", false)),
			},
		},
		{
			`id=1 AND deleted_at IS NIL (where package)`,
			query.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL (chained where package)`,
			query.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrWhere(t *testing.T) {
	tests := []struct {
		Case     string
		Build    query.Query
		Expected query.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			query.From("users").OrWhere(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			query.From("users").Where(query.FilterEq("id", 1)).OrWhere(query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterOr(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			query.Where(query.FilterEq("id", 1)).OrWhere(query.FilterNil("deleted_at")),
			query.Query{
				WhereClause: query.FilterOr(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			query.From("users").Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrWhere(query.FilterNe("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterNe("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			query.From("users").Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrWhere(query.FilterNe("active", false), query.FilterGte("score", 80)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			query.From("users").Where(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrWhere(query.FilterNe("active", false), query.FilterGte("score", 80)).Where(query.FilterLt("price", 10000)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))), query.FilterLt("price", 10000)),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000 (where package)`,
			query.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)).Where(where.Lt("price", 10000)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))), query.FilterLt("price", 10000)),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000 (chained where package)`,
			query.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")).OrWhere(where.Ne("active", false).AndGte("score", 80)).Where(where.Lt("price", 10000)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				WhereClause: query.FilterAnd(query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))), query.FilterLt("price", 10000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_Group(t *testing.T) {
	result := query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		GroupClause: query.GroupClause{
			Fields: []string{"active", "plan"},
		},
	}

	assert.Equal(t, result, query.From("users").Group("active", "plan"))
	assert.Equal(t, result, query.Group("active", "plan").From("users"))
}

func TestQuery_Having(t *testing.T) {
	tests := []struct {
		Case     string
		Build    query.Query
		Expected query.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")).Having(query.FilterNe("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at"), query.FilterNe("active", false)),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false (where package)`,
			query.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).Having(where.Ne("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at"), query.FilterNe("active", false)),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false (chained where package)`,
			query.From("users").Group("active", "plan").Having(where.Eq("id", 1).AndNil("deleted_at")).Having(where.Ne("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at"), query.FilterNe("active", false)),
				},
			},
		}}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrHaving(t *testing.T) {
	tests := []struct {
		Case     string
		Build    query.Query
		Expected query.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			query.From("users").Group("active", "plan").OrHaving(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
				},
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1)).OrHaving(query.FilterNil("deleted_at")),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterOr(query.FilterEq("id", 1), query.FilterNil("deleted_at")),
				},
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrHaving(query.FilterNe("active", false)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterNe("active", false)),
				},
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrHaving(query.FilterNe("active", false), query.FilterGte("score", 80)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))),
				},
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			query.From("users").Group("active", "plan").Having(query.FilterEq("id", 1), query.FilterNil("deleted_at")).OrHaving(query.FilterNe("active", false), query.FilterGte("score", 80)).Having(query.FilterLt("price", 10000)),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"active", "plan"},
					Filter: query.FilterAnd(query.FilterOr(query.FilterAnd(query.FilterEq("id", 1), query.FilterNil("deleted_at")), query.FilterAnd(query.FilterNe("active", false), query.FilterGte("score", 80))), query.FilterLt("price", 10000)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_Sort(t *testing.T) {
	tests := []struct {
		Case     string
		Build    query.Query
		Expected query.Query
	}{
		{
			"Sort",
			query.From("users").Sort("id"),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				SortClause: []query.SortClause{
					{
						Field: "id",
						Sort:  1,
					},
				},
			},
		},
		{
			"SortAsc",
			query.From("users").SortAsc("id", "name"),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				SortClause: []query.SortClause{
					{
						Field: "id",
						Sort:  1,
					},
					{
						Field: "name",
						Sort:  1,
					},
				},
			},
		},
		{
			"SortAsc",
			query.From("users").SortAsc("id", "name").SortDesc("age", "created_at"),
			query.Query{
				Collection: "users",
				SelectClause: query.SelectClause{
					Fields: []string{"users.*"},
				},
				SortClause: []query.SortClause{
					{
						Field: "id",
						Sort:  1,
					},
					{
						Field: "name",
						Sort:  1,
					},
					{
						Field: "age",
						Sort:  -1,
					},
					{
						Field: "created_at",
						Sort:  -1,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_Offset(t *testing.T) {
	assert.Equal(t, query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		OffsetClause: 10,
	}, query.From("users").Offset(10))
}

func TestQuery_Limit(t *testing.T) {
	assert.Equal(t, query.Query{
		Collection: "users",
		SelectClause: query.SelectClause{
			Fields: []string{"users.*"},
		},
		LimitClause: 10,
	}, query.From("users").Limit(10))
}

// func TestQuery_Lock_outsideTransaction(t *testing.T) {
// 	assert.Equal(t, repo.From("users").Lock(), Query{
// 		repo:       &repo,
// 		Collection: "users",
// 		Fields:     []string{"users.*"},
// 	})

// 	assert.Equal(t, repo.From("users").Lock("FOR SHARE"), Query{
// 		repo:       &repo,
// 		Collection: "users",
// 		Fields:     []string{"users.*"},
// 	})
// }

// func TestQuery_Lock_insideTransaction(t *testing.T) {
// 	repo := Repo{inTransaction: true}
// 	assert.Equal(t, repo.From("users").Lock(), Query{
// 		repo:       &repo,
// 		Collection: "users",
// 		Fields:     []string{"users.*"},
// 		LockClause: "FOR UPDATE",
// 	})

// 	assert.Equal(t, repo.From("users").Lock("FOR SHARE"), Query{
// 		repo:       &repo,
// 		Collection: "users",
// 		Fields:     []string{"users.*"},
// 		LockClause: "FOR SHARE",
// 	})
// }
