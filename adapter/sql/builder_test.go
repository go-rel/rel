package sql

import (
	"bytes"
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/sort"
	"github.com/Fs02/grimoire/where"
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
		query := grimoire.From("users").
			Select("id", "name").
			Join("transactions").
			Where(where.Eq("id", 10)).
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
		query = grimoire.From("users")
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
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
		query = grimoire.From("users")
	)

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
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

func TestBuilder_Aggregate(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
		query   = grimoire.From("users")
	)

	qs, args := builder.Aggregate(query, "count", "*")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT count(*) AS count FROM `users`;", qs)

	qs, args = builder.Aggregate(query, "sum", "transactions.total")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT sum(`transactions`.`total`) AS sum FROM `users`;", qs)

	qs, args = builder.Aggregate(query.Group("gender"), "sum", "transactions.total")
	assert.Nil(t, args)
	assert.Equal(t, "SELECT `gender`,sum(`transactions`.`total`) AS sum FROM `users` GROUP BY `gender`;", qs)
}

func TestBuilder_Insert(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
		qs, args = builder.Insert("users", changes)
	)

	assert.Equal(t, "INSERT INTO `users` (`name`,`age`,`agree`) VALUES (?,?,?);", qs)
	assert.Equal(t, []interface{}{"foo", 10, true}, args)
}

func TestBuilder_Insert_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder = NewBuilder(config)
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
		qs, args = builder.Returning("id").Insert("users", changes)
	)

	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\",\"agree\") VALUES ($1,$2,$3) RETURNING \"id\";", qs)
	assert.Equal(t, []interface{}{"foo", 10, true}, args)
}

func TestBuilder_Insert_defaultValuesDisabled(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "?",
			EscapeChar:          "`",
			InsertDefaultValues: false,
		}
		builder  = NewBuilder(config)
		changes  = grimoire.Changes{}
		qs, args = builder.Insert("users", changes)
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
		changes  = grimoire.Changes{}
		qs, args = builder.Returning("id").Insert("users", changes)
	)

	assert.Equal(t, "INSERT INTO `users` DEFAULT VALUES RETURNING `id`;", qs)
	assert.Equal(t, []interface{}{}, args)
}

func TestBuilder_InsertAll(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder    = NewBuilder(config)
		allchanges = []grimoire.Changes{
			grimoire.BuildChanges(
				grimoire.Set("name", "foo"),
			),
			grimoire.BuildChanges(
				grimoire.Set("age", 10),
			),
			grimoire.BuildChanges(
				grimoire.Set("name", "boo"),
				grimoire.Set("age", 20),
			),
		}
	)

	statement, args := builder.InsertAll("users", []string{"name"}, allchanges)
	assert.Equal(t, "INSERT INTO `users` (`name`) VALUES (?),(DEFAULT),(?);", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// with age
	statement, args = builder.InsertAll("users", []string{"name", "age"}, allchanges)
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
		builder    = NewBuilder(config)
		allchanges = []grimoire.Changes{
			grimoire.BuildChanges(
				grimoire.Set("name", "foo"),
			),
			grimoire.BuildChanges(
				grimoire.Set("age", 10),
			),
			grimoire.BuildChanges(
				grimoire.Set("name", "boo"),
				grimoire.Set("age", 20),
			),
		}
	)

	statement, args := builder.Returning("id").InsertAll("users", []string{"name"}, allchanges)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\") VALUES ($1),(DEFAULT),($2) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// with age
	builder.count = 0
	statement, args = builder.Returning("id").InsertAll("users", []string{"name", "age"}, allchanges)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,DEFAULT),(DEFAULT,$2),($3,$4) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", 10, "boo", 20}, args)
}

func TestBuilder_Update(t *testing.T) {
	var (
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
	)

	qs, qargs := builder.Update("users", changes, where.And())
	assert.Equal(t, "UPDATE `users` SET `name`=?,`age`=?,`agree`=?;", qs)
	assert.Equal(t, []interface{}{"foo", 10, true}, qargs)

	qs, qargs = builder.Update("users", changes, where.Eq("id", 1))
	assert.Equal(t, "UPDATE `users` SET `name`=?,`age`=?,`agree`=? WHERE `id`=?;", qs)
	assert.Equal(t, []interface{}{"foo", 10, true, 1}, qargs)
}

func TestBuilder_Update_ordinal(t *testing.T) {
	var (
		config = &Config{
			Placeholder:         "$",
			EscapeChar:          "\"",
			Ordinal:             true,
			InsertDefaultValues: true,
		}
		builder = NewBuilder(config)
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
	)

	qs, args := builder.Update("users", changes, where.And())
	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1,\"age\"=$2,\"agree\"=$3;", qs)
	assert.Equal(t, []interface{}{"foo", 10, true}, args)

	builder.count = 0
	qs, args = builder.Update("users", changes, where.Eq("id", 1))
	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1,\"age\"=$2,\"agree\"=$3 WHERE \"id\"=$4;", qs)
	assert.Equal(t, []interface{}{"foo", 10, true, 1}, args)
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
				buffer bytes.Buffer
			)

			builder.fields(&buffer, test.distinct, test.fields)
			assert.Equal(t, test.result, buffer.String())
		})
	}
}

func TestBuilder_From(t *testing.T) {
	var (
		buffer bytes.Buffer
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
		Query       grimoire.Query
	}{
		{
			"",
			grimoire.From("trxs"),
		},
		{
			" JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			grimoire.From("trxs").JoinOn("users", "user.id", "trxs.user_id"),
		},
		{
			" INNER JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			grimoire.From("trxs").JoinWith("INNER JOIN", "users", "user.id", "trxs.user_id"),
		},
		{
			" JOIN `users` ON `user`.`id`=`trxs`.`user_id` JOIN `payments` ON `payments`.`id`=`trxs`.`payment_id`",
			grimoire.From("trxs").JoinOn("users", "user.id", "trxs.user_id").
				JoinOn("payments", "payments.id", "trxs.payment_id"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				buffer  bytes.Buffer
				builder = NewBuilder(config)
				args    = builder.join(&buffer, grimoire.BuildQuery("", test.Query).JoinQuery...)
			)

			assert.Equal(t, test.QueryString, buffer.String())
			assert.Nil(t, args)
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
		Filter      grimoire.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"WHERE `field`=?",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"WHERE (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.where(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
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
		Filter      grimoire.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"WHERE \"field\"=$1",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"WHERE (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.where(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_GroupBy(t *testing.T) {
	var (
		buffer bytes.Buffer
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
		Filter      grimoire.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"HAVING `field`=?",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"HAVING (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.having(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
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
		Filter      grimoire.FilterQuery
	}{
		{
			"",
			nil,
			where.And(),
		},
		{
			"HAVING \"field\"=$1",
			[]interface{}{"value"},
			where.Eq("field", "value"),
		},
		{
			"HAVING (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			where.Eq("field1", "value1").AndEq("field2", "value2"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.having(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_OrderBy(t *testing.T) {
	var (
		buffer bytes.Buffer
		config = &Config{
			Placeholder: "?",
			EscapeChar:  "`",
		}
		builder = NewBuilder(config)
	)

	builder.orderBy(&buffer, []grimoire.SortQuery{sort.Asc("name")})
	assert.Equal(t, " ORDER BY `name` ASC", buffer.String())

	buffer.Reset()
	builder.orderBy(&buffer, []grimoire.SortQuery{sort.Asc("name"), sort.Desc("created_at")})
	assert.Equal(t, " ORDER BY `name` ASC, `created_at` DESC", buffer.String())
}

func TestBuilder_LimitOffset(t *testing.T) {
	var (
		buffer bytes.Buffer
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
		Filter      grimoire.FilterQuery
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
			grimoire.FilterQuery{Type: grimoire.FilterOp(9999)},
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.filter(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
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
		Filter      grimoire.FilterQuery
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
			grimoire.FilterQuery{Type: grimoire.FilterOp(9999)},
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			var (
				builder  = NewBuilder(config)
				qs, args = builder.filter(test.Filter)
			)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
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
		query    = grimoire.From("users").Lock(grimoire.ForUpdate())
		qs, args = builder.Find(query)
	)

	assert.Equal(t, "SELECT * FROM `users` FOR UPDATE;", qs)
	assert.Nil(t, args)
}
