package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/group"
	"github.com/Fs02/grimoire/join"
	"github.com/Fs02/grimoire/sort"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func TestQuerier(t *testing.T) {
	tests := []struct {
		name     string
		queriers [][]grimoire.Querier
		query    grimoire.Query
	}{
		{
			name: "where id=1",
			queriers: [][]grimoire.Querier{
				{
					where.Eq("id", 1),
				},
				{
					where.Eq("id", 1),
				},
			},
			query: grimoire.Query{
				WhereClause: where.Eq("id", 1),
			},
		},
		{
			name: "where id=1 and age<10",
			queriers: [][]grimoire.Querier{
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
			query: grimoire.Query{
				WhereClause: where.Eq("id", 1).AndLt("age", 10),
			},
		},
		{
			name: "where age>10 limit 10 offset 10 order by name asc, age desc",
			queriers: [][]grimoire.Querier{
				{
					grimoire.Where(where.Gt("age", 10)).Limit(10).Offset(10).Sort("name").SortDesc("age"),
				},
				{
					where.Gt("age", 10), grimoire.Limit(10), grimoire.Offset(10), grimoire.NewSortAsc("name"), grimoire.NewSortDesc("age"),
				},
				{
					where.Gt("age", 10), grimoire.Limit(10), grimoire.Offset(10), sort.Asc("name"), sort.Desc("age"),
				},
			},
			query: grimoire.Query{
				WhereClause:  where.Gt("age", 10),
				LimitClause:  10,
				OffsetClause: 10,
				SortQuery: []grimoire.SortQuery{
					grimoire.NewSortAsc("name"),
					grimoire.NewSortDesc("age"),
				},
			},
		},
		{
			name: "select sum(amount), name from transactions join users group by name offset 10 limit 5",
			queriers: [][]grimoire.Querier{
				{
					grimoire.From("transactions").Select("sum(amount)", "name").Join("users").Group("name").Having(where.Gt("amount", 10)).Offset(10).Limit(5),
				},
				{
					grimoire.From("transactions").Select("sum(amount)", "name"), grimoire.Join("users"), grimoire.Group("name").Having(where.Gt("amount", 10)), grimoire.Offset(10), grimoire.Limit(5),
				},
				{
					grimoire.From("transactions").Select("sum(amount)", "name"), join.Join("users"), group.By("name").Having(where.Gt("amount", 10)), grimoire.Offset(10), grimoire.Limit(5),
				},
				{
					join.Join("users"), group.By("name").Having(where.Gt("amount", 10)), grimoire.From("transactions").Select("sum(amount)", "name").Offset(10).Limit(5),
				},
			},
			query: grimoire.Query{
				SelectQuery: grimoire.SelectQuery{
					Fields: []string{"sum(amount)", "name"},
				},
				Collection: "transactions",
				JoinQuery: []grimoire.JoinQuery{
					{
						Mode:       "JOIN",
						Collection: "users",
						From:       "transactions.user_id",
						To:         "users.id",
					},
				},
				GroupQuery: grimoire.GroupQuery{
					Fields: []string{"name"},
					Filter: where.Gt("amount", 10),
				},
				OffsetClause: 10,
				LimitClause:  5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, b := range test.queriers {
				assert.Equal(t, test.query, grimoire.BuildQuery("", b...))
			}
		})
	}
}
