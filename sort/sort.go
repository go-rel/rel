// Package sort is syntatic sugar for building sort query.
package sort

import (
	"github.com/go-rel/rel"
)

var (
	// Asc creates a query that sort the result ascending by specified field.
	Asc = rel.NewSortAsc

	// Desc creates a query that sort the result descending by specified field.
	Desc = rel.NewSortDesc
)
