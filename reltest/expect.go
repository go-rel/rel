package reltest

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

var (
	// ErrConnectionClosed is alias for sql.ErrConnDone.
	ErrConnectionClosed = sql.ErrConnDone
)

// Expect is base behaviour for all reltest expectations.
type Expect struct {
	*mock.Call
}

// Error sets error to be returned.
func (e *Expect) Error(err error) {
	e.Return(err)
}

// ConnectionClosed sets this error to be returned.
func (e *Expect) ConnectionClosed() {
	e.Error(ErrConnectionClosed)
}

func newExpect(r *Repository, methodName string, args []interface{}, rets []interface{}) *Expect {
	return &Expect{
		Call: r.mock.On(methodName, args...).Return(rets...).Once(),
	}
}
