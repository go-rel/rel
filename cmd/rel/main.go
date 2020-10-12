package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-rel/rel/cmd/rel/internal"
	"github.com/subosito/gotenv"
)

var (
	version = ""
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
	case "version", "-v", "-version":
		fmt.Println("REL CLI " + version)
	case "-help":
		fmt.Println("Usage: rel [command] -help")
		fmt.Println("Available commands: migrate, rollback")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}
