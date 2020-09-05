package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	var (
		ctx = context.Background()
	)

	if len(os.Args) < 2 {
		fmt.Println("Available command are: migrate, rollback")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate", "up", "rollback", "down":
		migrate(ctx)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
