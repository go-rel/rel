package rel_test

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestJoinWith(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewJoinWith("JOIN", "transactions", "user_id", "id"))
}

func TestJoinWithFilter(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
		Filter: rel.FilterQuery{
			Type:  rel.FilterEqOp,
			Field: "deleted",
			Value: false,
		},
	}, rel.NewJoinWith("JOIN", "transactions", "user_id", "id", rel.Eq("deleted", false)))
}

func TestJoinWithMultipleFilters(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
		Filter: rel.FilterQuery{
			Type: rel.FilterAndOp,
			Inner: []rel.FilterQuery{
				{
					Type:  rel.FilterEqOp,
					Field: "user_id",
					Value: 5,
				},
				{
					Type:  rel.FilterEqOp,
					Field: "deleted",
					Value: false,
				},
			},
		},
	}, rel.NewJoin("transactions", rel.Eq("user_id", 5), rel.Eq("deleted", false)))
}

func TestJoinFragment(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:      "JOIN transactions ON id=?",
		Arguments: []any{1},
	}, rel.NewJoinFragment("JOIN transactions ON id=?", 1))

	assert.Equal(t, rel.JoinQuery{
		Mode:      "JOIN transactions ON transactions.user_id=users.id",
		Arguments: []any{},
	}, rel.NewJoinFragment("JOIN transactions ON transactions.user_id=users.id"))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
	}, rel.NewJoin("transactions"))
}

func TestJoinAssoc_hasOne(t *testing.T) {
	var (
		populated = rel.Build("", rel.Select("*", "address.*"), rel.NewJoinAssoc("address")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta())
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "user_addresses as address",
		To:    "address.user_id",
		From:  "users.id",
		Assoc: "address",
	}, populated.JoinQuery[0])
	assert.Equal(t, []string{
		"*",
		"address.id as address.id",
		"address.user_id as address.user_id",
		"address.street as address.street",
		"address.notes as address.notes",
		"address.deleted_at as address.deleted_at",
	}, populated.SelectQuery.Fields)
}

func TestJoinPopulate_hasOnePtr(t *testing.T) {
	var (
		populated = rel.Build("", rel.Select("id", "work_address.*", "name"), rel.NewJoinAssoc("work_address")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta())
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "user_addresses as work_address",
		To:    "work_address.user_id",
		From:  "users.id",
		Assoc: "work_address",
	}, populated.JoinQuery[0])
	assert.Equal(t, []string{
		"id",
		"name",
		"work_address.id as work_address.id",
		"work_address.user_id as work_address.user_id",
		"work_address.street as work_address.street",
		"work_address.notes as work_address.notes",
		"work_address.deleted_at as work_address.deleted_at",
	}, populated.SelectQuery.Fields)
}

func TestJoinPopulate_hasMany(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoinAssoc("transactions")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions as transactions",
		To:    "transactions.user_id",
		From:  "users.id",
		Assoc: "transactions",
	}, populated)
}

func TestJoinAssoc_belongsTo(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoinAssoc("user")).
			Populate(rel.NewDocument(&rel.Address{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "users as user",
		To:    "user.id",
		From:  "user_addresses.user_id",
		Assoc: "user",
	}, populated)
}

func TestJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewJoinOn("transactions", "user_id", "id"))
}

func TestInnerJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "INNER JOIN",
		Table: "transactions",
	}, rel.NewInnerJoin("transactions"))
}

func TestInnerJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "INNER JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewInnerJoinOn("transactions", "user_id", "id"))
}

func TestLeftJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "LEFT JOIN",
		Table: "transactions",
	}, rel.NewLeftJoin("transactions"))
}

func TestLeftJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "LEFT JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewLeftJoinOn("transactions", "user_id", "id"))
}

func TestRightJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "RIGHT JOIN",
		Table: "transactions",
	}, rel.NewRightJoin("transactions"))
}

func TestRightJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "RIGHT JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewRightJoinOn("transactions", "user_id", "id"))
}

func TestFullJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "FULL JOIN",
		Table: "transactions",
	}, rel.NewFullJoin("transactions"))
}

func TestFullJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "FULL JOIN",
		Table: "transactions",
		From:  "user_id",
		To:    "id",
	}, rel.NewFullJoinOn("transactions", "user_id", "id"))
}
