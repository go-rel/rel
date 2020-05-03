package sql

import (
	"fmt"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/sort"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBuilder_Find(b *testing.B) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	for n := 0; n < b.N; n++ {
		query := rel.From("users").
			Select("id", "name").
			Join("transactions").
			Where(where.Eq("id", 10), where.In("status", 1, 2, 3)).
			Group("type").Having(where.Gt("price", 1000)).
			SortAsc("created_at").SortDesc("id").
			Offset(10).Limit(10)

		builder.Find(query)
	}
}

func TestBuilder_Find(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		query = rel.From("users")
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       rel.Query
	}{
		{
			"SELECT * FROM `users`;",
			nil,
			query,
		},
		{
			"SELECT `users`.* FROM `users`;",
			nil,
			query.Select("users.*"),
		},
		{
			"SELECT `id`,`name` FROM `users`;",
			nil,
			query.Select("id", "name"),
		},
		{
			"SELECT `id`,FIELD(`gender`, \"male\") AS `order` FROM `users` ORDER BY `order` ASC;",
			nil,
			query.Select("id", "^FIELD(`gender`, \"male\") AS `order`").SortAsc("order"),
		},
		{
			"SELECT * FROM `users` JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			query.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM `users` WHERE `id`=?;",
			[]interface{}{10},
			query.Where(where.Eq("id", 10)),
		},
		{
			"SELECT DISTINCT * FROM `users` GROUP BY `type` HAVING `price`>?;",
			[]interface{}{1000},
			query.Distinct().Group("type").Having(where.Gt("price", 1000)),
		},
		{
			"SELECT * FROM `users` INNER JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			query.JoinWith("INNER JOIN", "transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC;",
			nil,
			query.SortAsc("created_at"),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC, `id` DESC;",
			nil,
			query.SortAsc("created_at").SortDesc("id"),
		},
		{
			"SELECT * FROM `users` LIMIT 10 OFFSET 10;",
			nil,
			query.Offset(10).Limit(10),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.Find(test.Query)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Find_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		query = rel.From("users")
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       rel.Query
	}{
		{
			"SELECT * FROM \"users\";",
			nil,
			query,
		},
		{
			"SELECT \"users\".* FROM \"users\";",
			nil,
			query.Select("users.*"),
		},
		{
			"SELECT \"id\",\"name\" FROM \"users\";",
			nil,
			query.Select("id", "name"),
		},
		{
			"SELECT * FROM \"users\" JOIN \"transactions\" ON \"transactions\".\"id\"=\"users\".\"transaction_id\";",
			nil,
			query.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM \"users\" WHERE \"id\"=$1;",
			[]interface{}{10},
			query.Where(where.Eq("id", 10)),
		},
		{
			"SELECT DISTINCT * FROM \"users\" GROUP BY \"type\" HAVING \"price\">$1;",
			[]interface{}{1000},
			query.Distinct().Group("type").Having(where.Gt("price", 1000)),
		},
		{
			"SELECT * FROM \"users\" JOIN \"transactions\" ON \"transactions\".\"id\"=\"users\".\"transaction_id\";",
			nil,
			query.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC;",
			nil,
			query.SortAsc("created_at"),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC, \"id\" DESC;",
			nil,
			query.SortAsc("created_at").SortDesc("id"),
		},
		{
			"SELECT * FROM \"users\" LIMIT 10 OFFSET 10;",
			nil,
			query.Offset(10).Limit(10),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.Find(test.Query)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Find_SQLQuery(t *testing.T) {
	var (
		config   = &Config{}
		builder  = NewBuilder(config)
		query    = rel.Build("", rel.SQL("SELECT * FROM `users` WHERE id=?;", 1))
		qs, args = builder.Find(query)
	)

	assert.Equal(t, "SELECT * FROM `users` WHERE id=?;", qs)
	assert.Equal(t, []interface{}{1}, args)
}

func BenchmarkBuilder_Aggregate(b *testing.B) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	for n := 0; n < b.N; n++ {
		builder.Aggregate(rel.From("users").Group("gender"), "sum", "transactions.total")
	}
}

func TestBuilder_Aggregate(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
		query   = rel.From("users")
	)

	qs, args := builder.Aggregate(query, "count", "*")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT count(*) AS count FROM `users`;", qs)

	qs, args = builder.Aggregate(query, "sum", "transactions.total")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT sum(`transactions`.`total`) AS sum FROM `users`;", qs)

	qs, args = builder.Aggregate(query.Group("gender"), "sum", "transactions.total")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT sum(`transactions`.`total`) AS sum,`gender` FROM `users` GROUP BY `gender`;", qs)
}

func BenchmarkBuilder_Insert(b *testing.B) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		mutates = map[string]rel.Mutate{
			"name":  rel.Set("name", "foo"),
			"age":   rel.Set("age", 10),
			"agree": rel.Set("agree", true),
		}
		builder = NewBuilder(config)
	)

	for n := 0; n < b.N; n++ {
		builder.Insert("users", mutates)
	}
}

func TestBuilder_Insert(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{
			"name":  rel.Set("name", "foo"),
			"age":   rel.Set("age", 10),
			"agree": rel.Set("agree", true),
		}
		qs, args = builder.Insert("users", mutates)
	)

	assert.Regexp(t, fmt.Sprint(`^INSERT INTO `, "`users`", ` \((`, "`", `\w*`, "`", `,?){3}\) VALUES \(\?,\?,\?\);`), qs)
	assert.Contains(t, qs, "name")
	assert.Contains(t, qs, "age")
	assert.Contains(t, qs, "agree")
	assert.ElementsMatch(t, []interface{}{"foo", 10, true}, args)
}

func TestBuilder_Insert_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{
			"name":  rel.Set("name", "foo"),
			"age":   rel.Set("age", 10),
			"agree": rel.Set("agree", true),
		}
		qs, args = builder.Returning("id").Insert("users", mutates)
	)

	assert.Regexp(t, `^INSERT INTO \"users\" \(("\w*",?){3}\) VALUES \(\$1,\$2,\$3\) RETURNING \"id\";`, qs)
	assert.Contains(t, qs, "name")
	assert.Contains(t, qs, "age")
	assert.Contains(t, qs, "agree")
	assert.ElementsMatch(t, []interface{}{"foo", 10, true}, args)
}

func TestBuilder_Insert_defaultValuesDisabled(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "?",
			EscapeChar:          "`",
			InsertDefaultValues: false,
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{}
		qs, args = builder.Insert("users", mutates)
	)

	assert.Equal(t, "INSERT INTO `users` () VALUES ();", qs)
	assert.Equal(t, []interface{}{}, args)
}

func TestBuilder_Insert_defaultValuesEnabled(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "?",
			InsertDefaultValues: true,
			EscapeChar:          "`",
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{}
		qs, args = builder.Returning("id").Insert("users", mutates)
	)

	assert.Equal(t, "INSERT INTO `users` DEFAULT VALUES RETURNING `id`;", qs)
	assert.Nil(t, args)
}

func BenchmarkBuilder_InsertAll(b *testing.B) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder      = NewBuilder(config)
		bulkMutates = []map[string]rel.Mutate{
			{
				"name": rel.Set("name", "foo"),
			},
			{
				"age": rel.Set("age", 10),
			},
			{
				"name": rel.Set("name", "boo"),
				"age":  rel.Set("age", 20),
			},
		}
	)

	for n := 0; n < b.N; n++ {
		builder.InsertAll("users", []string{"name"}, bulkMutates)
	}
}

func TestBuilder_InsertAll(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder      = NewBuilder(config)
		bulkMutates = []map[string]rel.Mutate{
			{
				"name": rel.Set("name", "foo"),
			},
			{
				"age": rel.Set("age", 10),
			},
			{
				"name": rel.Set("name", "boo"),
				"age":  rel.Set("age", 20),
			},
		}
	)

	statement, args := builder.InsertAll("users", []string{"name"}, bulkMutates)
	assert.Equal(t, "INSERT INTO `users` (`name`) VALUES (?),(DEFAULT),(?);", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// with age
	statement, args = builder.InsertAll("users", []string{"name", "age"}, bulkMutates)
	assert.Equal(t, "INSERT INTO `users` (`name`,`age`) VALUES (?,DEFAULT),(DEFAULT,?),(?,?);", statement)
	assert.Equal(t, []interface{}{"foo", 10, "boo", 20}, args)
}

func TestBuilder_InsertAll_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder      = NewBuilder(config)
		bulkMutates = []map[string]rel.Mutate{
			{
				"name": rel.Set("name", "foo"),
			},
			{
				"age": rel.Set("age", 10),
			},
			{
				"name": rel.Set("name", "boo"),
				"age":  rel.Set("age", 20),
			},
		}
	)

	statement, args := builder.Returning("id").InsertAll("users", []string{"name"}, bulkMutates)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\") VALUES ($1),(DEFAULT),($2) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// with age
	builder.count = 0
	statement, args = builder.Returning("id").InsertAll("users", []string{"name", "age"}, bulkMutates)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,DEFAULT),(DEFAULT,$2),($3,$4) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", 10, "boo", 20}, args)
}

func TestBuilder_Update(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{
			"name":  rel.Set("name", "foo"),
			"age":   rel.Set("age", 10),
			"agree": rel.Set("agree", true),
		}
	)

	qs, qargs := builder.Update("users", mutates, where.And())
	assert.Regexp(t, fmt.Sprint("UPDATE `users` SET `", `\w*`, "`=", `\?`, ",`", `\w*`, "`=", `\?`, ",`", `\w*`, "`=", `\?`, ";"), qs)
	assert.ElementsMatch(t, []interface{}{"foo", 10, true}, qargs)

	qs, qargs = builder.Update("users", mutates, where.Eq("id", 1))
	assert.Regexp(t, fmt.Sprint("UPDATE `users` SET `", `\w*`, "`=", `\?`, ",`", `\w*`, "`=", `\?`, ",`", `\w*`, "`=", `\?`, " WHERE `id`=", `\?`, ";"), qs)
	assert.ElementsMatch(t, []interface{}{"foo", 10, true, 1}, qargs)
}

func TestBuilder_Update_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder  = NewBuilder(config)
		mutates = map[string]rel.Mutate{
			"name":  rel.Set("name", "foo"),
			"age":   rel.Set("age", 10),
			"agree": rel.Set("agree", true),
		}
	)

	qs, args := builder.Update("users", mutates, where.And())
	assert.Regexp(t, `UPDATE "users" SET "\w*"=\$1,"\w*"=\$2,"\w*"=\$3;`, qs)
	assert.ElementsMatch(t, []interface{}{"foo", 10, true}, args)

	builder.count = 0
	qs, args = builder.Update("users", mutates, where.Eq("id", 1))
	assert.Regexp(t, `UPDATE "users" SET "\w*"=\$1,"\w*"=\$2,"\w*"=\$3 WHERE "id"=\$4;`, qs)
	assert.ElementsMatch(t, []interface{}{"foo", 10, true, 1}, args)
}

func TestBuilder_Update_incDecAndFragment(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	qs, qargs := builder.Update("users", map[string]rel.Mutate{"age": rel.Inc("age")}, where.And())
	assert.Equal(t, "UPDATE `users` SET `age`=`age`+?;", qs)
	assert.Equal(t, []interface{}{1}, qargs)

	qs, qargs = builder.Update("users", map[string]rel.Mutate{"age=?": rel.SetFragment("age=?", 10)}, where.And())
	assert.Equal(t, "UPDATE `users` SET age=?;", qs)
	assert.Equal(t, []interface{}{10}, qargs)
}

func TestBuilder_Delete(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	qs, args := builder.Delete("users", where.And())
	assert.Equal(t, "DELETE FROM `users`;", qs)
	assert.Equal(t, []interface{}(nil), args)

	qs, args = builder.Delete("users", where.Eq("id", 1))
	assert.Equal(t, "DELETE FROM `users` WHERE `id`=?;", qs)
	assert.Equal(t, []interface{}{1}, args)
}

func TestBuilder_Delete_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder = NewBuilder(config)
	)

	qs, args := builder.Delete("users", where.And())
	assert.Equal(t, "DELETE FROM \"users\";", qs)
	assert.Equal(t, []interface{}(nil), args)

	qs, args = builder.Delete("users", where.Eq("id", 1))
	assert.Equal(t, "DELETE FROM \"users\" WHERE \"id\"=$1;", qs)
	assert.Equal(t, []interface{}{1}, args)
}

func TestBuilder_Select(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	tests := []struct {
		result   string
		distinct bool
		fields   []string
	}{
		{
			result: "SELECT *",
		},
		{
			result: "SELECT *",
			fields: []string{"*"},
		},
		{
			result: "SELECT `id`,`name`",
			fields: []string{"id", "name"},
		},
		{
			result:   "SELECT DISTINCT *",
			distinct: true,
			fields:   []string{"*"},
		},
		{
			result:   "SELECT DISTINCT `id`,`name`",
			distinct: true,
			fields:   []string{"id", "name"},
		},
		{
			result: "SELECT COUNT(*) AS count",
			fields: []string{"COUNT(*) AS count"},
		},
		{
			result: "SELECT COUNT(`transactions`.*) AS count",
			fields: []string{"COUNT(transactions.*) AS count"},
		},
		{
			result: "SELECT SUM(`transactions`.`total`) AS total",
			fields: []string{"SUM(transactions.total) AS total"},
		},
	}

	for _, test := range tests {
		t.Run(test.result, func(t *testing.T) {
			var (
				buffer Buffer
			)

			builder.fields(&buffer, test.distinct, test.fields)
			assert.Equal(t, test.result, buffer.String())
		})
	}
}

func TestBuilder_From(t *testing.T) {
	var (
		buffer Buffer
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	builder.from(&buffer, "users")
	assert.Equal(t, " FROM `users`", buffer.String())
}

func TestBuilder_Join(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
	)

	tests := []struct {
		QueryString string
		Query       rel.Query
	}{
		{
			"",
			rel.From("transactions"),
		},
		{
			" JOIN `users` ON `transactions`.`user_id`=`users`.`id`",
			rel.From("transactions").Join("users"),
		},
		{
			" JOIN `users` ON `users`.`id`=`transactions`.`user_id`",
			rel.From("transactions").JoinOn("users", "users.id", "transactions.user_id"),
		},
		{
			" INNER JOIN `users` ON `users`.`id`=`transactions`.`user_id`",
			rel.From("transactions").JoinWith("INNER JOIN", "users", "users.id", "transactions.user_id"),
		},
		{
			" JOIN `users` ON `users`.`id`=`transactions`.`user_id` JOIN `payments` ON `payments`.`id`=`transactions`.`payment_id`",
			rel.From("transactions").JoinOn("users", "users.id", "transactions.user_id").
				JoinOn("payments", "payments.id", "transactions.payment_id"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.join(&buffer, "transactions", rel.Build("", test.Query).JoinQuery)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Nil(t, buffer.Arguments)
		})
	}
}

func TestBuilder_Where(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			" WHERE `field`=?",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			" WHERE (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.where(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_Where_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			" WHERE \"field\"=$1",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			" WHERE (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.where(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_GroupBy(t *testing.T) {
	var (
		buffer Buffer
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	builder.groupBy(&buffer, []string{"city"})
	assert.Equal(t, " GROUP BY `city`", buffer.String())

	buffer.Reset()
	builder.groupBy(&buffer, []string{"city", "nation"})
	assert.Equal(t, " GROUP BY `city`,`nation`", buffer.String())
}

func TestBuilder_Having(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			" HAVING `field`=?",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			" HAVING (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.having(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_Having_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			" HAVING \"field\"=$1",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			" HAVING (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.having(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_OrderBy(t *testing.T) {
	var (
		buffer Buffer
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	builder.orderBy(&buffer, []rel.SortQuery{sort.Asc("name")})
	assert.Equal(t, " ORDER BY `name` ASC", buffer.String())

	buffer.Reset()
	builder.orderBy(&buffer, []rel.SortQuery{sort.Asc("name"), sort.Desc("created_at")})
	assert.Equal(t, " ORDER BY `name` ASC, `created_at` DESC", buffer.String())
}

func TestBuilder_LimitOffset(t *testing.T) {
	var (
		buffer Buffer
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	builder.limitOffset(&buffer, 10, 0)
	assert.Equal(t, " LIMIT 10", buffer.String())

	buffer.Reset()
	builder.limitOffset(&buffer, 10, 10)
	assert.Equal(t, " LIMIT 10 OFFSET 10", buffer.String())
}

func TestBuilder_Filter(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"`field`=?",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"`field`<>?",
			[]interface{}{"value"},
			where.Ne("field", "value"),
		},
		{
			"`field`<?",
			[]interface{}{10},
			where.Lt("field", 10),
		},
		{
			"`field`<=?",
			[]interface{}{10},
			where.Lte("field", 10),
		},
		{
			"`field`>?",
			[]interface{}{10},
			where.Gt("field", 10),
		},
		{
			"`field`>=?",
			[]interface{}{10},
			where.Gte("field", 10),
		},
		{
			"`field` IS NULL",
			nil,
			where.Nil("field"),
		},
		{
			"`field` IS NOT NULL",
			nil,
			where.NotNil("field"),
		},
		{
			"`field` IN (?)",
			[]interface{}{"value1"},
			where.In("field", "value1"),
		},
		{
			"`field` IN (?,?)",
			[]interface{}{"value1", "value2"},
			where.In("field", "value1", "value2"),
		},
		{
			"`field` IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			where.In("field", "value1", "value2", "value3"),
		},
		{
			"`field` NOT IN (?)",
			[]interface{}{"value1"},
			where.Nin("field", "value1"),
		},
		{
			"`field` NOT IN (?,?)",
			[]interface{}{"value1", "value2"},
			where.Nin("field", "value1", "value2"),
		},
		{
			"`field` NOT IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			where.Nin("field", "value1", "value2", "value3"),
		},
		{
			"`field` LIKE ?",
			[]interface{}{"%value%"},
			where.Like("field", "%value%"),
		},
		{
			"`field` NOT LIKE ?",
			[]interface{}{"%value%"},
			where.NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			where.Fragment("FRAGMENT"),
		},
		{
			"(`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
		{
			"(`field1`=? AND `field2`=? AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			where.Eq("field1", "value1").AndEq("field2", "value2").AndEq("field3", "value3"),
		},
		{
			"(`field1`=? OR `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").OrEq("field2", "value2"),
		},
		{
			"(`field1`=? OR `field2`=? OR `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			where.Eq("field1", "value1").OrEq("field2", "value2").OrEq("field3", "value3"),
		},
		{
			"NOT (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2")),
		},
		{
			"NOT (`field1`=? AND `field2`=? AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2"), where.Eq("field3", "value3")),
		},
		{
			"((`field1`=? OR `field2`=?) AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			where.And(where.Or(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Eq("field3", "value3")),
		},
		{
			"((`field1`=? OR `field2`=?) AND (`field3`=? OR `field4`=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			where.And(where.Or(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Or(where.Eq("field3", "value3"), where.Eq("field4", "value4"))),
		},
		{
			"(NOT (`field1`=? AND `field2`=?) AND NOT (`field3`=? OR `field4`=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			where.And(where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Not(where.Or(where.Eq("field3", "value3"), where.Eq("field4", "value4")))),
		},
		{
			"NOT (`field1`=? AND (`field2`=? OR `field3`=?) AND NOT (`field4`=? OR `field5`=?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			where.And(where.Not(where.Eq("field1", "value1"), where.Or(where.Eq("field2", "value2"), where.Eq("field3", "value3")), where.Not(where.Or(where.Eq("field4", "value4"), where.Eq("field5", "value5"))))),
		},
		{
			"((`field1` IN (?,?) OR `field2` NOT IN (?)) AND `field3` IN (?,?,?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			where.And(where.Or(where.In("field1", "value1", "value2"), where.Nin("field2", "value3")), where.In("field3", "value4", "value5", "value6")),
		},
		{
			"(`field1` LIKE ? AND `field2` NOT LIKE ?)",
			[]interface{}{"%value1%", "%value2%"},
			where.And(where.Like("field1", "%value1%"), where.NotLike("field2", "%value2%")),
		},
		{
			"",
			nil,
			rel.FilterQuery{Type: rel.FilterOp(9999)},
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.filter(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_Filter_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Filter      rel.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"\"field\"=$1",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"\"field\"<>$1",
			[]interface{}{"value"},
			where.Ne("field", "value"),
		},
		{
			"\"field\"<$1",
			[]interface{}{10},
			where.Lt("field", 10),
		},
		{
			"\"field\"<=$1",
			[]interface{}{10},
			where.Lte("field", 10),
		},
		{
			"\"field\">$1",
			[]interface{}{10},
			where.Gt("field", 10),
		},
		{
			"\"field\">=$1",
			[]interface{}{10},
			where.Gte("field", 10),
		},
		{
			"\"field\" IS NULL",
			nil,
			where.Nil("field"),
		},
		{
			"\"field\" IS NOT NULL",
			nil,
			where.NotNil("field"),
		},
		{
			"\"field\" IN ($1)",
			[]interface{}{"value1"},
			where.In("field", "value1"),
		},
		{
			"\"field\" IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			where.In("field", "value1", "value2"),
		},
		{
			"\"field\" IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.In("field", "value1", "value2", "value3"),
		},
		{
			"\"field\" NOT IN ($1)",
			[]interface{}{"value1"},
			where.Nin("field", "value1"),
		},
		{
			"\"field\" NOT IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			where.Nin("field", "value1", "value2"),
		},
		{
			"\"field\" NOT IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.Nin("field", "value1", "value2", "value3"),
		},
		{
			"\"field\" LIKE $1",
			[]interface{}{"%value%"},
			where.Like("field", "%value%"),
		},
		{
			"\"field\" NOT LIKE $1",
			[]interface{}{"%value%"},
			where.NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			where.Fragment("FRAGMENT"),
		},
		{
			"(\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
		{
			"(\"field1\"=$1 AND \"field2\"=$2 AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.Eq("field1", "value1").AndEq("field2", "value2").AndEq("field3", "value3"),
		},
		{
			"(\"field1\"=$1 OR \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").OrEq("field2", "value2"),
		},
		{
			"(\"field1\"=$1 OR \"field2\"=$2 OR \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.Eq("field1", "value1").OrEq("field2", "value2").OrEq("field3", "value3"),
		},
		{
			"NOT (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2")),
		},
		{
			"NOT (\"field1\"=$1 AND \"field2\"=$2 AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2"), where.Eq("field3", "value3")),
		},
		{
			"((\"field1\"=$1 OR \"field2\"=$2) AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			where.And(where.Or(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Eq("field3", "value3")),
		},
		{
			"((\"field1\"=$1 OR \"field2\"=$2) AND (\"field3\"=$3 OR \"field4\"=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			where.And(where.Or(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Or(where.Eq("field3", "value3"), where.Eq("field4", "value4"))),
		},
		{
			"(NOT (\"field1\"=$1 AND \"field2\"=$2) AND NOT (\"field3\"=$3 OR \"field4\"=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			where.And(where.Not(where.Eq("field1", "value1"), where.Eq("field2", "value2")), where.Not(where.Or(where.Eq("field3", "value3"), where.Eq("field4", "value4")))),
		},
		{
			"NOT (\"field1\"=$1 AND (\"field2\"=$2 OR \"field3\"=$3) AND NOT (\"field4\"=$4 OR \"field5\"=$5))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			where.And(where.Not(where.Eq("field1", "value1"), where.Or(where.Eq("field2", "value2"), where.Eq("field3", "value3")), where.Not(where.Or(where.Eq("field4", "value4"), where.Eq("field5", "value5"))))),
		},
		{
			"((\"field1\" IN ($1,$2) OR \"field2\" NOT IN ($3)) AND \"field3\" IN ($4,$5,$6))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			where.And(where.Or(where.In("field1", "value1", "value2"), where.Nin("field2", "value3")), where.In("field3", "value4", "value5", "value6")),
		},
		{
			"(\"field1\" LIKE $1 AND \"field2\" NOT LIKE $2)",
			[]interface{}{"%value1%", "%value2%"},
			where.And(where.Like("field1", "%value1%"), where.NotLike("field2", "%value2%")),
		},
		{
			"",
			nil,
			rel.FilterQuery{Type: rel.FilterOp(9999)},
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  Buffer
				builder = NewBuilder(config)
			)

			builder.filter(&buffer, test.Filter)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Equal(t, test.Args, buffer.Arguments)
		})
	}
}

func TestBuilder_Lock(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder  = NewBuilder(config)
		query    = rel.From("users").Lock("FOR UPDATE")
		qs, args = builder.Find(query)
	)

	assert.Equal(t, "SELECT * FROM `users` FOR UPDATE;", qs)
	assert.Nil(t, args)
}
