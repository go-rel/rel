package schema

import (
	"context"

	"github.com/Fs02/rel"
)

// Adapter interface
type Adapter interface {
	rel.Adapter
	Apply(ctx context.Context, table Table) error
}
