package main

import (
	"context"
	"io"

	"github.com/Fs02/rel"
)

// User model example.
type User struct {
	ID   int
	Name string
}

// SendPromotionEmail tp demonstrate Iteration.
func SendPromotionEmail(*User) {}

// Iteration docs example.
func Iteration(ctx context.Context, repo rel.Repository) error {
	/// [batch-iteration]
	var (
		user User
		iter = repo.Iterate(ctx, rel.From("users"), rel.BatchSize(500))
	)

	// make sure iterator is closed after process is finish.
	defer iter.Close()
	for {
		// retrieve next user.
		if err := iter.Next(&user); err != nil {
			if err == io.EOF {
				break
			}

			// handle error
			return err
		}

		// process user
		SendPromotionEmail(&user)
	}
	/// [batch-iteration]

	return nil
}
