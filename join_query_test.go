package rel_test

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestJoinWith(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewJoinWith("JOIN", "transactions", "user_id", "id"))
}

func TestJoinFragment(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:      "JOIN transactions ON id=?",
		Arguments: []interface{}{1},
	}, rel.NewJoinFragment("JOIN transactions ON id=?", 1))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "JOIN",
		Collection: "transactions",
	}, rel.NewJoin("transactions"))
}

func TestJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewJoinOn("transactions", "user_id", "id"))
}

func TestInnerJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "INNER JOIN",
		Collection: "transactions",
	}, rel.NewInnerJoin("transactions"))
}

func TestInnerJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "INNER JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewInnerJoinOn("transactions", "user_id", "id"))
}

func TestLeftJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
	}, rel.NewLeftJoin("transactions"))
}

func TestLeftJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewLeftJoinOn("transactions", "user_id", "id"))
}

func TestRightJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
	}, rel.NewRightJoin("transactions"))
}

func TestRightJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewRightJoinOn("transactions", "user_id", "id"))
}

func TestFullJoin(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "FULL JOIN",
		Collection: "transactions",
	}, rel.NewFullJoin("transactions"))
}

func TestFullJoinOn(t *testing.T) {
	assert.Equal(t, rel.JoinQuery{
		Mode:       "FULL JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, rel.NewFullJoinOn("transactions", "user_id", "id"))
}
