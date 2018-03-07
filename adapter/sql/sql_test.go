package sql_test

import (
	"github.com/Fs02/grimoire/adapter/sql"
	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	builder := sql.QueryBuilder{}

	assert.Equal(t, "SELECT *", builder.Select(false, "*"))
	assert.Equal(t, "SELECT id, name", builder.Select(false, "id", "name"))

	assert.Equal(t, "SELECT DISTINCT *", builder.Select(true, "*"))
	assert.Equal(t, "SELECT DISTINCT id, name", builder.Select(true, "id", "name"))
}

func TestFrom(t *testing.T) {
	builder := sql.QueryBuilder{}

	assert.Equal(t, "FROM users", builder.From("users"))
}

func TestJoin(t *testing.T) {
	t.Skip("PENDING")
}

func TestWhere(t *testing.T) {
	t.Skip("PENDING")
}

func TestGroupBy(t *testing.T) {
	t.Skip("PENDING")
}

func TestHaving(t *testing.T) {
	t.Skip("PENDING")
}

func TestOrderBy(t *testing.T) {
	t.Skip("PENDING")
}

func TestOffset(t *testing.T) {
	builder := sql.QueryBuilder{}

	assert.Equal(t, "", builder.Offset(0))
	assert.Equal(t, "OFFSET 10", builder.Offset(10))
}

func TestLimit(t *testing.T) {
	builder := sql.QueryBuilder{}

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
			"field = ?",
			[]interface{}{"value"},
			Eq("field", "value"),
		},
		{
			"field <> ?",
			[]interface{}{"value"},
			Ne("field", "value"),
		},
		{
			"field < ?",
			[]interface{}{10},
			Lt("field", 10),
		},
		{
			"field <= ?",
			[]interface{}{10},
			Lte("field", 10),
		},
		{
			"field > ?",
			[]interface{}{10},
			Gt("field", 10),
		},
		{
			"field >= ?",
			[]interface{}{10},
			Gte("field", 10),
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
			"field LIKE \"%value%\"",
			nil,
			Like("field", "%value%"),
		},
		{
			"field NOT LIKE \"%value%\"",
			nil,
			NotLike("field", "%value%"),
		},
		{
			"FRAGMENT(\"field\", \"value\"",
			nil,
			Fragment("field", "FRAGMENT(\"field\", \"value\""),
		},
		{
			"(field1 = ? AND field2 = ?)",
			[]interface{}{"value1", "value2"},
			And(Eq("field1", "value1"), Eq("field2", "value2")),
		},
		{
			"(field1 = ? AND field2 = ? AND field3 = ?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Eq("field1", "value1"), Eq("field2", "value2"), Eq("field3", "value3")),
		},
		{
			"(field1 = ? OR field2 = ?)",
			[]interface{}{"value1", "value2"},
			Or(Eq("field1", "value1"), Eq("field2", "value2")),
		},
		{
			"(field1 = ? OR field2 = ? OR field3 = ?)",
			[]interface{}{"value1", "value2", "value3"},
			Or(Eq("field1", "value1"), Eq("field2", "value2"), Eq("field3", "value3")),
		},
		{
			"(field1 = ? XOR field2 = ?)",
			[]interface{}{"value1", "value2"},
			Xor(Eq("field1", "value1"), Eq("field2", "value2")),
		},
		{
			"(field1 = ? XOR field2 = ? XOR field3 = ?)",
			[]interface{}{"value1", "value2", "value3"},
			Xor(Eq("field1", "value1"), Eq("field2", "value2"), Eq("field3", "value3")),
		},
		{
			"NOT (field1 = ? AND field2 = ?)",
			[]interface{}{"value1", "value2"},
			Not(Eq("field1", "value1"), Eq("field2", "value2")),
		},
		{
			"NOT (field1 = ? AND field2 = ? AND field3 = ?)",
			[]interface{}{"value1", "value2", "value3"},
			Not(Eq("field1", "value1"), Eq("field2", "value2"), Eq("field3", "value3")),
		},
		{
			"((field1 = ? OR field2 = ?) AND field3 = ?)",
			[]interface{}{"value1", "value2", "value3"},
			And(Or(Eq("field1", "value1"), Eq("field2", "value2")), Eq("field3", "value3")),
		},
		{
			"((field1 = ? OR field2 = ?) AND (field3 = ? OR field4 = ?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Or(Eq("field1", "value1"), Eq("field2", "value2")), Or(Eq("field3", "value3"), Eq("field4", "value4"))),
		},
		{
			"(NOT (field1 = ? AND field2 = ?) AND NOT (field3 = ? OR field4 = ?))",
			[]interface{}{"value1", "value2", "value3", "value4"},
			And(Not(Eq("field1", "value1"), Eq("field2", "value2")), Not(Or(Eq("field3", "value3"), Eq("field4", "value4")))),
		},
		{
			"NOT (field1 = ? AND (field2 = ? OR field3 = ?) AND NOT (field4 = ? OR field5 = ?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5"},
			And(Not(Eq("field1", "value1"), Or(Eq("field2", "value2"), Eq("field3", "value3")), Not(Or(Eq("field4", "value4"), Eq("field5", "value5"))))),
		},
		{
			"((field1 IN (?,?) OR field2 NOT IN (?)) AND field3 IN (?,?,?))",
			[]interface{}{"value1", "value2", "value3", "value4", "value5", "value6"},
			And(Or(In("field1", "value1", "value2"), Nin("field2", "value3")), In("field3", "value4", "value5", "value6")),
		},
		{
			"(field1 LIKE \"%value1%\" AND field2 NOT LIKE \"%value2%\")",
			nil,
			And(Like("field1", "%value1%"), NotLike("field2", "%value2%")),
		},
	}

	builder := sql.QueryBuilder{}

	for _, tt := range tests {
		t.Run(tt.QueryString, func(t *testing.T) {
			qs, args := builder.Condition(tt.Condition)
			assert.Equal(t, tt.QueryString, qs)
			assert.Equal(t, tt.Args, args)
		})
	}
}
