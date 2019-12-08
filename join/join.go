// Package join is syntatic sugar for building join query.
package join

import (
	"github.com/Fs02/rel"
)

var (
	Join    = rel.NewJoin
	On      = rel.NewJoinOn
	Inner   = rel.NewInnerJoin
	InnerOn = rel.NewInnerJoinOn
	Left    = rel.NewLeftJoin
	LeftOn  = rel.NewLeftJoinOn
	Right   = rel.NewRightJoin
	RightOn = rel.NewRightJoinOn
	Full    = rel.NewFullJoin
	FullOn  = rel.NewFullJoinOn
)
