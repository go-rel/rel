package sql

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/sort"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Find(t *testing.T) {
	users := grimoire.From("users")

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
	}{
		{
			"SELECT * FROM `users`;",
			nil,
			users,
		},
		{
			"SELECT `users`.* FROM `users`;",
			nil,
			users.Select("users.*"),
		},
		{
			"SELECT `id`,`name` FROM `users`;",
			nil,
			users.Select("id", "name"),
		},
		{
			"SELECT `id`,FIELD(`gender`, \"male\") AS `order` FROM `users` ORDER BY `order` ASC;",
			nil,
			users.Select("id", "^FIELD(`gender`, \"male\") AS `order`").SortAsc("order"),
		},
		{
			"SELECT * FROM `users` JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			users.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM `users` WHERE `id`=?;",
			[]interface{}{10},
			users.Where(where.Eq("id", 10)),
		},
		{
			"SELECT DISTINCT * FROM `users` GROUP BY `type` HAVING `price`>?;",
			[]interface{}{1000},
			users.Distinct().Group("type").Having(where.Gt("price", 1000)),
		},
		{
			"SELECT * FROM `users` INNER JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			users.JoinWith("INNER JOIN", "transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC;",
			nil,
			users.SortAsc("created_at"),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC, `id` DESC;",
			nil,
			users.SortAsc("created_at").SortDesc("id"),
		},
		{
			"SELECT * FROM `users` LIMIT 10 OFFSET 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).Find(test.Query)
			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Find_ordinal(t *testing.T) {
	users := grimoire.From("users")

	tests := []struct {
		QueryString string
		Args        []interface{}
		Query       grimoire.Query
	}{
		{
			"SELECT * FROM \"users\";",
			nil,
			users,
		},
		{
			"SELECT \"users\".* FROM \"users\";",
			nil,
			users.Select("users.*"),
		},
		{
			"SELECT \"id\",\"name\" FROM \"users\";",
			nil,
			users.Select("id", "name"),
		},
		{
			"SELECT * FROM \"users\" JOIN \"transactions\" ON \"transactions\".\"id\"=\"users\".\"transaction_id\";",
			nil,
			users.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM \"users\" WHERE \"id\"=$1;",
			[]interface{}{10},
			users.Where(where.Eq("id", 10)),
		},
		{
			"SELECT DISTINCT * FROM \"users\" GROUP BY \"type\" HAVING \"price\">$1;",
			[]interface{}{1000},
			users.Distinct().Group("type").Having(where.Gt("price", 1000)),
		},
		{
			"SELECT * FROM \"users\" JOIN \"transactions\" ON \"transactions\".\"id\"=\"users\".\"transaction_id\";",
			nil,
			users.JoinOn("transactions", "transactions.id", "users.transaction_id"),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC;",
			nil,
			users.SortAsc("created_at"),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC, \"id\" DESC;",
			nil,
			users.SortAsc("created_at").SortDesc("id"),
		},
		{
			"SELECT * FROM \"users\" LIMIT 10 OFFSET 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).Find(test.Query)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

// func TestBuilder_Aggregate(t *testing.T) {
// 	builder := NewBuilder(&Config{
// 		Placeholder: "?",
// 		EscapeChar:  "`",
// 	})

// 	users := grimoire.From("users")

// 	users.AggregateMode = "count"
// 	users.AggregateField = "*"

// 	qs, args := builder.Aggregate(users)
// 	assert.Nil(t, args)
// 	assert.Equal(t, "SELECT count(*) AS count FROM `users`;", qs)

// 	users.AggregateMode = "sum"
// 	users.AggregateField = "transactions.total"

// 	qs, args = builder.Aggregate(users)
// 	assert.Nil(t, args)
// 	assert.Equal(t, "SELECT sum(`transactions`.`total`) AS sum FROM `users`;", qs)

// 	qs, args = builder.Aggregate(users.Group("gender"))
// 	assert.Nil(t, args)
// 	assert.Equal(t, "SELECT `gender`,sum(`transactions`.`total`) AS sum FROM `users` GROUP BY `gender`;", qs)
// }

func TestBuilder_Insert(t *testing.T) {
	var (
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
		args = []interface{}{"foo", 10, true}
	)

	qs, qargs := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Insert("users", changes)

	assert.Equal(t, "INSERT INTO `users` (`name`,`age`,`agree`) VALUES (?,?,?);", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Returning("id").Insert("users", changes)

	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\",\"agree\") VALUES ($1,$2,$3) RETURNING \"id\";", qs)
	assert.Equal(t, args, qargs)
}

func TestBuilder_Insert_defaultValues(t *testing.T) {
	var (
		changes = grimoire.Changes{}
		args    = []interface{}{}
	)

	qs, qargs := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Insert("users", changes)

	assert.Equal(t, "INSERT INTO `users` () VALUES ();", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "?",
		InsertDefaultValues: true,
		EscapeChar:          "`",
	}).Returning("id").Insert("users", changes)

	assert.Equal(t, "INSERT INTO `users` DEFAULT VALUES RETURNING `id`;", qs)
	assert.Equal(t, args, qargs)
}

// func TestBuilder_InsertAll(t *testing.T) {
// 	fields := []string{"name"}
// 	allchanges := []map[string]interface{}{
// 		{"name": "foo"},
// 		{"age": 12},
// 		{"name": "boo"},
// 	}

// 	statement, args := NewBuilder(&Config{
// 		Placeholder: "?",
// 		EscapeChar:  "`",
// 	}).InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO `users` (`name`) VALUES (?),(DEFAULT),(?);", statement)
// 	assert.Equal(t, []interface{}{"foo", "boo"}, args)

// 	// ordinal
// 	statement, args = NewBuilder(&Config{
// 		Placeholder:         "$",
// 		EscapeChar:          "\"",
// 		Ordinal:             true,
// 		InsertDefaultValues: true,
// 	}).Returning("id").InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO \"users\" (\"name\") VALUES ($1),(DEFAULT),($2) RETURNING \"id\";", statement)
// 	assert.Equal(t, []interface{}{"foo", "boo"}, args)

// 	// with age
// 	fields = append(fields, "age")
// 	statement, args = NewBuilder(&Config{
// 		Placeholder: "?",
// 		EscapeChar:  "`",
// 	}).InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO `users` (`name`,`age`) VALUES (?,DEFAULT),(DEFAULT,?),(?,DEFAULT);", statement)
// 	assert.Equal(t, []interface{}{"foo", 12, "boo"}, args)

// 	// ordinal
// 	statement, args = NewBuilder(&Config{
// 		Placeholder:         "$",
// 		EscapeChar:          "\"",
// 		Ordinal:             true,
// 		InsertDefaultValues: true,
// 	}).Returning("id").InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,DEFAULT),(DEFAULT,$2),($3,DEFAULT) RETURNING \"id\";", statement)
// 	assert.Equal(t, []interface{}{"foo", 12, "boo"}, args)

// 	// all changes have value
// 	allchanges = []map[string]interface{}{
// 		{"name": "foo", "age": 10},
// 		{"name": "zoo", "age": 12},
// 		{"name": "boo", "age": 20},
// 	}

// 	statement, args = NewBuilder(&Config{
// 		Placeholder: "?",
// 		EscapeChar:  "`",
// 	}).InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO `users` (`name`,`age`) VALUES (?,?),(?,?),(?,?);", statement)
// 	assert.Equal(t, []interface{}{"foo", 10, "zoo", 12, "boo", 20}, args)

// 	// ordinal
// 	statement, args = NewBuilder(&Config{
// 		Placeholder:         "$",
// 		EscapeChar:          "\"",
// 		Ordinal:             true,
// 		InsertDefaultValues: true,
// 	}).InsertAll("users", fields, allchanges)
// 	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,$2),($3,$4),($5,$6);", statement)
// 	assert.Equal(t, []interface{}{"foo", 10, "zoo", 12, "boo", 20}, args)
// }

func TestBuilder_Update(t *testing.T) {
	var (
		changes = grimoire.BuildChanges(
			grimoire.Set("name", "foo"),
			grimoire.Set("age", 10),
			grimoire.Set("agree", true),
		)
		args = []interface{}{"foo", 10, true}
		cond = where.And()
	)

	qs, qargs := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Update("users", changes, cond)

	assert.Equal(t, "UPDATE `users` SET `name`=?,`age`=?,`agree`=?;", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Update("users", changes, cond)

	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1,\"age\"=$2,\"agree\"=$3;", qs)
	assert.Equal(t, args, qargs)

	args = []interface{}{"foo", 10, true, 1}
	cond = where.Eq("id", 1)

	qs, qargs = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Update("users", changes, cond)

	assert.Equal(t, "UPDATE `users` SET `name`=?,`age`=?,`agree`=? WHERE `id`=?;", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Update("users", changes, cond)
	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1,\"age\"=$2,\"agree\"=$3 WHERE \"id\"=$4;", qs)
	assert.Equal(t, args, qargs)
}

func TestBuilder_Delete(t *testing.T) {
	qs, args := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Delete("users", where.And())
	assert.Equal(t, "DELETE FROM `users`;", qs)
	assert.Equal(t, []interface{}(nil), args)

	qs, args = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Delete("users", where.Eq("id", 1))
	assert.Equal(t, "DELETE FROM `users` WHERE `id`=?;", qs)
	assert.Equal(t, []interface{}{1}, args)

	qs, args = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Delete("users", where.Eq("id", 1))
	assert.Equal(t, "DELETE FROM \"users\" WHERE \"id\"=$1;", qs)
	assert.Equal(t, []interface{}{1}, args)
}

func TestBuilder_Select(t *testing.T) {
	assert.Equal(t, "SELECT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, nil))

	assert.Equal(t, "SELECT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, []string{"*"}))

	assert.Equal(t, "SELECT `id`,`name`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, []string{"id", "name"}))

	assert.Equal(t, "SELECT DISTINCT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(true, []string{"*"}))

	assert.Equal(t, "SELECT DISTINCT `id`,`name`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(true, []string{"id", "name"}))

	assert.Equal(t, "SELECT COUNT(*) AS count", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, []string{"COUNT(*) AS count"}))

	assert.Equal(t, "SELECT COUNT(`transactions`.*) AS count", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, []string{"COUNT(transactions.*) AS count"}))

	assert.Equal(t, "SELECT SUM(`transactions`.`total`) AS total", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, []string{"SUM(transactions.total) AS total"}))
}

func TestBuilder_From(t *testing.T) {
	assert.Equal(t, "FROM `users`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).from("users"))
}

func TestBuilder_Join(t *testing.T) {
	tests := []struct {
		QueryString string
		Query       grimoire.Query
	}{
		{
			"",
			grimoire.From("trxs"),
		},
		{
			"JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			grimoire.From("trxs").JoinOn("users", "user.id", "trxs.user_id"),
		},
		{
			"INNER JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			grimoire.From("trxs").JoinWith("INNER JOIN", "users", "user.id", "trxs.user_id"),
		},
		{
			"JOIN `users` ON `user`.`id`=`trxs`.`user_id` JOIN `payments` ON `payments`.`id`=`trxs`.`payment_id`",
			grimoire.From("trxs").JoinOn("users", "user.id", "trxs.user_id").
				JoinOn("payments", "payments.id", "trxs.payment_id"),
		},
	}

	for _, test := range tests {
		t.Run(test.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).join(grimoire.BuildQuery("", test.Query).JoinQuery...)
			assert.Equal(t, test.QueryString, qs)
			assert.Nil(t, args)
		})
	}
}

func TestBuilder_Where(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).where(test.Filter)
			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Where_ordinal(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).where(test.Filter)
			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_GroupBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy(nil))

	assert.Equal(t, "GROUP BY `city`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy([]string{"city"}))

	assert.Equal(t, "GROUP BY `city`,`nation`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy([]string{"city", "nation"}))
}

func TestBuilder_Having(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).having(test.Filter)
			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Having_ordinal(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).having(test.Filter)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_OrderBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy(nil))

	assert.Equal(t, "ORDER BY `name` ASC", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy([]grimoire.SortQuery{sort.Asc("name")}))

	assert.Equal(t, "ORDER BY `name` ASC, `created_at` DESC", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy([]grimoire.SortQuery{sort.Asc("name"), sort.Desc("created_at")}))
}

func TestBuilder_LimitOffset(t *testing.T) {
	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).limitOffset(0, 0))

	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).limitOffset(0, 10))

	assert.Equal(t, "LIMIT 10", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).limitOffset(10, 0))

	assert.Equal(t, "LIMIT 10 OFFSET 10", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).limitOffset(10, 10))
}

func TestBuilder_Filter(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).filter(test.Filter)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Filter_ordinal(t *testing.T) {
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
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).filter(test.Filter)

			assert.Equal(t, test.QueryString, qs)
			assert.Equal(t, test.Args, args)
		})
	}
}

func TestBuilder_Lock(t *testing.T) {
	users := grimoire.From("users").Lock(grimoire.ForUpdate())

	qs, args := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Find(users)

	assert.Equal(t, "SELECT * FROM `users` FOR UPDATE;", qs)
	assert.Nil(t, args)
}
