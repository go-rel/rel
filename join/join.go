// Package join is syntatic sugar for building join query.
package join

import (
	"github.com/Fs02/rel"
)

func Join(collection string) rel.JoinQuery {
	return rel.NewJoin(collection)
}

func On(collection string, from string, to string) rel.JoinQuery {
	return rel.NewJoinOn(collection, from, to)
}

func Inner(collection string) rel.JoinQuery {
	return rel.NewInnerJoin(collection)
}

func InnerOn(collection string, from string, to string) rel.JoinQuery {
	return rel.NewInnerJoinOn(collection, from, to)
}

func Left(collection string) rel.JoinQuery {
	return rel.NewLeftJoin(collection)
}

func LeftOn(collection string, from string, to string) rel.JoinQuery {
	return rel.NewLeftJoinOn(collection, from, to)
}

func Right(collection string) rel.JoinQuery {
	return rel.NewRightJoin(collection)
}

func RightOn(collection string, from string, to string) rel.JoinQuery {
	return rel.NewRightJoinOn(collection, from, to)
}

func Full(collection string) rel.JoinQuery {
	return rel.NewFullJoin(collection)
}

func FullOn(collection string, from string, to string) rel.JoinQuery {
	return rel.NewFullJoinOn(collection, from, to)
}
