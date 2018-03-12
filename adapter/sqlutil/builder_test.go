package sqlutil

import (
	"testing"

	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "SELECT *", builder.Select(false, "*"))
	assert.Equal(t, "SELECT id, name", builder.Select(false, "id", "name"))

	assert.Equal(t, "SELECT DISTINCT *", builder.Select(true, "*"))
	assert.Equal(t, "SELECT DISTINCT id, name", builder.Select(true, "id", "name"))
}

func TestFrom(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "FROM users", builder.From("users"))
}

func TestJoin(t *testing.T) {
	tests := []struct {
		QueryString string
		Args        []interface{}
		JoinClause  []JoinClause
	}{
		{
			"",
			nil,
			nil,
		},
		{
			"JOIN users ON user.id=trxs.user_id",
			nil,
			From("trxs").Join("users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"INNER JOIN users ON user.id=trxs.user_id",
			nil,
			From("trxs").JoinWith("INNER JOIN", "users", Eq(I("user.id"), I("trxs.user_id"))).JoinClause,
		},
		{
			"JOIN users ON user.id=trxs.user_id JOIN payments ON payments.id=trxs.payment_id",
			nil,
			From("trxs").Join("users", Eq(I("user.id"), I("trxs.user_id"))).
				Join("payments", Eq(I("payments.id"), I("trxs.payment_id"))).JoinClause,
		},
	}

	builder := Builder{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := builder.Join(tt.JoinClause...)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestWhere(t *testing.T) {
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
			"WHERE field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"WHERE (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	builder := Builder{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := builder.Where(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestGroupBy(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "", builder.GroupBy())
	assert.Equal(t, "GROUP BY city", builder.GroupBy("city"))
	assert.Equal(t, "GROUP BY city, nation", builder.GroupBy("city", "nation"))
}

func TestHaving(t *testing.T) {
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
			"HAVING field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"HAVING (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
	}

	builder := Builder{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := builder.Having(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}

func TestOrderBy(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "", builder.OrderBy())
	assert.Equal(t, "ORDER BY name ASC", builder.OrderBy(Asc("name")))
	assert.Equal(t, "ORDER BY name ASC, created_at DESC", builder.OrderBy(Asc("name"), Desc("created_at")))
}

func TestOffset(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "", builder.Offset(0))
	assert.Equal(t, "OFFSET 10", builder.Offset(10))
}

func TestLimit(t *testing.T) {
	builder := Builder{}

	assert.Equal(t, "", builder.Limit(0))
	assert.Equal(t, "LIMIT 10", builder.Limit(10))
}

func TestCondition(t *testing.T) {
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
			"field=?",
			[]interface{}{"value"},
			Eq(I("field"), "value"),
		},
		{
			"?=field",
			[]interface{}{"value"},
			Eq("value", I("field")),
		},
		{
			"?=?",
			[]interface{}{"value1", "value2"},
			Eq("value1", "value2"),
		},
		{
			"field<>?",
			[]interface{}{"value"},
			Ne(I("field"), "value"),
		},
		{
			"?<>field",
			[]interface{}{"value"},
			Ne("value", I("field")),
		},
		{
			"?<>?",
			[]interface{}{"value1", "value2"},
			Ne("value1", "value2"),
		},
		{
			"field<?",
			[]interface{}{10},
			Lt(I("field"), 10),
		},
		{
			"?<field",
			[]interface{}{"value"},
			Lt("value", I("field")),
		},
		{
			"?<?",
			[]interface{}{"value1", "value2"},
			Lt("value1", "value2"),
		},
		{
			"field<=?",
			[]interface{}{10},
			Lte(I("field"), 10),
		},
		{
			"?<=field",
			[]interface{}{"value"},
			Lte("value", I("field")),
		},
		{
			"?<=?",
			[]interface{}{"value1", "value2"},
			Lte("value1", "value2"),
		},
		{
			"field>?",
			[]interface{}{10},
			Gt(I("field"), 10),
		},
		{
			"?>field",
			[]interface{}{"value"},
			Gt("value", I("field")),
		},
		{
			"?>?",
			[]interface{}{"value1", "value2"},
			Gt("value1", "value2"),
		},
		{
			"field>=?",
			[]interface{}{10},
			Gte(I("field"), 10),
		},
		{
			"?>=field",
			[]interface{}{"value"},
			Gte("value", I("field")),
		},
		{
			"?>=?",
			[]interface{}{"value1", "value2"},
			Gte("value1", "value2"),
		},
		{
			"field IS NULL",
			nil,
			Nil("field"),
		},
		{
			"field IS NOT NULL",
			nil,
			NotNil("field"),
		},
		{
			"field IN (?)",
			[]interface{}{"value1"},
			In("field", "value1"),
		},
		{
			"field IN (?,?)",
			[]interface{}{"value1", "value2"},
			In("field", "value1", "value2"),
		},
		{
			"field IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			In("field", "value1", "value2", "value3"),
		},
		{
			"field NOT IN (?)",
			[]interface{}{"value1"},
			Nin("field", "value1"),
		},
		{
			"field NOT IN (?,?)",
			[]interface{}{"value1", "value2"},
			Nin("field", "value1", "value2"),
		},
		{
			"field NOT IN (?,?,?)",
			[]interface{}{"value1", "value2", "value3"},
			Nin("field", "value1", "value2", "value3"),
		},
		{
			"field LIKE ?",
			[]interface{}{"%value%"},
			Like("field", "%value%"),
		},
		{
			"field NOT LIKE ?",
			[]interface{}{"%value%"},
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT",
			nil,
			Fragment("FRAGMENT"),
		},
		{
			"(field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? AND field2=? AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=? OR field2=?)",
			[]interface{}{"value1", "value2"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? OR field2=? OR field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"(field1=? XOR field2=?)",
			[]interface{}{"value1", "value2"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"(field1=? XOR field2=? XOR field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Xor(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"NOT (field1=? AND field2=?)",
			[]interface{}{"value1", "value2"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")),
		},
		{
			"NOT (field1=? AND field2=? AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2"), Eq(I("field3"), "value3")),
		},
		{
			"((field1=? OR field2=?) AND field3=?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Eq(I("field3"), "value3")),
		},
		{
			"((field1=? OR field2=?) AND (field3=? OR field4=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4"))),
		},
		{
			"(NOT (field1=? AND field2=?) AND NOT (field3=? OR field4=?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq(I("field1"), "value1"), Eq(I("field2"), "value2")), Not(Or(Eq(I("field3"), "value3"), Eq(I("field4"), "value4")))),
		},
		{
			"NOT (field1=? AND (field2=? OR field3=?) AND NOT (field4=? OR field5=?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq(I("field1"), "value1"), Or(Eq(I("field2"), "value2"), Eq(I("field3"), "value3")), Not(Or(Eq(I("field4"), "value4"), Eq(I("field5"), "value5"))))),
		},
		{
			"((field1 IN (?,?) OR field2 NOT IN (?)) AND field3 IN (?,?,?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(field1 LIKE ? AND field2 NOT LIKE ?)",
			[]interface{}{"%value1%", "%value2%"},
			And(Like(I("field1"), "%value1%"), NotLike(I(I("field2")), "%value2%")),
		},
	}

	builder := Builder{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := builder.Condition(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}
