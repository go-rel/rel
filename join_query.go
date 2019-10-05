package rel

import (
	"strings"
)

// JoinQuery defines join information in query.
type JoinQuery struct {
	Mode       string
	Collection string
	From       string
	To         string
	Arguments  []interface{}
}

func (jq JoinQuery) Build(query *Query) {
	query.JoinQuery = append(query.JoinQuery, jq)
}

func (jq *JoinQuery) buildJoin(query Query) {
	if jq.Arguments == nil && (jq.From == "" || jq.To == "") {
		jq.From = query.Collection + "." + strings.TrimSuffix(jq.Collection, "s") + "_id"
		jq.To = jq.Collection + ".id"
	}
}

func NewJoinWith(mode string, collection string, from string, to string) JoinQuery {
	return JoinQuery{
		Mode:       mode,
		Collection: collection,
		From:       from,
		To:         to,
	}
}

func NewJoinFragment(expr string, args ...interface{}) JoinQuery {
	return JoinQuery{
		Mode:      expr,
		Arguments: args,
	}
}

func NewJoin(collection string) JoinQuery {
	return NewJoinWith("JOIN", collection, "", "")
}

func NewJoinOn(collection string, from string, to string) JoinQuery {
	return NewJoinWith("JOIN", collection, from, to)
}

func NewInnerJoin(collection string) JoinQuery {
	return NewInnerJoinOn(collection, "", "")
}

func NewInnerJoinOn(collection string, from string, to string) JoinQuery {
	return NewJoinWith("INNER JOIN", collection, from, to)
}

func NewLeftJoin(collection string) JoinQuery {
	return NewLeftJoinOn(collection, "", "")
}

func NewLeftJoinOn(collection string, from string, to string) JoinQuery {
	return NewJoinWith("LEFT JOIN", collection, from, to)
}

func NewRightJoin(collection string) JoinQuery {
	return NewRightJoinOn(collection, "", "")
}

func NewRightJoinOn(collection string, from string, to string) JoinQuery {
	return NewJoinWith("RIGHT JOIN", collection, from, to)
}

func NewFullJoin(collection string) JoinQuery {
	return NewFullJoinOn(collection, "", "")
}

func NewFullJoinOn(collection string, from string, to string) JoinQuery {
	return NewJoinWith("FULL JOIN", collection, from, to)
}
