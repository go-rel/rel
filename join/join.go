// Package join is syntatic sugar for building join query.
package join

import (
	"github.com/go-rel/rel"
)

var (
	// Join is alias for rel.NewJoin
	Join = rel.NewJoin
	// On is alias for rel.NewJoinOn
	On = rel.NewJoinOn
	// Inner is alias for rel.NewInnerJoin
	Inner = rel.NewInnerJoin
	// InnerOn is alias for rel.NewInnerJoinOn
	InnerOn = rel.NewInnerJoinOn
	// Left is alias for rel.NewLeftJoin
	Left = rel.NewLeftJoin
	// LeftOn is alias for rel.NewLeftJoinOn
	LeftOn = rel.NewLeftJoinOn
	// Right is alias for rel.NewRightJoin
	Right = rel.NewRightJoin
	// RightOn is alias for rel.NewRightJoinOn
	RightOn = rel.NewRightJoinOn
	// Full is alias for rel.NewFullJoin
	Full = rel.NewFullJoin
	// FullOn is alias for rel.NewFullJoinOn
	FullOn = rel.NewFullJoinOn
	// AssocWith is alias for rel.NewJoinAssocWith
	AssocWith = rel.NewJoinAssocWith
	// Assoc is alias for rel.NewJoinAssoc
	Assoc = rel.NewJoinAssoc
)
