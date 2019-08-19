package join

import (
	"github.com/Fs02/grimoire"
)

func Join(collection string) grimoire.JoinQuery {
	return grimoire.NewJoin(collection)
}

func On(collection string, from string, to string) grimoire.JoinQuery {
	return grimoire.NewJoinOn(collection, from, to)
}

func Inner(collection string) grimoire.JoinQuery {
	return grimoire.NewInnerJoin(collection)
}

func InnerOn(collection string, from string, to string) grimoire.JoinQuery {
	return grimoire.NewInnerJoinOn(collection, from, to)
}

func Left(collection string) grimoire.JoinQuery {
	return grimoire.NewLeftJoin(collection)
}

func LeftOn(collection string, from string, to string) grimoire.JoinQuery {
	return grimoire.NewLeftJoinOn(collection, from, to)
}

func Right(collection string) grimoire.JoinQuery {
	return grimoire.NewRightJoin(collection)
}

func RightOn(collection string, from string, to string) grimoire.JoinQuery {
	return grimoire.NewRightJoinOn(collection, from, to)
}

func Full(collection string) grimoire.JoinQuery {
	return grimoire.NewFullJoin(collection)
}

func FullOn(collection string, from string, to string) grimoire.JoinQuery {
	return grimoire.NewFullJoinOn(collection, from, to)
}
