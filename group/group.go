// Package group is syntatic sugar for building group query.
package group

import (
	"github.com/go-rel/rel"
)

var (
	// By is alias for rel.NewGroup
	By = rel.NewGroup
	// Fields is alias for rel.NewGroup
	Fields = rel.NewGroup
)
