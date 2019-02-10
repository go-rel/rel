package query

import (
	"strings"
)

// JoinClause defines join information in query.
type JoinClause struct {
	Mode       string
	Collection string
	From       string
	To         string
	Arguments  []interface{}
}

func (j JoinClause) Build(query *Query) {
	if j.Arguments == nil && (j.From == "" || j.To == "") {
		j.From = query.Collection + "." + strings.TrimSuffix(j.Collection, "s") + "_id"
		j.To = j.Collection + ".id"
	}

	query.JoinClause = append(query.JoinClause, j)
}

func NewJoinWith(mode string, collection string, from string, to string) JoinClause {
	return JoinClause{
		Mode:       mode,
		Collection: collection,
		From:       from,
		To:         to,
	}
}

func NewJoinFragment(expr string, args ...interface{}) JoinClause {
	return JoinClause{
		Mode:      expr,
		Arguments: args,
	}
}

func NewJoin(collection string) JoinClause {
	return NewJoinWith("JOIN", collection, "", "")
}

func NewJoinOn(collection string, from string, to string) JoinClause {
	return NewJoinWith("JOIN", collection, from, to)
}

func NewInnerJoin(collection string) JoinClause {
	return NewInnerJoinOn(collection, "", "")
}

func NewInnerJoinOn(collection string, from string, to string) JoinClause {
	return NewJoinWith("INNER JOIN", collection, from, to)
}

func NewLeftJoin(collection string) JoinClause {
	return NewLeftJoinOn(collection, "", "")
}

func NewLeftJoinOn(collection string, from string, to string) JoinClause {
	return NewJoinWith("LEFT JOIN", collection, from, to)
}

func NewRightJoin(collection string) JoinClause {
	return NewRightJoinOn(collection, "", "")
}

func NewRightJoinOn(collection string, from string, to string) JoinClause {
	return NewJoinWith("RIGHT JOIN", collection, from, to)
}

func NewFullJoin(collection string) JoinClause {
	return NewFullJoinOn(collection, "", "")
}

func NewFullJoinOn(collection string, from string, to string) JoinClause {
	return NewJoinWith("FULL JOIN", collection, from, to)
}
