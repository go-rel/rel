// Package sort is syntatic sugar for building sort query.
package sort

import (
	"github.com/Fs02/rel"
)

// Asc creates a query that sort the result ascending by specified field.
func Asc(field string) rel.SortQuery {
	return rel.NewSortAsc(field)
}

// Desc creates a query that sort the result descending by specified field.
func Desc(field string) rel.SortQuery {
	return rel.NewSortDesc(field)
}
