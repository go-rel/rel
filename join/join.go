package join

import (
	"github.com/Fs02/grimoire/query"
)

func Join(collection string) query.JoinClause {
	return query.NewJoin(collection)
}

func On(collection string, from string, to string) query.JoinClause {
	return query.NewJoinOn(collection, from, to)
}

func Inner(collection string) query.JoinClause {
	return query.NewInnerJoin(collection)
}

func InnerOn(collection string, from string, to string) query.JoinClause {
	return query.NewInnerJoinOn(collection, from, to)
}

func Left(collection string) query.JoinClause {
	return query.NewLeftJoin(collection)
}

func LeftOn(collection string, from string, to string) query.JoinClause {
	return query.NewLeftJoinOn(collection, from, to)
}

func Right(collection string) query.JoinClause {
	return query.NewRightJoin(collection)
}

func RightOn(collection string, from string, to string) query.JoinClause {
	return query.NewRightJoinOn(collection, from, to)
}

func Full(collection string) query.JoinClause {
	return query.NewFullJoin(collection)
}

func FullOn(collection string, from string, to string) query.JoinClause {
	return query.NewFullJoinOn(collection, from, to)
}
