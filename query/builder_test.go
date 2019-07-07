package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/group"
	"github.com/Fs02/grimoire/join"
	"github.com/Fs02/grimoire/query"
	"github.com/Fs02/grimoire/sort"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builders [][]query.Builder
		query    query.Query
	}{
		{
			name: "where id=1",
			builders: [][]query.Builder{
				{
					query.FilterEq("id", 1),
				},
				{
					where.Eq("id", 1),
				},
			},
			query: query.Query{
				WhereClause: query.FilterEq("id", 1),
			},
		},
		{
			name: "where id=1 and age<10",
			builders: [][]query.Builder{
				{
					query.FilterEq("id", 1).AndLt("age", 10),
				},
				{
					query.FilterEq("id", 1), query.FilterLt("age", 10),
				},
				{
					where.Eq("id", 1).AndLt("age", 10),
				},
				{
					where.Eq("id", 1), where.Lt("age", 10),
				},
			},
			query: query.Query{
				WhereClause: query.FilterEq("id", 1).AndLt("age", 10),
			},
		},
		{
			name: "where age>10 limit 10 offset 10 order by name asc, age desc",
			builders: [][]query.Builder{
				{
					query.Where(query.FilterGt("age", 10)).Limit(10).Offset(10).Sort("name").SortDesc("age"),
				},
				{
					query.FilterGt("age", 10), query.Limit(10), query.Offset(10), query.NewSortAsc("name"), query.NewSortDesc("age"),
				},
				{
					where.Gt("age", 10), query.Limit(10), query.Offset(10), sort.Asc("name"), sort.Desc("age"),
				},
			},
			query: query.Query{
				WhereClause:  query.FilterGt("age", 10),
				LimitClause:  10,
				OffsetClause: 10,
				SortClause: []query.SortClause{
					query.NewSortAsc("name"),
					query.NewSortDesc("age"),
				},
			},
		},
		{
			name: "select sum(amount), name from transactions join users group by name offset 10 limit 5",
			builders: [][]query.Builder{
				{
					query.From("transactions").Select("sum(amount)", "name").Join("users").Group("name").Having(query.FilterGt("amount", 10)).Offset(10).Limit(5),
				},
				{
					query.From("transactions").Select("sum(amount)", "name"), query.Join("users"), query.Group("name").Having(query.FilterGt("amount", 10)), query.Offset(10), query.Limit(5),
				},
				{
					query.From("transactions").Select("sum(amount)", "name"), join.Join("users"), group.By("name").Having(query.FilterGt("amount", 10)), query.Offset(10), query.Limit(5),
				},
				{
					join.Join("users"), group.By("name").Having(query.FilterGt("amount", 10)), query.From("transactions").Select("sum(amount)", "name").Offset(10).Limit(5),
				},
			},
			query: query.Query{
				SelectClause: query.SelectClause{
					Fields: []string{"sum(amount)", "name"},
				},
				Collection: "transactions",
				JoinClause: []query.JoinClause{
					{
						Mode:       "JOIN",
						Collection: "users",
						From:       "transactions.user_id",
						To:         "users.id",
					},
				},
				GroupClause: query.GroupClause{
					Fields: []string{"name"},
					Filter: query.FilterGt("amount", 10),
				},
				OffsetClause: 10,
				LimitClause:  5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, b := range test.builders {
				assert.Equal(t, test.query, query.Build("", b...))
			}
		})
	}
}
