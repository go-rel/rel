package grimoire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var repo = Repo{}

func TestFrom(t *testing.T) {
	assert.Equal(t, repo.From("users"), Query{
		Collection: "users",
		Fields:     []string{"*"},
	})
}
