package rel_test

import (
	"testing"

	"github.com/Fs02/rel"
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
