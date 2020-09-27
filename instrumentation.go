package rel

import (
	"context"
	"log"
	"strings"
	"time"
)

// Instrumenter defines function type that can be used for instrumetation.
// This function should return a function with no argument as a callback for finished execution.
type Instrumenter func(ctx context.Context, op string, message string) func(err error)

// Observe operation.
func (i Instrumenter) Observe(ctx context.Context, op string, message string) func(err error) {
	if i != nil {
		return i(ctx, op, message)
	}

	return func(err error) {}
}

// DefaultLogger instrumentation to log queries and rel operation.
func DefaultLogger(ctx context.Context, op string, message string) func(err error) {
	// no op for rel functions.
	if strings.HasPrefix(op, "rel-") {
		return func(error) {}
	}

	t := time.Now()

	return func(err error) {
		duration := time.Since(t)
		if err != nil {
			log.Print("[duration: ", duration, " op: ", op, "] ", message, " - ", err)
		} else {
			log.Print("[duration: ", duration, " op: ", op, "] ", message)
		}
	}
}
