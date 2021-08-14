package reltest

import (
	"context"
	"database/sql"
)

var (
	// ErrConnectionClosed is alias for sql.ErrConnDone.
	ErrConnectionClosed = sql.ErrConnDone
)

type Assert struct {
	ctxData       ctxData
	repeatability int // 0 means not limited
	totalCalls    int
	optional      bool
}

func (a *Assert) Once() {
	a.Times(1)
}

func (a *Assert) Times(times int) {
	a.repeatability = times
}

func (a *Assert) Maybe() {
	a.optional = true
}

// this function needs to be called as last condition of if
// otherwise recorded total calls will be wrong
func (a *Assert) call(ctx context.Context) bool {
	if a.ctxData != fetchContext(ctx) || (a.repeatability != 0 && a.totalCalls >= a.repeatability) {
		return false
	}

	a.totalCalls++
	return true
}

func (a Assert) assert() bool {
	return a.optional ||
		(a.repeatability == 0 && a.totalCalls > 0) ||
		(a.repeatability != 0 && a.totalCalls >= a.repeatability)
}
