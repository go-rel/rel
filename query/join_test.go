package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestJoinWith(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewJoinWith("JOIN", "transactions", "user_id", "id"))
}

func TestJoinFragment(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:      "JOIN transactions ON id=?",
		Arguments: []interface{}{1},
	}, query.NewJoinFragment("JOIN transactions ON id=?", 1))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
	}, query.NewJoin("transactions"))
}

func TestJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewJoinOn("transactions", "user_id", "id"))
}

func TestInnerJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
	}, query.NewInnerJoin("transactions"))
}

func TestInnerJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewInnerJoinOn("transactions", "user_id", "id"))
}

func TestLeftJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
	}, query.NewLeftJoin("transactions"))
}

func TestLeftJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewLeftJoinOn("transactions", "user_id", "id"))
}

func TestRightJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
	}, query.NewRightJoin("transactions"))
}

func TestRightJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewRightJoinOn("transactions", "user_id", "id"))
}

func TestFullJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
	}, query.NewFullJoin("transactions"))
}

func TestFullJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.NewFullJoinOn("transactions", "user_id", "id"))
}
