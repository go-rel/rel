# Migration

## Overview

Migration is a feature that allows you to evolve your database schema over time, REL provides DSL that allows you to write migration in Golang.

## Defining Migration

Migration package usually located inside `your-repo/db/migrations` package. It's a standalone package that should not be imported by the rest of your application.
Each migration file is named as `number_name.go`, and each migration file must define a pair of migration and rollback functions: `MigrateName` and `RollbackName`. 
Migrate and rollback function name is the camel cased file name without version.

!!! note
    Sample project that demonstrate this setup can be found at https://github.com/Fs02/go-todo-backend

{{ embed_code("examples/migrations/20202806225100_create_todos.go") }}

## Running Migration

REL provides CLI that can be used to run your migration, it can be installed using `go get` or downloaded from [release page](https://github.com/go-rel/rel/releases).

```bash
go get github.com/go-rel/rel/cmd/rel
```

*Verify installation:*

```bash
rel -version
```

*Migrate to the latest version:*

```bash
rel migrate
```

*Rollback one migration step:*

```bash
rel rollback
```

## Configuring Database Connection

By default, REL will try to use database connection info that available as environment variable.

| Variable              | Description                                                   |
|-----------------------|---------------------------------------------------------------|
| `DATABASE_URL`        | Database connection string (Optional)                         |
| `DATABASE_ADAPTER`    | Adapter package (Required if `DATABASE_URL` specified)        |
| `DATABASE_DRIVER`     | Driver package (Required if `DATABASE_URL` specified)         |
| `MYSQL_HOST`          | MySQL host (Optional)                                         |
| `MYSQL_PORT`          | MySQL port (Optional)                                         |
| `MYSQL_DATABASE`      | MySQL database (Required, if `MYSQL_HOST` specified)          |
| `MYSQL_USERNAME`      | MySQL host (Required, if `MYSQL_HOST` specified)              |
| `MYSQL_PASSWORD`      | MySQL host (Optional)                                         |
| `POSTGRES_HOST`       | PostgreSQL host (Optional)                                    |
| `POSTGRES_PORT`       | PostgreSQL port (Optional)                                    |
| `POSTGRES_DATABASE`   | PostgreSQL database (Required, if `POSTGRES_HOST` specified)  |
| `POSTGRES_USERNAME`   | PostgreSQL username (Required, if `POSTGRES_HOST` specified   |
| `POSTGRES_PASSWORD`   | PostgreSQL password (Optional)                                |
| `SQLITE3_DATABASE`    | SQLite3 database (Optional)                                   |

*Database connection info can also be specified using command line options: `dsn`, `adapter` and `driver`:*

```bash
rel migrate -adapter=github.com/go-rel/rel/adapter/sqlite3 -driver=github.com/mattn/go-sqlite3 -dsn=:memory:
```
