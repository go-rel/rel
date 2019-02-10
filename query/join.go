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

func JoinWith(mode string, collection string, from string, to string) JoinClause {
	return JoinClause{
		Mode:       mode,
		Collection: collection,
		From:       from,
		To:         to,
	}
}

func JoinFragment(expr string, args ...interface{}) JoinClause {
	return JoinClause{
		Mode:      expr,
		Arguments: args,
	}
}

func Join(collection string) JoinClause {
	return JoinWith("JOIN", collection, "", "")
}

func JoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("JOIN", collection, from, to)
}

func InnerJoin(collection string) JoinClause {
	return InnerJoinOn(collection, "", "")
}

func InnerJoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("INNER JOIN", collection, from, to)
}

func LeftJoin(collection string) JoinClause {
	return LeftJoinOn(collection, "", "")
}

func LeftJoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("LEFT JOIN", collection, from, to)
}

func RightJoin(collection string) JoinClause {
	return RightJoinOn(collection, "", "")
}

func RightJoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("RIGHT JOIN", collection, from, to)
}

func FullJoin(collection string) JoinClause {
	return FullJoinOn(collection, "", "")
}

func FullJoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("FULL JOIN", collection, from, to)
}
