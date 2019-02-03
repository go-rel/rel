package sort

import (
	"github.com/Fs02/grimoire/query"
)

func Asc(field string) query.Sort {
	return query.SortAsc(field)
}

func Desc(field string) query.Sort {
	return query.SortDesc(field)
}
