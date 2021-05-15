package reltest

// Exec asserts and simulate exec function for test.
type Exec struct {
	*Expect
}

// Result sets the result of this query.
func (e *Exec) Result(lastInsertedId int, rowsAffected int) {
	e.Return(lastInsertedId, rowsAffected, nil)
}

// Error sets error to be returned.
func (e *Exec) Error(err error) {
	e.Return(0, 0, err)
}

// ConnectionClosed sets this error to be returned.
func (e *Exec) ConnectionClosed() {
	e.Error(ErrConnectionClosed)
}

// ExpectExec to be called with given field and queries.
func ExpectExec(r *Repository, statement string, args []interface{}) *Exec {
	return &Exec{
		Expect: newExpect(r, "Exec",
			[]interface{}{r.ctxData, statement, args},
			[]interface{}{0, 0, nil},
		),
	}
}
