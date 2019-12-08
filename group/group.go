// Package group is syntatic sugar for building group query.
package group

import (
	"github.com/Fs02/rel"
)

var (
	By     = rel.NewGroup
	Fields = rel.NewGroup
)
