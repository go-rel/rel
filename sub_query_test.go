package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	assert.Equal(t, SubQuery{
		Prefix: "ALL",
		Query:  Query{},
	}, All(Query{}))
}

func TestAny(t *testing.T) {
	assert.Equal(t, SubQuery{
		Prefix: "ANY",
		Query:  Query{},
	}, Any(Query{}))
}
