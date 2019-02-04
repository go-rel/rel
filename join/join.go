package join

import (
	"github.com/Fs02/grimoire/query"
)

func Join(collection string) query.Join {
	return query.JoinWith("JOIN", collection, "", "")
}

func On(collection string, from string, to string) query.Join {
	return query.JoinWith("JOIN", collection, from, to)
}

func Inner(collection string) query.Join {
	return query.JoinInnerOn(collection, "", "")
}

func InnerOn(collection string, from string, to string) query.Join {
	return query.JoinWith("INNER JOIN", collection, from, to)
}

func Left(collection string) query.Join {
	return query.JoinLeftOn(collection, "", "")
}

func LeftOn(collection string, from string, to string) query.Join {
	return query.JoinWith("LEFT JOIN", collection, from, to)
}

func Right(collection string) query.Join {
	return query.JoinRightOn(collection, "", "")
}

func RightOn(collection string, from string, to string) query.Join {
	return query.JoinWith("RIGHT JOIN", collection, from, to)
}

func Full(collection string) query.Join {
	return query.JoinFullOn(collection, "", "")
}

func FullOn(collection string, from string, to string) query.Join {
	return query.JoinWith("FULL JOIN", collection, from, to)
}
