package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getDsnFromEnv() (string, string, string) {
	var adapter, driver, url string
	switch {
	case os.Getenv("POSTGRESQL_HOST") != "":
		adapter = "github.com/Fs02/rel/adapter/postgres"
		driver = "github.com/lib/pq"
		url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRESQL_USERNAME"),
			os.Getenv("POSTGRESQL_PASSWORD"),
			os.Getenv("POSTGRESQL_HOST"),
			os.Getenv("POSTGRESQL_PORT"),
			os.Getenv("POSTGRESQL_DATABASE"))
	case os.Getenv("MYSQL_HOST") != "":
		adapter = "github.com/Fs02/rel/adapter/mysql"
		driver = "github.com/go-sql-driver/mysql"
		url = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("MYSQL_USERNAME"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"))
	case os.Getenv("SQLITE3_DATABASE") != "":
		adapter = "github.com/Fs02/rel/adapter/sqlite3"
		driver = "github.com/mattn/go-sqlite3"
		url = os.Getenv("SQLITE3_DATABASE")
	}

	return adapter, driver, url
}

func getModule() string {
	if module := getModuleFromGomod(); module != "" {
		return module
	}

	return getModuleFromGopath()
}

func getModuleFromGomod() string {
	data, err := ioutil.ReadFile("go.mod")
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
		gopath = os.Getenv("GOPATH")
		wd, _  = os.Getwd()
	)

	return strings.TrimPrefix(wd, gopath)
}
