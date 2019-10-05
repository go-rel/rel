package rel

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		DefaultLogger("", time.Second, nil)
		DefaultLogger("", time.Second, errors.New("error"))
	})
}

func TestLog(t *testing.T) {
	assert.NotPanics(t, func() {
		Log([]Logger{DefaultLogger}, "", time.Second, nil)
	})
}
