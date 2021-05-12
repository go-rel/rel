package rel_test

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/group"
	"github.com/go-rel/rel/join"
	"github.com/go-rel/rel/sort"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestQuerier(t *testing.T) {
	tests := []struct {
		name     string
		queriers [][]rel.Querier
		query    rel.Query
	}{
		{
			name: "where id=1",
			queriers: [][]rel.Querier{
				{
					where.Eq("id", 1),
				},
				{
					where.Eq("id", 1),
				},
			},
			query: rel.Query{
				WhereQuery:   where.Eq("id", 1),
				CascadeQuery: true,
			},
		},
		{
			name: "where id=1 and age<10",
			queriers: [][]rel.Querier{
				{
					where.Eq("id", 1).AndLt("age", 10),
				},
				{
					where.Eq("id", 1), where.Lt("age", 10),
				},
				{
					where.Eq("id", 1).AndLt("age", 10),
				},
				{
					where.Eq("id", 1), where.Lt("age", 10),
				},
			},
			query: rel.Query{
				WhereQuery:   where.Eq("id", 1).AndLt("age", 10),
				CascadeQuery: true,
			},
		},
		{
			name: "where age>10 limit 10 offset 10 order by name asc, age desc",
			queriers: [][]rel.Querier{
				{
					rel.Where(where.Gt("age", 10)).Limit(10).Offset(10).Sort("name").SortDesc("age"),
				},
				{
					where.Gt("age", 10), rel.Limit(10), rel.Offset(10), rel.NewSortAsc("name"), rel.NewSortDesc("age"),
				},
				{
					where.Gt("age", 10), rel.Limit(10), rel.Offset(10), sort.Asc("name"), sort.Desc("age"),
				},
			},
			query: rel.Query{
				WhereQuery:  where.Gt("age", 10),
				LimitQuery:  10,
				OffsetQuery: 10,
				SortQuery: []rel.SortQuery{
					rel.NewSortAsc("name"),
					rel.NewSortDesc("age"),
				},
				CascadeQuery: true,
			},
		},
		{
			name: "select sum(amount), name from transactions join users group by name offset 10 limit 5",
			queriers: [][]rel.Querier{
				{
					rel.From("transactions").Select("sum(amount)", "name").Join("users").Group("name").Having(where.Gt("amount", 10)).Offset(10).Limit(5),
				},
				{
					rel.From("transactions").Select("sum(amount)", "name"), join.Join("users"), group.By("name").Having(where.Gt("amount", 10)), rel.Offset(10), rel.Limit(5),
				},
				{
					join.Join("users"), group.By("name").Having(where.Gt("amount", 10)), rel.From("transactions").Select("sum(amount)", "name").Offset(10).Limit(5),
				},
			},
			query: rel.Query{
				SelectQuery: rel.SelectQuery{
					Fields: []string{"sum(amount)", "name"},
				},
				Table: "transactions",
				JoinQuery: []rel.JoinQuery{
					{
						Mode:  "JOIN",
						Table: "users",
					},
				},
				GroupQuery: rel.GroupQuery{
					Fields: []string{"name"},
					Filter: where.Gt("amount", 10),
				},
				OffsetQuery:  10,
				LimitQuery:   5,
				CascadeQuery: true,
			},
		},
		{
			name: "where id=1 unscoped",
			queriers: [][]rel.Querier{
				{
					where.Eq("id", 1), rel.Unscoped(true),
				},
			},
			query: rel.Query{
				WhereQuery:    where.Eq("id", 1),
				UnscopedQuery: true,
				CascadeQuery:  true,
			},
		},
		{
			name: "where id=1 preload",
			queriers: [][]rel.Querier{
				{
					where.Eq("id", 1), rel.Preload("users"), rel.Cascade(true),
				},
			},
			query: rel.Query{
				WhereQuery:   where.Eq("id", 1),
				PreloadQuery: []string{"users"},
				CascadeQuery: true,
			},
		},
		{
			name: "where id=1 for update",
			queriers: [][]rel.Querier{
				{
					where.Eq("id", 1), rel.ForUpdate(),
				},
			},
			query: rel.Query{
				WhereQuery:   where.Eq("id", 1),
				LockQuery:    "FOR UPDATE",
				CascadeQuery: true,
			},
		},
		{
			name: "where id=1, from user group age for update",
			queriers: [][]rel.Querier{
				{
					where.Nil("deleted_at"),
					rel.Select("status", "count(id)").Group("status").Where(where.Ne("status", "paid")).Lock("FOR UPDATE"),
				},
			},
			query: rel.Query{
				SelectQuery:  rel.SelectQuery{Fields: []string{"status", "count(id)"}},
				GroupQuery:   rel.GroupQuery{Fields: []string{"status"}},
				WhereQuery:   where.Nil("deleted_at").AndNe("status", "paid"),
				LockQuery:    "FOR UPDATE",
				CascadeQuery: true,
			},
		},
		{
			name: "sql query",
			queriers: [][]rel.Querier{
				{
					rel.SQL("SELECT 1;"),
				},
			},
			query: rel.Query{
				SQLQuery:     rel.SQL("SELECT 1;"),
				CascadeQuery: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, b := range test.queriers {
				assert.Equal(t, test.query, rel.Build("", b...))
			}
		})
	}
}

func TestQuery_Build(t *testing.T) {
	q := rel.From("users").Select("*")
	assert.Equal(t, q, rel.Build("", q))
}

func TestQuery_Select(t *testing.T) {
	assert.Equal(t, rel.Query{
		Table: "users",
		SelectQuery: rel.SelectQuery{
			Fields: []string{"*"},
		},
		CascadeQuery: true,
	}, rel.From("users").Select("*"))

	assert.Equal(t, rel.Query{
		Table: "users",
		SelectQuery: rel.SelectQuery{
			Fields: []string{"id", "name", "email"},
		},
		CascadeQuery: true,
	}, rel.From("users").Select("id", "name", "email"))
}

func TestQuery_Distinct(t *testing.T) {
	assert.Equal(t, rel.Query{
		Table: "users",
		SelectQuery: rel.SelectQuery{
			Fields:       []string{"*"},
			OnlyDistinct: true,
		},
		CascadeQuery: true,
	}, rel.From("users").Select("*").Distinct())
}

func TestQuery_Join(t *testing.T) {
	result := rel.Query{
		Table: "users",
		JoinQuery: []rel.JoinQuery{
			{
				Mode:  "JOIN",
				Table: "transactions",
			},
		},
		CascadeQuery: true,
	}

	assert.Equal(t, result, rel.Build("", rel.From("users").Join("transactions")))
	assert.Equal(t, result, rel.Build("", rel.Join("transactions").From("users")))
	assert.Equal(t, result, rel.Build("users", rel.Join("transactions")))
}

func TestQuery_JoinOn(t *testing.T) {
	result := rel.Query{
		Table: "users",
		JoinQuery: []rel.JoinQuery{
			{
				Mode:  "JOIN",
				Table: "transactions",
				From:  "users.transaction_id",
				To:    "transactions.id",
			},
		},
		CascadeQuery: true,
	}

	assert.Equal(t, result, rel.From("users").JoinOn("transactions", "users.transaction_id", "transactions.id"))
	assert.Equal(t, result, rel.JoinOn("transactions", "users.transaction_id", "transactions.id").From("users"))
}

func TestQuery_Joinf(t *testing.T) {
	result := rel.Query{
		Table: "users",
		JoinQuery: []rel.JoinQuery{
			{
				Mode:      "JOIN transactions ON transacations.id=?",
				Arguments: []interface{}{1},
			},
		},
		CascadeQuery: true,
	}

	assert.Equal(t, result, rel.From("users").Joinf("JOIN transactions ON transacations.id=?", 1))
	assert.Equal(t, result, rel.Joinf("JOIN transactions ON transacations.id=?", 1).From("users"))
}

func TestQuery_Where(t *testing.T) {
	tests := []struct {
		Case     string
		Build    rel.Query
		Expected rel.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`where id=1 AND deleted_at IS NIL`,
			rel.Where(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).Where(where.Ne("active", false)),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND deleted_at IS NIL (where package)`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND deleted_at IS NIL (chained where package)`,
			rel.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`id=1`,
			rel.From("users").Wheref("id=?", 1),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Fragment("id=?", 1)),
				CascadeQuery: true,
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
		Build    rel.Query
		Expected rel.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			rel.From("users").OrWhere(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			rel.From("users").Where(where.Eq("id", 1)).OrWhere(where.Nil("deleted_at")),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			rel.Where(where.Eq("id", 1)).OrWhere(where.Nil("deleted_at")),
			rel.Query{
				WhereQuery:   where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
				CascadeQuery: true,
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false)),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.Ne("active", false)),
				CascadeQuery: true,
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))),
				CascadeQuery: true,
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			rel.From("users").Where(where.Eq("id", 1), where.Nil("deleted_at")).OrWhere(where.Ne("active", false), where.Gte("score", 80)).Where(where.Lt("price", 10000)),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
				CascadeQuery: true,
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000 (chained where package)`,
			rel.From("users").Where(where.Eq("id", 1).AndNil("deleted_at")).OrWhere(where.Ne("active", false).AndGte("score", 80)).Where(where.Lt("price", 10000)),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
				CascadeQuery: true,
			},
		},
		{
			`id=1`,
			rel.From("users").Where(where.Nil("deleted_at")).OrWheref("id=?", 1),
			rel.Query{
				Table:        "users",
				WhereQuery:   where.Or(where.Nil("deleted_at"), where.Fragment("id=?", 1)),
				CascadeQuery: true,
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
	result := rel.Query{
		Table: "users",
		GroupQuery: rel.GroupQuery{
			Fields: []string{"active", "plan"},
		},
		CascadeQuery: true,
	}

	assert.Equal(t, result, rel.From("users").Group("active", "plan"))
}

func TestQuery_Having(t *testing.T) {
	tests := []struct {
		Case     string
		Build    rel.Query
		Expected rel.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				},
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).Having(where.Ne("active", false)),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
				},
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false (chained where package)`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1).AndNil("deleted_at")).Having(where.Ne("active", false)),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at"), where.Ne("active", false)),
				},
				CascadeQuery: true,
			},
		},
		{
			`id=1`,
			rel.From("users").Group("active", "plan").Havingf("id=?", 1),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Fragment("id=?", 1)),
				},
				CascadeQuery: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrHaving(t *testing.T) {
	tests := []struct {
		Case     string
		Build    rel.Query
		Expected rel.Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			rel.From("users").Group("active", "plan").OrHaving(where.Eq("id", 1), where.Nil("deleted_at")),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Eq("id", 1), where.Nil("deleted_at")),
				},
				CascadeQuery: true,
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1)).OrHaving(where.Nil("deleted_at")),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.Eq("id", 1), where.Nil("deleted_at")),
				},
				CascadeQuery: true,
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false)),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.Ne("active", false)),
				},
				CascadeQuery: true,
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false), where.Gte("score", 80)),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))),
				},
				CascadeQuery: true,
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			rel.From("users").Group("active", "plan").Having(where.Eq("id", 1), where.Nil("deleted_at")).OrHaving(where.Ne("active", false), where.Gte("score", 80)).Having(where.Lt("price", 10000)),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Or(where.And(where.Eq("id", 1), where.Nil("deleted_at")), where.And(where.Ne("active", false), where.Gte("score", 80))), where.Lt("price", 10000)),
				},
				CascadeQuery: true,
			},
		},
		{
			`id=1 AND`,
			rel.From("users").Group("active", "plan").OrHavingf("id=?", 1),
			rel.Query{
				Table: "users",
				GroupQuery: rel.GroupQuery{
					Fields: []string{"active", "plan"},
					Filter: where.And(where.Fragment("id=?", 1)),
				},
				CascadeQuery: true,
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
		Build    rel.Query
		Expected rel.Query
	}{
		{
			"Sort",
			rel.From("users").Sort("id"),
			rel.Query{
				Table: "users",
				SortQuery: []rel.SortQuery{
					{
						Field: "id",
						Sort:  1,
					},
				},
				CascadeQuery: true,
			},
		},
		{
			"SortAsc",
			rel.From("users").SortAsc("id", "name"),
			rel.Query{
				Table: "users",
				SortQuery: []rel.SortQuery{
					{
						Field: "id",
						Sort:  1,
					},
					{
						Field: "name",
						Sort:  1,
					},
				},
				CascadeQuery: true,
			},
		},
		{
			"SortAsc",
			rel.From("users").SortAsc("id", "name").SortDesc("age", "created_at"),
			rel.Query{
				Table: "users",
				SortQuery: []rel.SortQuery{
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
				CascadeQuery: true,
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
	assert.Equal(t, rel.Query{
		Table:        "users",
		OffsetQuery:  10,
		CascadeQuery: true,
	}, rel.From("users").Offset(10))
}

func TestQuery_Limit(t *testing.T) {
	assert.Equal(t, rel.Query{
		Table:        "users",
		LimitQuery:   10,
		CascadeQuery: true,
	}, rel.From("users").Limit(10))
}

func TestQuery_Lock_outsideTransaction(t *testing.T) {
	assert.Equal(t, rel.Query{
		Table:        "users",
		LockQuery:    "FOR UPDATE",
		CascadeQuery: true,
	}, rel.From("users").Lock("FOR UPDATE"))
}
