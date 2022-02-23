package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/serenize/snaker"
)

const migrationTemplate = `
package main

import (
	"context"
	"log"
	"strings"
	"time"

	_ "{{.Driver}}"
	db "{{.Adapter}}"
	"github.com/go-rel/rel"
	"github.com/go-rel/migration"

	"{{.Package}}"
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
		} else if {{.Verbose}} {
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

	adapter, err := db.Open("{{.DSN}}")
	if err != nil {
		log.Fatal(err)
	}

	var (
		repo = rel.New(adapter)
		m    = migration.New(repo)
	)

	log.SetFlags(0)
	repo.Instrumentation(logger)
	m.Instrumentation(logger)

	{{range .Migrations}}
	m.Register({{.Version}}, migrations.Migrate{{.Name}}, migrations.Rollback{{.Name}})
	{{end}}

	{{.Command}}
}
`

var (
	tempdir           = ""
	stdout  io.Writer = os.Stdout
	stderr  io.Writer = os.Stderr
)

// ExecMigrate command.
// assumes args already validated.
func ExecMigrate(ctx context.Context, args []string) error {
	var (
		defAdapter, defDriver, defDSN = getDatabaseInfo()
		fs                            = flag.NewFlagSet(args[1], flag.ExitOnError)
		command                       = getMigrateCommand(args[1])
		dir                           = fs.String("dir", "db/migrations", "Path to directory containing migration files")
		module                        = fs.String("module", getModule(), "Module of the main package")
		adapter                       = fs.String("adapter", defAdapter, "Adapter package")
		driver                        = fs.String("driver", defDriver, "Driver package")
		dsn                           = fs.String("dsn", defDSN, "DSN for database connection")
		verbose                       = fs.Bool("verbose", false, "Show logs from REL")
		tmpl                          = template.Must(template.New("migration").Parse(migrationTemplate))
	)

	fs.Parse(args[2:])

	if *adapter == "" || *driver == "" || *dsn == "" {
		return fmt.Errorf("rel: missing required parameters:\n\tadapter: %s\n\tdriver: %s\n\tdsn: %s", *adapter, *driver, *dsn)
	}

	file, err := ioutil.TempFile(tempdir, "rel-*.go")
	check(err)
	defer os.Remove(file.Name())

	migrations, err := scanMigration(*dir)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, struct {
		Package    string
		Command    string
		Adapter    string
		Driver     string
		DSN        string
		Migrations []migration
		Verbose    bool
	}{
		Package:    *module + "/" + *dir,
		Command:    command,
		Adapter:    *adapter,
		Driver:     *driver,
		DSN:        *dsn,
		Migrations: migrations,
		Verbose:    *verbose,
	})
	check(err)
	check(file.Close())

	cmd := exec.CommandContext(ctx, "go", "run", "-mod=mod", file.Name(), "migrate")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

type migration struct {
	Version string
	Name    string
}

func scanMigration(dir string) ([]migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.New("rel: error accessing read migration directory: " + dir)
	}

	mFiles := make([]migration, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		result := reMigrationFile.FindStringSubmatch(f.Name())
		if len(result) < 3 {
			return nil, errors.New("rel: invalid migration file: " + f.Name())
		}

		mFiles = append(mFiles, migration{
			Version: result[1],
			Name:    snaker.SnakeToCamel(result[2]),
		})
	}

	return mFiles, err
}

func getMigrateCommand(cmd string) string {
	switch cmd {
	case "rollback", "down":
		return "m.Rollback(ctx)"
	default:
		return "m.Migrate(ctx)"
	}
}
