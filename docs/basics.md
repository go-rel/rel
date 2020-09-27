# Basics

## Full Example

Below is a very basic example on how to utilize REL using mysql adapter.
Testing database query using REL can be done using [reltest](/reference/reltest/) package.

=== "main.go"

	{{ embed_code("examples/basics.go", prefix="\t") }}

=== "main_test.go"

	{{ embed_code("examples/basics_test.go", "example", "\t") }}

<!-- tabs:end -->

### Other Examples

- [go-todo-backend](https://github.com/Fs02/go-todo-backend) - Todo Backend

## Conventions

### Schema Definition

REL uses a struct as the schema to infer `table name`, `columns` and `primary field`.

```go
// Table name: books
type Book struct {
	ID        int       // id
	Title     string    // title
	Category  string    // category
	CreatedAt time.Time // created_at
	UpdatedAt time.Time // updated_at
}
```

### Table Name

Table name will be the pluralized struct name in snake case, you may create a `Table() string` method to override the default table name.

```go
// Default table name is `books`
type Book struct {}

// Override table name to be `ebooks`
func (b Book) Table() string {
	return "ebooks"
}
```

### Column Name

Column name will be the struct field name in snake case, you may override the column name by using using `db` tag.

```go
type Book struct {
	ID       int                // this field will be mapped to `id` column.
	Title    string `db:"name"` // this field will be mapped to `name` column.
	Category string `db:"-"`    // this field will be skipped
}
```

### Primary Key

REL requires every struct to have at least `primary` key. by default field named `id` will be used as primary key. To use other field as primary key, you may define it as `primary` using `db` tag. Defining multiple field as primary will forms composite primary key.


```go
type Book struct {
	UUID string `db:"uuid,primary"` // or just `db:",primary"`
}
```

### Timestamp

REL automatically track created and updated time of each struct if `CreatedAt` or `UpdatedAt` field exists.
