package internal

import (
	"context"
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/serenize/snaker"
)

const migrationTemplate = `
package main

import (
	"context"
	"log"
	{{if not .Verbose}} "io/ioutil" {{end}}

	_ "{{.Driver}}"
	db "{{.Adapter}}"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/migrator"

	"{{.Package}}"
)

var (
	shutdowns []func() error
)

func main() {
	var (
		ctx = context.Background()
	)

	{{if not .Verbose}}
	log.SetOutput(ioutil.Discard)
	{{end}}

	adapter, err := db.Open("{{.DSN}}")
	if err != nil {
		log.Fatal(err)
	}

	var (
		repo = rel.New(adapter)
		m    = migrator.New(repo)
	)

	{{range .Migrations}}
	m.Register({{.Version}}, migrations.Migrate{{.Name}}, migrations.Rollback{{.Name}})
	{{end}}

	{{.Command}}
}
`

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
		verbose                       = fs.Bool("verbose", true, "Show logs from REL")
		tmpl                          = template.Must(template.New("migration").Parse(migrationTemplate))
	)

	fs.Parse(args[2:])

	file, err := ioutil.TempFile(os.TempDir(), "rel-*.go")
	check(err)
	defer os.Remove(file.Name())

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
		Migrations: scanMigration(*dir),
		Verbose:    *verbose,
	})
	check(err)
	check(file.Close())

	cmd := exec.CommandContext(ctx, "go", "run", file.Name(), "migrate")
	output, err := cmd.CombinedOutput()
	print(string(output))
	return err
}

type migration struct {
	Version string
	Name    string
}

func scanMigration(dir string) []migration {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("rel: cannot access migration directory: " + err.Error())
	}

	mFiles := make([]migration, len(files))
	for i, f := range files {
		result := reMigrationFile.FindStringSubmatch(f.Name())
		if len(result) < 3 {
			panic("rel: invalid migration file: " + f.Name())
		}

		mFiles[i] = migration{
			Version: result[1],
			Name:    snaker.SnakeToCamel(result[2]),
		}
	}

	return mFiles
}

func getMigrateCommand(cmd string) string {
	switch cmd {
	case "rollback", "down":
		return "m.Rollback(ctx)"
	default:
		return "m.Migrate(ctx)"
	}
}
