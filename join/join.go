package join

import (
	"github.com/Fs02/grimoire/query"
)

func Join(collection string) query.JoinClause {
	return query.Join(collection)
}

func On(collection string, from string, to string) query.JoinClause {
	return query.JoinOn(collection, from, to)
}

func Inner(collection string) query.JoinClause {
	return query.InnerJoin(collection)
}

func InnerOn(collection string, from string, to string) query.JoinClause {
	return query.InnerJoinOn(collection, from, to)
}

func Left(collection string) query.JoinClause {
	return query.LeftJoin(collection)
}

func LeftOn(collection string, from string, to string) query.JoinClause {
	return query.LeftJoinOn(collection, from, to)
}

func Right(collection string) query.JoinClause {
	return query.RightJoin(collection)
}

func RightOn(collection string, from string, to string) query.JoinClause {
	return query.RightJoinOn(collection, from, to)
}

func Full(collection string) query.JoinClause {
	return query.FullJoin(collection)
}

func FullOn(collection string, from string, to string) query.JoinClause {
	return query.FullJoinOn(collection, from, to)
}
