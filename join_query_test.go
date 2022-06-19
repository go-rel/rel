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
		Arguments: []interface{}{1},
	}, rel.NewJoinFragment("JOIN transactions ON id=?", 1))

	assert.Equal(t, rel.JoinQuery{
		Mode:      "JOIN transactions ON transactions.user_id=users.id",
		Arguments: []interface{}{},
	}, rel.NewJoinFragment("JOIN transactions ON transactions.user_id=users.id"))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
	}, rel.NewJoin("transactions"))
}

func TestJoinPopulate_hasOne(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoin("address")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "user_addresses",
		To:    "user_addresses.user_id",
		From:  "users.id",
	}, populated)
}

func TestJoinPopulate_hasOnePtr(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoin("work_address")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "user_addresses",
		To:    "user_addresses.user_id",
		From:  "users.id",
	}, populated)
}

func TestJoinPopulate_hasMany(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoin("transactions")).
			Populate(rel.NewDocument(&rel.User{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "transactions",
		To:    "transactions.user_id",
		From:  "users.id",
	}, populated)
}

func TestJoinPopulate_belongsTo(t *testing.T) {
	var (
		populated = rel.Build("", rel.NewJoin("user")).
			Populate(rel.NewDocument(&rel.Address{}, false).Meta()).
			JoinQuery[0]
	)

	assert.Equal(t, rel.JoinQuery{
		Mode:  "JOIN",
		Table: "users",
		To:    "users.id",
		From:  "user_addresses.user_id",
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
