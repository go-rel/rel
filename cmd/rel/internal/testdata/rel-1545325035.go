
package main

import (
	"context"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	db "github.com/go-rel/sqlite3"
	"github.com/go-rel/rel"
	"github.com/go-rel/migration"

	"github.com/go-rel/rel/cmd/rel/internal/testdata/migrations"
)

var (
	shutdowns []func() error
)

func logger(ctx context.Context, op string, message string) func(err error) {
	// no op for rel functions.
	if strings.HasPrefix(op, "rel-") {
		return func(error) {}
	}

	if op == "migrate" || op == "rollback" {
		log.Print("Running: ", op, " ", message)
	}

	t := time.Now()
	return func(err error) {
		duration := time.Since(t)
		if op == "migrate" || op == "rollback" {
			log.Print("=> Done: ", op, " ", message, " in ", duration)
		} else if false {
			log.Print("\t[duration: ", duration, " op: ", op, "] ", message)
		}

		if err != nil {
			log.Println("\tError: ", op, " ", err)
		}
	}
}

func main() {
	var (
		ctx = context.Background()
	)

	adapter, err := db.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	var (
		repo = rel.New(adapter)
		m    = migrator.New(repo)
	)

	log.SetFlags(0)
	repo.Instrumentation(logger)
	m.Instrumentation(logger)

	
	m.Register(1, migrations.MigrateCreateSamples, migrations.RollbackCreateSamples)
	

	m.Migrate(ctx)
}
