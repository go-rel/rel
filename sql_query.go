package rel

import "strings"

// SQLQuery allows querying using native query supported by database.
type SQLQuery struct {
	Statement string
	Values    []interface{}
}

// Build Raw Query.
func (sq SQLQuery) Build(query *Query) {
	query.SQLQuery = sq
}

func (sq SQLQuery) String() string {
	var builder strings.Builder
	builder.WriteString("rel.SQL(\"")
	builder.WriteString(sq.Statement)
	builder.WriteString("\"")

	if len(sq.Values) != 0 {
		builder.WriteString(", ")
		builder.WriteString(fmtifaces(sq.Values))
	}

	builder.WriteString(")")
	return builder.String()
}

// SQL Query.
func SQL(statement string, values ...interface{}) SQLQuery {
	return SQLQuery{
		Statement: statement,
		Values:    values,
	}
}
