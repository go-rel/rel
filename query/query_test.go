package query_test

import (
	"testing"

	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	assert.Equal(t, From("users"), Query{
		Collection: "users",
		Fields:     []string{"*"},
	})
}

func TestSelect(t *testing.T) {
	assert.Equal(t, From("users").Select("*"), Query{
		Collection: "users",
		Fields:     []string{"*"},
	})

	assert.Equal(t, From("users").Select("id", "name", "email"), Query{
		Collection: "users",
		Fields:     []string{"id", "name", "email"},
	})
}

func TestJoin(t *testing.T) {
	t.Skip("PENDING")
}

func TestJoinWith(t *testing.T) {
	t.Skip("PENDING")
}

func TestWhere(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")).Where(Ne("active", false)),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at"), Ne("active", false)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestOrWhere(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").OrWhere(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			From("users").Where(Eq("id", 1)).OrWhere(Nil("deleted_at")),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false)),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)).Where(Lt("price", 10000)),
			Query{
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))), Lt("price", 10000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestGroup(t *testing.T) {
	assert.Equal(t, From("users").Group("active", "plan"), Query{
		Collection:  "users",
		Fields:      []string{"*"},
		GroupFields: []string{"active", "plan"},
	})
}

func TestHaving(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")).Having(Ne("active", false)),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at"), Ne("active", false)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestOrHaving(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			From("users").OrHaving(Eq("id", 1), Nil("deleted_at")),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			From("users").Having(Eq("id", 1)).OrHaving(Nil("deleted_at")),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false)),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)).Having(Lt("price", 10000)),
			Query{
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))), Lt("price", 10000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestOrderBy(t *testing.T) {
	t.Skip("PENDING")
}

func TestOffset(t *testing.T) {
	assert.Equal(t, From("users").Offset(10), Query{
		Collection:   "users",
		Fields:       []string{"*"},
		OffsetResult: 10,
	})
}

func TestLimit(t *testing.T) {
	assert.Equal(t, From("users").Limit(10), Query{
		Collection:  "users",
		Fields:      []string{"*"},
		LimitResult: 10,
	})
}

func TestAsc(t *testing.T) {
	asc := Asc("id")

	assert.Equal(t, asc, OrderClause{
		Field: "id",
		Order: 1,
	})
	assert.True(t, asc.Asc())
}

func TestDesc(t *testing.T) {
	desc := Desc("id")

	assert.Equal(t, desc, OrderClause{
		Field: "id",
		Order: -1,
	})
	assert.True(t, desc.Desc())
}
