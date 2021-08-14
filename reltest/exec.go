package reltest

import (
	"context"
	"reflect"
)

type exec []*MockExec

func (e *exec) register(ctxData ctxData, statement string, args ...interface{}) *MockExec {
	me := &MockExec{ctxData: ctxData, argStatement: statement, argArgs: args}
	*e = append(*e, me)
	return me
}

func (e exec) execute(ctx context.Context, statement string, args ...interface{}) (int, int, error) {
	for _, me := range e {
		if fetchContext(ctx) == me.ctxData &&
			me.argStatement == statement &&
			reflect.DeepEqual(me.argArgs, args) {
			return me.retLastInsertedId, me.retRowsAffected, me.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockExec asserts and simulate UpdateAny function for test.
type MockExec struct {
	ctxData           ctxData
	argStatement      string
	argArgs           []interface{}
	retLastInsertedId int
	retRowsAffected   int
	retError          error
}

// Result sets the result of this query.
func (me *MockExec) Result(lastInsertedId int, rowsAffected int) {
	me.retLastInsertedId = lastInsertedId
	me.retRowsAffected = rowsAffected
}

// Error sets error to be returned.
func (me *MockExec) Error(err error) {
	me.retError = err
}

// ConnectionClosed sets this error to be returned.
func (me *MockExec) ConnectionClosed() {
	me.Error(ErrConnectionClosed)
}
