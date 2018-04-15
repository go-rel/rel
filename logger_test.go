package grimoire

import (
	"testing"
	"time"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		logger("", time.Second, nil)
		logger("", time.Second, errors.UnexpectedError("error"))
	})
}
