package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnConflict(t *testing.T) {
	assert.Equal(t, OnConflict{Ignore: true}, OnConflictIgnore())
	assert.Equal(t, OnConflict{Keys: []string{"id"}, Ignore: true}, OnConflictKeyIgnore("id"))
	assert.Equal(t, OnConflict{Keys: []string{"id"}, Ignore: true}, OnConflictKeysIgnore([]string{"id"}))

	assert.Equal(t, OnConflict{Replace: true}, OnConflictReplace())
	assert.Equal(t, OnConflict{Keys: []string{"id"}, Replace: true}, OnConflictKeyReplace("id"))
	assert.Equal(t, OnConflict{Keys: []string{"id"}, Replace: true}, OnConflictKeysReplace([]string{"id"}))

	assert.Equal(t, OnConflict{Fragment: "sql", FragmentArgs: []any{1}}, OnConflictFragment("sql", 1))
}
