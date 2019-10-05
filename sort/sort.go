package sort

import (
	"github.com/Fs02/rel"
)

func Asc(field string) rel.SortQuery {
	return rel.NewSortAsc(field)
}

func Desc(field string) rel.SortQuery {
	return rel.NewSortDesc(field)
}
