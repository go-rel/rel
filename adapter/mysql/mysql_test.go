package mysql_test

import (
	"github.com/Fs02/grimoire/adapter/mysql"
	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       Query
	}{
		{
			"SELECT * FROM users;",
			nil,
			From("users"),
		},
		{
			"SELECT id, name FROM users;",
			nil,
			From("users").Select("id", "name"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id;",
			nil,
			From("users").Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM users WHERE id=?;",
			[]interface{}{10},
			From("users").Where(Eq(I("id"), 10)),
		},
		{
			"SELECT DISTINCT * FROM users GROUP BY type;",
			nil,
			From("users").Distinct().Group("type"),
		},
		{
			"SELECT * FROM users JOIN transactions ON transactions.id=users.transaction_id HAVING price>?;",
			[]interface{}{1000},
			From("users").Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))).Having(Gt(I("price"), 1000)),
		},
		{
			"SELECT * FROM users ORDER BY created_at ASC;",
			nil,
			From("users").Order(Asc("created_at")),
		},
		{
			"SELECT * FROM users OFFSET 10 LIMIT 10;",
			nil,
			From("users").Offset(10).Limit(10),
		},
	}

	adapter := mysql.Adapter{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := adapter.All(tt.Query)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}
