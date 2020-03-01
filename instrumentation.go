package rel

import (
	"context"
	"log"
	"time"
)

// Instrumenter defines function type that can be used for instrumetation.
// This function should return a function with no argument as a callback for finished execution.
type Instrumenter func(ctx context.Context, op string, message string) func(err error)

// DefaultLogger log query suing standard log library.
func DefaultLogger(ctx context.Context, op string, message string) func(err error) {
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
