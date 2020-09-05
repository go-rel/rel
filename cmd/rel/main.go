package main

import (
	"context"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	var (
		ctx = context.Background()
	)

	migrate(ctx, "db/migrations")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
