package sort

import (
	"github.com/Fs02/grimoire"
)

func Asc(field string) grimoire.SortQuery {
	return grimoire.NewSortAsc(field)
}

func Desc(field string) grimoire.SortQuery {
	return grimoire.NewSortDesc(field)
}
