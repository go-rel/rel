package grimoire

import (
	"testing"
	"time"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		DefaultLogger("", time.Second, nil)
		DefaultLogger("", time.Second, errors.NewUnexpected("error"))
	})
}

func TestLog(t *testing.T) {
	assert.NotPanics(t, func() {
		Log([]Logger{DefaultLogger}, "", time.Second, nil)
	})
}
