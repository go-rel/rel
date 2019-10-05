package group

import (
	"github.com/Fs02/rel"
)

func By(fields ...string) rel.GroupQuery {
	return rel.NewGroup(fields...)
}

func Fields(fields ...string) rel.GroupQuery {
	return By(fields...)
}
