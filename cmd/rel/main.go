package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Fs02/rel/cmd/rel/internal"
	"github.com/subosito/gotenv"
)

func main() {
	log.SetFlags(0)
	gotenv.Load()

	var (
		err error
		ctx = context.Background()
	)

	if len(os.Args) < 2 {
		fmt.Println("Available command are: migrate, rollback")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate", "up", "rollback", "down":
		err = internal.ExecMigrate(ctx, os.Args)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}
