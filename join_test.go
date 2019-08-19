package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestJoinWith(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewJoinWith("JOIN", "transactions", "user_id", "id"))
}

func TestJoinFragment(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:      "JOIN transactions ON id=?",
		Arguments: []interface{}{1},
	}, grimoire.NewJoinFragment("JOIN transactions ON id=?", 1))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
	}, grimoire.NewJoin("transactions"))
}

func TestJoinOn(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewJoinOn("transactions", "user_id", "id"))
}

func TestInnerJoin(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
	}, grimoire.NewInnerJoin("transactions"))
}

func TestInnerJoinOn(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "INNER JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewInnerJoinOn("transactions", "user_id", "id"))
}

func TestLeftJoin(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
	}, grimoire.NewLeftJoin("transactions"))
}

func TestLeftJoinOn(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "LEFT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewLeftJoinOn("transactions", "user_id", "id"))
}

func TestRightJoin(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
	}, grimoire.NewRightJoin("transactions"))
}

func TestRightJoinOn(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "RIGHT JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewRightJoinOn("transactions", "user_id", "id"))
}

func TestFullJoin(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
	}, grimoire.NewFullJoin("transactions"))
}

func TestFullJoinOn(t *testing.T) {
	assert.Equal(t, grimoire.JoinClause{
		Mode:       "FULL JOIN",
		Collection: "transactions",
		From:       "user_id",
		To:         "id",
	}, grimoire.NewFullJoinOn("transactions", "user_id", "id"))
}
