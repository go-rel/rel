package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	reMigrationFile = regexp.MustCompile(`^(\d+)_([a-z_]+)\.go$`)
	reGomod         = regexp.MustCompile(`module\s(\S+)`)
	gomod           = "go.mod"
)

func getDatabaseInfo() (string, string, string) {
	var adapter, driver, dsn string
	switch {
	case os.Getenv("DATABASE_URL") != "":
		adapter = os.Getenv("DATABASE_ADAPTER")
		driver = os.Getenv("DATABASE_DRIVER")
		dsn = os.Getenv("DATABASE_URL")
	case os.Getenv("SQLITE3_DATABASE") != "":
		adapter = "github.com/go-rel/sqlite3"
		driver = "github.com/mattn/go-sqlite3"
		dsn = os.Getenv("SQLITE3_DATABASE")
	case os.Getenv("MYSQL_HOST") != "":
		adapter = "github.com/go-rel/mysql"
		driver = "github.com/go-sql-driver/mysql"
		dsn = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("MYSQL_USERNAME"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"))
	case os.Getenv("POSTGRES_HOST") != "":
		adapter = "github.com/go-rel/postgres"
		driver = "github.com/lib/pq"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRES_USERNAME"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DATABASE"))
	case os.Getenv("POSTGRESQL_HOST") != "":
		adapter = "github.com/go-rel/postgres"
		driver = "github.com/lib/pq"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRESQL_USERNAME"),
			os.Getenv("POSTGRESQL_PASSWORD"),
			os.Getenv("POSTGRESQL_HOST"),
			os.Getenv("POSTGRESQL_PORT"),
			os.Getenv("POSTGRESQL_DATABASE"))
	}

	return adapter, driver, dsn
}

func getModule() string {
	if module := getModuleFromGomod(); module != "" {
		return module
	}

	return getModuleFromGopath()
}

func getModuleFromGomod() string {
	data, err := ioutil.ReadFile(gomod)
	if err != nil {
		return ""
	}

	result := reGomod.FindSubmatch(data)
	if len(result) < 2 {
		return ""
	}

	return string(result[1])
}

func getModuleFromGopath() string {
	var (
		gopath = os.Getenv("GOPATH") + "/src/"
		wd, _  = os.Getwd()
	)

	return strings.TrimPrefix(wd, gopath)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
