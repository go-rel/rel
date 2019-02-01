package sql

import (
	"strings"
	"testing"

	"github.com/Fs02/grimoire"
	. "github.com/Fs02/grimoire/c"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Find(t *testing.T) {
	users := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
	}

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
			"SELECT * FROM `users` JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM `users` WHERE `id`=?;",
			[]interface{}{10},
			users.Where(Eq(I("id"), 10)),
		},
		{
			"SELECT DISTINCT * FROM `users` GROUP BY `type` HAVING `price`>?;",
			[]interface{}{1000},
			users.Distinct().Group("type").Having(Gt(I("price"), 1000)),
		},
		{
			"SELECT * FROM `users` JOIN `transactions` ON `transactions`.`id`=`users`.`transaction_id`;",
			nil,
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC;",
			nil,
			users.Order(Asc("created_at")),
		},
		{
			"SELECT * FROM `users` ORDER BY `created_at` ASC, `id` DESC;",
			nil,
			users.Order(Asc("created_at"), Desc("id")),
		},
		{
			"SELECT * FROM `users` LIMIT 10 OFFSET 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).Find(tt.Query)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Find_ordinal(t *testing.T) {
	users := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
	}

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
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM \"users\" WHERE \"id\"=$1;",
			[]interface{}{10},
			users.Where(Eq(I("id"), 10)),
		},
		{
			"SELECT DISTINCT * FROM \"users\" GROUP BY \"type\" HAVING \"price\">$1;",
			[]interface{}{1000},
			users.Distinct().Group("type").Having(Gt(I("price"), 1000)),
		},
		{
			"SELECT * FROM \"users\" JOIN \"transactions\" ON \"transactions\".\"id\"=\"users\".\"transaction_id\";",
			nil,
			users.Join("transactions", Eq(I("transactions.id"), I("users.transaction_id"))),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC;",
			nil,
			users.Order(Asc("created_at")),
		},
		{
			"SELECT * FROM \"users\" ORDER BY \"created_at\" ASC, \"id\" DESC;",
			nil,
			users.Order(Asc("created_at"), Desc("id")),
		},
		{
			"SELECT * FROM \"users\" LIMIT 10 OFFSET 10;",
			nil,
			users.Offset(10).Limit(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).Find(tt.Query)

			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Aggregate(t *testing.T) {
	builder := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	})

	users := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
	}

	users.AggregateMode = "count"
	users.AggregateField = "*"

	qs, args := builder.Aggregate(users)
	assert.Nil(t, args)
	assert.Equal(t, "SELECT count(*) AS count FROM `users`;", qs)

	users.AggregateMode = "sum"
	users.AggregateField = "transactions.total"

	qs, args = builder.Aggregate(users)
	assert.Nil(t, args)
	assert.Equal(t, "SELECT `transactions`.`total`,sum(`transactions`.`total`) AS sum FROM `users`;", qs)
}

func TestBuilder_Insert(t *testing.T) {
	changes := map[string]interface{}{
		"name": "foo",
	}
	args := []interface{}{"foo"}

	qs, qargs := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Insert("users", changes)
	assert.Equal(t, "INSERT INTO `users` (`name`) VALUES (?);", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Returning("id").Insert("users", changes)

	assert.Equal(t, "INSERT INTO \"users\" (\"name\") VALUES ($1) RETURNING \"id\";", qs)
	assert.Equal(t, args, qargs)

	// test for multiple changes since map is randomly ordered
	changes["age"] = 10
	changes["agree"] = true
	qs, _ = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Insert("users", changes)

	assert.True(t, strings.HasPrefix(qs, "INSERT INTO `users` ("))
	assert.True(t, strings.Contains(qs, "`name`"))
	assert.True(t, strings.Contains(qs, "`age`"))
	assert.True(t, strings.Contains(qs, "`agree`"))
	assert.True(t, strings.HasSuffix(qs, ";"))
}

func TestBuilder_Insert_defaultValues(t *testing.T) {
	changes := map[string]interface{}{}
	args := []interface{}{}

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

func TestBuilder_InsertAll(t *testing.T) {
	fields := []string{"name"}
	allchanges := []map[string]interface{}{
		{"name": "foo"},
		{"age": 12},
		{"name": "boo"},
	}

	statement, args := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO `users` (`name`) VALUES (?),(DEFAULT),(?);", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// ordinal
	statement, args = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Returning("id").InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\") VALUES ($1),(DEFAULT),($2) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", "boo"}, args)

	// with age
	fields = append(fields, "age")
	statement, args = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO `users` (`name`,`age`) VALUES (?,DEFAULT),(DEFAULT,?),(?,DEFAULT);", statement)
	assert.Equal(t, []interface{}{"foo", 12, "boo"}, args)

	// ordinal
	statement, args = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Returning("id").InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,DEFAULT),(DEFAULT,$2),($3,DEFAULT) RETURNING \"id\";", statement)
	assert.Equal(t, []interface{}{"foo", 12, "boo"}, args)

	// all changes have value
	allchanges = []map[string]interface{}{
		{"name": "foo", "age": 10},
		{"name": "zoo", "age": 12},
		{"name": "boo", "age": 20},
	}

	statement, args = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO `users` (`name`,`age`) VALUES (?,?),(?,?),(?,?);", statement)
	assert.Equal(t, []interface{}{"foo", 10, "zoo", 12, "boo", 20}, args)

	// ordinal
	statement, args = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).InsertAll("users", fields, allchanges)
	assert.Equal(t, "INSERT INTO \"users\" (\"name\",\"age\") VALUES ($1,$2),($3,$4),($5,$6);", statement)
	assert.Equal(t, []interface{}{"foo", 10, "zoo", 12, "boo", 20}, args)
}

func TestBuilder_Update(t *testing.T) {
	changes := map[string]interface{}{
		"name": "foo",
	}
	args := []interface{}{"foo"}
	cond := And()

	qs, qargs := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Update("users", changes, cond)
	assert.Equal(t, "UPDATE `users` SET `name`=?;", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Update("users", changes, cond)
	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1;", qs)
	assert.Equal(t, args, qargs)

	args = []interface{}{"foo", 1}
	cond = Eq(I("id"), 1)

	qs, qargs = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Update("users", changes, cond)
	assert.Equal(t, "UPDATE `users` SET `name`=? WHERE `id`=?;", qs)
	assert.Equal(t, args, qargs)

	qs, qargs = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Update("users", changes, cond)
	assert.Equal(t, "UPDATE \"users\" SET \"name\"=$1 WHERE \"id\"=$2;", qs)
	assert.Equal(t, args, qargs)

	// test for multiple changes since map is randomly ordered
	changes["age"] = 10
	changes["agree"] = true
	qs, _ = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Update("users", changes, And())

	assert.True(t, strings.HasPrefix(qs, "UPDATE \"users\" SET "))
	assert.True(t, strings.Contains(qs, "\"name\""))
	assert.True(t, strings.Contains(qs, "\"age\""))
	assert.True(t, strings.Contains(qs, "\"agree\""))
	assert.True(t, strings.HasSuffix(qs, ";"))
}

func TestBuilder_Delete(t *testing.T) {
	qs, args := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Delete("users", And())
	assert.Equal(t, "DELETE FROM `users`;", qs)
	assert.Equal(t, []interface{}(nil), args)

	qs, args = NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Delete("users", Eq(I("id"), 1))
	assert.Equal(t, "DELETE FROM `users` WHERE `id`=?;", qs)
	assert.Equal(t, []interface{}{1}, args)

	qs, args = NewBuilder(&Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
	}).Delete("users", Eq(I("id"), 1))
	assert.Equal(t, "DELETE FROM \"users\" WHERE \"id\"=$1;", qs)
	assert.Equal(t, []interface{}{1}, args)
}

func TestBuilder_Select(t *testing.T) {
	assert.Equal(t, "SELECT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false))

	assert.Equal(t, "SELECT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, "*"))

	assert.Equal(t, "SELECT `id`,`name`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, "id", "name"))

	assert.Equal(t, "SELECT DISTINCT *", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(true, "*"))

	assert.Equal(t, "SELECT DISTINCT `id`,`name`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(true, "id", "name"))

	assert.Equal(t, "SELECT COUNT(*) AS count", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, "COUNT(*) AS count"))

	assert.Equal(t, "SELECT COUNT(`transactions`.*) AS count", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, "COUNT(transactions.*) AS count"))

	assert.Equal(t, "SELECT SUM(`transactions`.`total`) AS total", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).fields(false, "SUM(transactions.total) AS total"))
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
		Args        []interface{}
		JoinClause  []Join
	}{
		{
			"",
			nil,
			nil,
		},
		{
			"JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			nil,
			grimoire.Query{Collection: "trxs"}.Join("users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"INNER JOIN `users` ON `user`.`id`=`trxs`.`user_id`",
			nil,
			grimoire.Query{Collection: "trxs"}.JoinWith("INNER JOIN", "users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"JOIN `users` ON `user`.`id`=`trxs`.`user_id` JOIN `payments` ON `payments`.`id`=`trxs`.`payment_id`",
			nil,
			grimoire.Query{Collection: "trxs"}.Join("users", Eq(I("user.id"), I("trxs.user_id"))).
				Join("payments", Eq(I("payments.id"), I("trxs.payment_id"))).JoinClause,
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).join(tt.JoinClause...)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Where(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"WHERE `field`=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"WHERE (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).where(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Where_ordinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"WHERE \"field\"=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"WHERE (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).where(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_GroupBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy())

	assert.Equal(t, "GROUP BY `city`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy("city"))

	assert.Equal(t, "GROUP BY `city`,`nation`", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).groupBy("city", "nation"))
}

func TestBuilder_Having(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"HAVING `field`=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"HAVING (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).having(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Having_ordinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"HAVING \"field\"=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"HAVING (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).having(tt.Condition)

			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_OrderBy(t *testing.T) {
	assert.Equal(t, "", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy())

	assert.Equal(t, "ORDER BY `name` ASC", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy(Asc("name")))

	assert.Equal(t, "ORDER BY `name` ASC, `created_at` DESC", NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).orderBy(Asc("name"), Desc("created_at")))
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

func TestBuilder_Condition(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"`field`=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"?=`field`",
			[]interface{}{"value"},
			Eq("value", I("field")),
		},
		{
			"?=?",
			[]interface{}{"value1", "value2"},
			Eq("value1", "value2"),
		},
		{
			"`field`<>?",
			[]interface{}{"value"},
			Ne(I("field"), "value"),
		},
		{
			"?<>`field`",
			[]interface{}{"value"},
			Ne("value", I("field")),
		},
		{
			"?<>?",
			[]interface{}{"value1", "value2"},
			Ne("value1", "value2"),
		},
		{
			"`field`<?",
			[]interface{}{10},
			Lt(I("field"), 10),
		},
		{
			"?<`field`",
			[]interface{}{"value"},
			Lt("value", I("field")),
		},
		{
			"?<?",
			[]interface{}{"value1", "value2"},
			Lt("value1", "value2"),
		},
		{
			"`field`<=?",
			[]interface{}{10},
			Lte(I("field"), 10),
		},
		{
			"?<=`field`",
			[]interface{}{"value"},
			Lte("value", I("field")),
		},
		{
			"?<=?",
			[]interface{}{"value1", "value2"},
			Lte("value1", "value2"),
		},
		{
			"`field`>?",
			[]interface{}{10},
			Gt(I("field"), 10),
		},
		{
			"?>`field`",
			[]interface{}{"value"},
			Gt("value", I("field")),
		},
		{
			"?>?",
			[]interface{}{"value1", "value2"},
			Gt("value1", "value2"),
		},
		{
			"`field`>=?",
			[]interface{}{10},
			Gte(I("field"), 10),
		},
		{
			"?>=`field`",
			[]interface{}{"value"},
			Gte("value", I("field")),
		},
		{
			"?>=?",
			[]interface{}{"value1", "value2"},
			Gte("value1", "value2"),
		},
		{
			"`field` IS NULL",
			nil,
			Nil("field"),
		},
		{
			"`field` IS NOT NULL",
			nil,
			NotNil("field"),
		},
		{
			"`field` IN (?)",
			[]interface{}{"value1"},
			In("field", "value1"),
		},
		{
			"`field` IN (?,?)",
			[]interface{}{"value1", "value2"},
			In("field", "value1", "value2"),
		},
		{
			"`field` IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			In("field", "value1", "value2", "value3"),
		},
		{
			"`field` NOT IN (?)",
			[]interface{}{"value1"},
			Nin("field", "value1"),
		},
		{
			"`field` NOT IN (?,?)",
			[]interface{}{"value1", "value2"},
			Nin("field", "value1", "value2"),
		},
		{
			"`field` NOT IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			Nin("field", "value1", "value2", "value3"),
		},
		{
			"`field` LIKE ?",
			[]interface{}{"%value%"},
			Like("field", "%value%"),
		},
		{
			"`field` NOT LIKE ?",
			[]interface{}{"%value%"},
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			Fragment("FRAGMENT"),
		},
		{
			"(`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(`field1`=? AND `field2`=? AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(`field1`=? OR `field2`=?)",
			[]interface{}{"value1", "value2"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(`field1`=? OR `field2`=? OR `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"NOT (`field1`=? AND `field2`=?)",
			[]interface{}{"value1", "value2"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"NOT (`field1`=? AND `field2`=? AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"((`field1`=? OR `field2`=?) AND `field3`=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Eq(I("field3"), "value3")),
		},
		{
			"((`field1`=? OR `field2`=?) AND (`field3`=? OR `field4`=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4"))),
		},
		{
			"(NOT (`field1`=? AND `field2`=?) AND NOT (`field3`=? OR `field4`=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Not(Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4")))),
		},
		{
			"NOT (`field1`=? AND (`field2`=? OR `field3`=?) AND NOT (`field4`=? OR `field5`=?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq(I("field1"), "value1"), Or(Eq(I("field2"), "value2"), Eq(I("field3"), "value3")), Not(Or(Eq(I("field4"), "value4"), Eq(I("field5"), "value5"))))),
		},
		{
			"((`field1` IN (?,?) OR `field2` NOT IN (?)) AND `field3` IN (?,?,?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(`field1` LIKE ? AND `field2` NOT LIKE ?)",
			[]interface{}{"%value1%", "%value2%"},
			And(Like(I("field1"), "%value1%"), NotLike(I(I("field2")), "%value2%")),
		},
		{
			"",
			nil,
			Condition{Type: ConditionType(9999)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder: "?",
				EscapeChar:  "`",
			}).condition(tt.Condition)

			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Condition_ordinal(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		Condition   Condition
	}{
		{
			"",
			nil,
			And(),
		},
		{
			"\"field\"=$1",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"$1=\"field\"",
			[]interface{}{"value"},
			Eq("value", I("field")),
		},
		{
			"$1=$2",
			[]interface{}{"value1", "value2"},
			Eq("value1", "value2"),
		},
		{
			"\"field\"<>$1",
			[]interface{}{"value"},
			Ne(I("field"), "value"),
		},
		{
			"$1<>\"field\"",
			[]interface{}{"value"},
			Ne("value", I("field")),
		},
		{
			"$1<>$2",
			[]interface{}{"value1", "value2"},
			Ne("value1", "value2"),
		},
		{
			"\"field\"<$1",
			[]interface{}{10},
			Lt(I("field"), 10),
		},
		{
			"$1<\"field\"",
			[]interface{}{"value"},
			Lt("value", I("field")),
		},
		{
			"$1<$2",
			[]interface{}{"value1", "value2"},
			Lt("value1", "value2"),
		},
		{
			"\"field\"<=$1",
			[]interface{}{10},
			Lte(I("field"), 10),
		},
		{
			"$1<=\"field\"",
			[]interface{}{"value"},
			Lte("value", I("field")),
		},
		{
			"$1<=$2",
			[]interface{}{"value1", "value2"},
			Lte("value1", "value2"),
		},
		{
			"\"field\">$1",
			[]interface{}{10},
			Gt(I("field"), 10),
		},
		{
			"$1>\"field\"",
			[]interface{}{"value"},
			Gt("value", I("field")),
		},
		{
			"$1>$2",
			[]interface{}{"value1", "value2"},
			Gt("value1", "value2"),
		},
		{
			"\"field\">=$1",
			[]interface{}{10},
			Gte(I("field"), 10),
		},
		{
			"$1>=\"field\"",
			[]interface{}{"value"},
			Gte("value", I("field")),
		},
		{
			"$1>=$2",
			[]interface{}{"value1", "value2"},
			Gte("value1", "value2"),
		},
		{
			"\"field\" IS NULL",
			nil,
			Nil("field"),
		},
		{
			"\"field\" IS NOT NULL",
			nil,
			NotNil("field"),
		},
		{
			"\"field\" IN ($1)",
			[]interface{}{"value1"},
			In("field", "value1"),
		},
		{
			"\"field\" IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			In("field", "value1", "value2"),
		},
		{
			"\"field\" IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			In("field", "value1", "value2", "value3"),
		},
		{
			"\"field\" NOT IN ($1)",
			[]interface{}{"value1"},
			Nin("field", "value1"),
		},
		{
			"\"field\" NOT IN ($1,$2)",
			[]interface{}{"value1", "value2"},
			Nin("field", "value1", "value2"),
		},
		{
			"\"field\" NOT IN ($1,$2,$3)",
			[]interface{}{"value1", "value2", "value3"},
			Nin("field", "value1", "value2", "value3"),
		},
		{
			"\"field\" LIKE $1",
			[]interface{}{"%value%"},
			Like("field", "%value%"),
		},
		{
			"\"field\" NOT LIKE $1",
			[]interface{}{"%value%"},
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			Fragment("FRAGMENT"),
		},
		{
			"(\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(\"field1\"=$1 AND \"field2\"=$2 AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(\"field1\"=$1 OR \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(\"field1\"=$1 OR \"field2\"=$2 OR \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"NOT (\"field1\"=$1 AND \"field2\"=$2)",
			[]interface{}{"value1", "value2"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"NOT (\"field1\"=$1 AND \"field2\"=$2 AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"((\"field1\"=$1 OR \"field2\"=$2) AND \"field3\"=$3)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Eq(I("field3"), "value3")),
		},
		{
			"((\"field1\"=$1 OR \"field2\"=$2) AND (\"field3\"=$3 OR \"field4\"=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4"))),
		},
		{
			"(NOT (\"field1\"=$1 AND \"field2\"=$2) AND NOT (\"field3\"=$3 OR \"field4\"=$4))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Not(Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4")))),
		},
		{
			"NOT (\"field1\"=$1 AND (\"field2\"=$2 OR \"field3\"=$3) AND NOT (\"field4\"=$4 OR \"field5\"=$5))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq(I("field1"), "value1"), Or(Eq(I("field2"), "value2"), Eq(I("field3"), "value3")), Not(Or(Eq(I("field4"), "value4"), Eq(I("field5"), "value5"))))),
		},
		{
			"((\"field1\" IN ($1,$2) OR \"field2\" NOT IN ($3)) AND \"field3\" IN ($4,$5,$6))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(\"field1\" LIKE $1 AND \"field2\" NOT LIKE $2)",
			[]interface{}{"%value1%", "%value2%"},
			And(Like(I("field1"), "%value1%"), NotLike(I(I("field2")), "%value2%")),
		},
		{
			"",
			nil,
			Condition{Type: ConditionType(9999)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := NewBuilder(&Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
			}).condition(tt.Condition)

			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestBuilder_Lock(t *testing.T) {
	query := grimoire.Query{
		Collection: "users",
		Fields:     []string{"*"},
		LockClause: "FOR UPDATE",
	}

	qs, args := NewBuilder(&Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}).Find(query)

	assert.Equal(t, "SELECT * FROM `users` FOR UPDATE;", qs)
	assert.Nil(t, args)
}
