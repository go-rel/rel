package reltest

import (
	"context"
	"fmt"
	"reflect"
)

// T is an interface wrapper around *testing.T
type T interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

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

func (a Assert) assert(t T, mock interface{}) bool {
	if a.optional ||
		(a.repeatability == 0 && a.totalCalls > 0) ||
		(a.repeatability != 0 && a.totalCalls >= a.repeatability) {
		return true
	}

	// TODO: stacktrace not correct
	if a.repeatability > 0 {
		t.Errorf("FAIL: Need to make %d more call(s) to satisfy mock:\n\t%s", a.repeatability-a.totalCalls, mock)
	} else {
		t.Errorf("FAIL: Mock defined but not called:\n\t%s", mock)
	}

	return false
}

func failExecuteMessage(call interface{}, mocks interface{}) string {
	var (
		mocksStr      string
		callStr       = call.(interface{ String() string }).String()
		expectCallStr = call.(interface{ ExpectString() string }).ExpectString()
		rv            = reflect.ValueOf(mocks)
	)

	for i := 0; i < rv.Len(); i++ {
		mocksStr += fmt.Sprintf("\n\t- %s", rv.Index(i).Interface())
	}

	if mocksStr == "" {
		mocksStr = "None"
	}

	return fmt.Sprintf("FAIL: this call is not mocked:\n\t%s\nMaybe try adding mock:\n\t%s\n\nMocked calls:%s\n\n", callStr, expectCallStr, mocksStr)
}
