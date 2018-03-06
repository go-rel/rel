package sql_test

import (
	"github.com/Fs02/grimoire/adapter/sql"
	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
