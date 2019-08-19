package group

import (
	"github.com/Fs02/grimoire"
)

func By(fields ...string) grimoire.GroupQuery {
	return grimoire.NewGroup(fields...)
}

func Fields(fields ...string) grimoire.GroupQuery {
	return By(fields...)
}
