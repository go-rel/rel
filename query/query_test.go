package query_test

import (
	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrom(t *testing.T) {
	assert.Equal(t, From("users"), Query{
		Collection: "users",
		Fields: []string{"*"},
	})
}

func TestSelect(t *testing.T) {
	assert.Equal(t, From("users").Select("*"), Query{
		Collection: "users",
		Fields: []string{"*"},
	})

	assert.Equal(t, From("users").Select("id", "name", "email"), Query{
		Collection: "users",
		Fields: []string{"id", "name", "email"},
	})
}

func TestJoin(t *testing.T) {
	t.Skip("PENDING")
}

func TestJoinWith(t *testing.T) {
	t.Skip("PENDING")
}

func TestWhere(t *testing.T) {
	t.Skip("PENDING")
}

func TestOrWhere(t *testing.T) {
	t.Skip("PENDING")
}

func TestGroupBy(t *testing.T) {
	assert.Equal(t, From("users").GroupBy("active", "plan"), Query{
		Collection: "users",
		Fields: []string{"*"},
		GroupFields: []string{"active", "plan"},
	})
}

func TestHaving(t *testing.T) {
	t.Skip("PENDING")
}

func TestOrHaving(t *testing.T) {
	t.Skip("PENDING")
}

func TestOrderBy(t *testing.T) {
	t.Skip("PENDING")
}

func TestOffset(t *testing.T) {
	assert.Equal(t, From("users").Offset(10), Query{
		Collection: "users",
		Fields: []string{"*"},
		OffsetResult: 10,
	})
}

func TestLimit(t *testing.T) {
	assert.Equal(t, From("users").Limit(10), Query{
		Collection: "users",
		Fields: []string{"*"},
		LimitResult: 10,
	})
}

func TestAsc(t *testing.T) {
	asc := Asc("id")

	assert.Equal(t, asc, OrderQuery{
		Field: "id",
		Order: 1,
	})
	assert.True(t, asc.Asc())
}

func TestDesc(t *testing.T) {
	desc := Desc("id")

	assert.Equal(t, desc, OrderQuery{
		Field: "id",
		Order: -1,
	})
	assert.True(t, desc.Desc())
}
