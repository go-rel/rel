package main

import (
	"context"
	"log"
	"time"

	"github.com/Fs02/rel"
)

// Instrumentation docs example.
func Instrumentation(ctx context.Context, repo rel.Repository) {
	/// [instrumentation]
	repo.Instrumentation(func(ctx context.Context, op string, message string) func(err error) {
		t := time.Now()

		return func(err error) {
			duration := time.Since(t)
			log.Print("[duration: ", duration, " op: ", op, "] ", message, " - ", err)
		}
	})
	/// [instrumentation]
}
