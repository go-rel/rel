
## Comparison with other ORMs

Most (if not all) orm for golang is written as a chainable API, meaning all of the query need to be called before performing actual action as a chain of method invocations. example:

```go
db.Where("id = ?", 1).First(&user)
```

Chainable api is very hard to be unit tested without writing a wrapper. One way to make it testable is to make an interface that also acts as a wrapper, which is usually ends up as its own repository package resides somewhere in your project:

```go
// mockable interface.
type UserRepository interface {
	Find(user *User, int id) error
}

// actual implementation
type userRepository struct{
	db *DB
}

func (ur userRepository) Find(user *User, int id) error {
	return db.Where("id = ?", 1).First(&user)
}
```

Compared to other orm, rel api is built with [testability](https://godoc.org/github.com/Fs02/rel/reltest) in mind. rel uses [interface](https://godoc.org/github.com/Fs02/rel#Repository) to define contract of every database query or execution, all while making a chainable query possible. The ultimate goal of rel is to be **your repository package without the needs of making your own wrapper**. example:

```golang
// rel repository
repo.Find(&user, where.Eq("id", 1))
```
