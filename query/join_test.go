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
	}, query.JoinWith("JOIN", "transactions", "user_id", "id"))
}

func TestJoinFragment(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:      "JOIN transactions ON id=?",
		Arguments: []interface{}{1},
	}, query.JoinFragment("JOIN transactions ON id=?", 1))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
	}, query.Join("transactions"))
}

func TestJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.JoinOn("transactions", "user_id", "id"))
}

func TestInnerJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
	}, query.InnerJoin("transactions"))
}

func TestInnerJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.InnerJoinOn("transactions", "user_id", "id"))
}

func TestLeftJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
	}, query.LeftJoin("transactions"))
}

func TestLeftJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.LeftJoinOn("transactions", "user_id", "id"))
}

func TestRightJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
	}, query.RightJoin("transactions"))
}

func TestRightJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.RightJoinOn("transactions", "user_id", "id"))
}

func TestFullJoin(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
	}, query.FullJoin("transactions"))
}

func TestFullJoinOn(t *testing.T) {
	assert.Equal(t, query.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, query.FullJoinOn("transactions", "user_id", "id"))
}
