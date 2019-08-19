package sort

import (
	"github.com/Fs02/grimoire"
)

func Asc(field string) grimoire.SortClause {
	return grimoire.NewSortAsc(field)
}

func Desc(field string) grimoire.SortClause {
	return grimoire.NewSortDesc(field)
}
