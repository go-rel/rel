package query

// JoinClause defines join information in query.
type JoinClause struct {
	Mode       string
	Collection string
	From       string
	To         string
	Arguments  []interface{}
}

func (j JoinClause) Build(query *Query) {
	// TODO: infer from and to when not specified
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

// func Join(collection string) JoinClause {
// 	return JoinWith("JOIN", collection, "", "")
// }

func JoinOn(collection string, from string, to string) JoinClause {
	return JoinWith("JOIN", collection, from, to)
}

func JoinInner(collection string) JoinClause {
	return JoinInnerOn(collection, "", "")
}

func JoinInnerOn(collection string, from string, to string) JoinClause {
	return JoinWith("INNER JOIN", collection, from, to)
}

func JoinLeft(collection string) JoinClause {
	return JoinLeftOn(collection, "", "")
}

func JoinLeftOn(collection string, from string, to string) JoinClause {
	return JoinWith("LEFT JOIN", collection, from, to)
}

func JoinRight(collection string) JoinClause {
	return JoinRightOn(collection, "", "")
}

func JoinRightOn(collection string, from string, to string) JoinClause {
	return JoinWith("RIGHT JOIN", collection, from, to)
}

func JoinFull(collection string) JoinClause {
	return JoinFullOn(collection, "", "")
}

func JoinFullOn(collection string, from string, to string) JoinClause {
	return JoinWith("FULL JOIN", collection, from, to)
}
