package group

import (
	"github.com/Fs02/grimoire/query"
)

func By(fields ...string) query.Group {
	return query.GroupBy(fields...)
}

func Fields(fields ...string) query.Group {
	return By(fields...)
}
