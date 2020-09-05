package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/serenize/snaker"
)

const migrationTemplate = `
package main

import (
	"context"

	_ "{{.Driver}}"
	db "{{.Adapter}}"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/migrator"
	"go.uber.org/zap"

	"{{.Package}}"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
	shutdowns []func() error
)

func main() {
	var (
		ctx = context.Background()
	)

	adapter, err := db.Open("{{.DSN}}")
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
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

var (
	reMigrationFile = regexp.MustCompile(`^(\d+)_([a-z_]+)\.go$`)
	reGomod         = regexp.MustCompile(`module\s(\S+)`)
)

type migrationTemplateData struct {
	Package    string
	Command    string
	Adapter    string
	Driver     string
	DSN        string
	Migrations []migration
}

func migrate(ctx context.Context, dir string) {
	var (
		err  error
		file *os.File
		tmpl = template.Must(template.New("migration").Parse(migrationTemplate))
		data = migrationTemplateData{
			Package:    getModule() + "/" + dir,
			Command:    "m.Migrate(ctx)",
			Migrations: scanMigration(dir),
		}
	)

	data.Adapter, data.Driver, data.DSN = getDsnFromEnv()

	file, err = ioutil.TempFile(os.TempDir(), "rel-*.go")
	check(err)
	defer os.RemoveAll(file.Name())

	err = tmpl.Execute(file, data)
	check(err)
	check(file.Close())

	cmd := exec.CommandContext(ctx, "go", "run", file.Name(), "migrate")
	output, err := cmd.CombinedOutput()
	print(string(output))
	os.Exit(cmd.ProcessState.ExitCode())
}

type migration struct {
	Version string
	Name    string
}

func scanMigration(dir string) []migration {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	mFiles := make([]migration, len(files))
	for i, f := range files {
		result := reMigrationFile.FindStringSubmatch(f.Name())
		if len(result) < 3 {
			log.Fatal("invalid migration file: ", result)
		}

		mFiles[i] = migration{
			Version: result[1],
			Name:    snaker.SnakeToCamel(result[2]),
		}
	}

	return mFiles
}
