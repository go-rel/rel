package group

import (
	"github.com/Fs02/grimoire"
)

func By(fields ...string) grimoire.GroupClause {
	return grimoire.NewGroup(fields...)
}

func Fields(fields ...string) grimoire.GroupClause {
	return By(fields...)
}
