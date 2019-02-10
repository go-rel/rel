package sort

import (
	"github.com/Fs02/grimoire/query"
)

func Asc(field string) query.SortClause {
	return query.NewSortAsc(field)
}

func Desc(field string) query.SortClause {
	return query.NewSortDesc(field)
}
