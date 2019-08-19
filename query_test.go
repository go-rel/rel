package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Build(t *testing.T) {
	q := grimoire.From("users").Select("*")
	assert.Equal(t, q, grimoire.BuildQuery("", q))
}

func TestQuery_Select(t *testing.T) {
	assert.Equal(t, grimoire.Query{
		Collection: "users",
		SelectQuery: grimoire.SelectQuery{
			Fields: []string{"*"},
		},
	}, grimoire.From("users").Select("*"))

	assert.Equal(t, grimoire.Query{
		Collection: "users",
		SelectQuery: grimoire.SelectQuery{
			Fields: []string{"id", "name", "email"},
		},
	}, grimoire.From("users").Select("id", "name", "email"))
}

func TestQuery_Distinct(t *testing.T) {
	assert.Equal(t, grimoire.Query{
		Collection: "users",
		SelectQuery: grimoire.SelectQuery{
			Fields:       []string{"*"},
			OnlyDistinct: true,
		},
	}, grimoire.From("users").Select("*").Distinct())
}

func TestQuery_Join(t *testing.T) {
	result := grimoire.Query{
		Collection: "users",
		JoinQuery: []grimoire.JoinQuery{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				From:       "users.transaction_id",
				To:         "transactions.id",
			},
		},
	}

	assert.Equal(t, result, grimoire.BuildQuery("", grimoire.From("users").Join("transactions")))
	assert.Equal(t, result, grimoire.BuildQuery("", grimoire.Join("transactions").From("users")))
	assert.Equal(t, result, grimoire.BuildQuery("users", grimoire.Join("transactions")))
}

func TestQuery_JoinOn(t *testing.T) {
	result := grimoire.Query{
		Collection: "users",
		JoinQuery: []grimoire.JoinQuery{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				From:       "users.transaction_id",
				To:         "transactions.id",
			},
		},
	}

	assert.Equal(t, result, grimoire.From("users").JoinOn("transactions", "users.transaction_id", "transactions.id"))
	assert.Equal(t, result, grimoire.JoinOn("transactions", "users.transaction_id", "transactions.id").From("users"))
}

func TestQuery_JoinFragment(t *testing.T) {
	result := grimoire.Query{
		Collection: "users",
		JoinQuery: []grimoire.JoinQuery{
			{
				Mode:      "JOIN transactions ON transacations.id=?",
				Arguments: []interface{}{1},
			},
		},
	}

	assert.Equal(t, result, grimoire.From("users").JoinFragment("JOIN transactions ON transacations.id=?", 1))
	assert.Equal(t, result, grimoire.JoinFragment("JOIN transactions ON transacations.id=?", 1).From("users"))
}

func TestQuery_Where(t *testing.T) {
	tests := []struct {
		Case     string
		Build    grimoire.Query
		Expected grimoire.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`where id=1 AND deleted_at IS NIL`,
			grimoire.Where(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).Where(where.Ne("active", false)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
			},
		},
		{
			`id=1 AND deleted_at IS NIL (where package)`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL (chained where package)`,
			grimoire.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
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
		Build    grimoire.Query
		Expected grimoire.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			grimoire.From("users").OrWhere(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			grimoire.From("users").Where(where.Eq("id", 1)).OrWhere(where.Nil("deleted_at")),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			grimoire.Where(where.Eq("id", 1)).OrWhere(where.Nil("deleted_at")),
			grimoire.Query{
				WhereClause: where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)).Where(where.Lt("price", 10000)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000 (where package)`,
			grimoire.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)).Where(where.Lt("price", 10000)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000 (chained where package)`,
			grimoire.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")).OrWhere(where.Ne("active", false).AndGte("score", 80)).Where(where.Lt("price", 10000)),
			grimoire.Query{
				Collection:  "users",
				WhereClause: where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
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
	result := grimoire.Query{
		Collection: "users",
		GroupQuery: grimoire.GroupQuery{
			Fields: []string{"active", "plan"},
		},
	}

	assert.Equal(t, result, grimoire.From("users").Group("active", "plan"))
	assert.Equal(t, result, grimoire.Group("active", "plan").From("users"))
}

func TestQuery_Having(t *testing.T) {
	tests := []struct {
		Case     string
		Build    grimoire.Query
		Expected grimoire.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).Having(where.Ne("active", false)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false (where package)`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).Having(where.Ne("active", false)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
				},
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false (chained where package)`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1).AndNil("deleted_at")).Having(where.Ne("active", false)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
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
		Build    grimoire.Query
		Expected grimoire.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			grimoire.From("users").Group("active", "plan").OrHaving(where.Eq("id", 1), where.Nil("deleted_at")),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				},
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1)).OrHaving(where.Nil("deleted_at")),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
				},
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.Ne("active", false)),
				},
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false), where.Gte("score", 80)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))),
				},
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			grimoire.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false), where.Gte("score", 80)).Having(where.Lt("price", 10000)),
			grimoire.Query{
				Collection: "users",
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
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
		Build    grimoire.Query
		Expected grimoire.Query
	}{
		{
			"Sort",
			grimoire.From("users").Sort("id"),
			grimoire.Query{
				Collection: "users",
				SortQuery: []grimoire.SortQuery{
					{
						Field: "id",
						Sort:  1,
					},
				},
			},
		},
		{
			"SortAsc",
			grimoire.From("users").SortAsc("id", "name"),
			grimoire.Query{
				Collection: "users",
				SortQuery: []grimoire.SortQuery{
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
			grimoire.From("users").SortAsc("id", "name").SortDesc("age", "created_at"),
			grimoire.Query{
				Collection: "users",
				SortQuery: []grimoire.SortQuery{
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
	assert.Equal(t, grimoire.Query{
		Collection:   "users",
		OffsetClause: 10,
	}, grimoire.From("users").Offset(10))
}

func TestQuery_Limit(t *testing.T) {
	assert.Equal(t, grimoire.Query{
		Collection:  "users",
		LimitClause: 10,
	}, grimoire.From("users").Limit(10))
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
