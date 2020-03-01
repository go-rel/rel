package rel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		DefaultLogger(context.TODO(), "test", "test log")(nil)
		DefaultLogger(context.TODO(), "test", "test log with error")(errors.New("error"))
	})
}
