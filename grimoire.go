package grimoire

// type Query struct {
// 	Collection      string
// 	Fields          []string
// 	AsDistinct      bool
// 	JoinClause      []JoinClause
// 	Condition       Condition
// 	GroupFields     []string
// 	HavingCondition Condition
// 	OrderClause     []OrderClause
// 	OffsetResult    int
// 	LimitResult     int
// }

// type JoinClause struct {
// 	Mode       string
// 	Collection string
// 	Condition  Condition
// }

// type OrderClause struct {
// 	Field string
// 	Order int
// }

// type Repo struct {
// 	adapter *Adapter
// 	query Query
// }

// func (repo Repo) From(collection string) Repo {
// 	repo.query = Query{
// 		Collection: collection,
// 		Fields:     []string{"*"},
// 	}

// 	return repo
// }

// func (r Repo) All(entities interface{}, q query.Query) error {
// 	qs, args := r.adapter.All(q)
// 	return r.adapter.Query(entities, qs, args)
// }

// func (r Repo) Insert(entity interface{}, ch *Changeset) error {
// 	qs, args := r.adapter.Insert(ch)
// 	id, _, err := r.adapter.Exec(qs, args)
// 	if err != nil {
// 		return err
// 	}

// 	q := From(ch.Collection).Where(Eq(I("id"), id))
// 	qs, args = r.adapter.All(q)
// 	return r.adapter.Query(entities, qs, args)
// }

// func (r Repo) Preload(field string, q query.Query) error {
// 	return nil
// }

// db := grimoire.New("mysql://blabla")
// db.From("books").Find(&book, 1)
// db.From("books").All(&books)
// db.From("books").Where("1 > ?", 2).All(&books)
// db.From("books").Where("1 > ?", 2).UpdateAll(&books)
// db.From("books").Join("").Where("1 > ?", 2).All(&books)

// err := db.From("books").Find(1).One(&book) // automatically limit 1
// err := db.From("books").Find(1).All(&book)
// err, count := db.From("books").Find(1).Count() // sql adapter automatically select count(*)

// err := db.From("books").Insert(&book, changes)
// err := db.From("books").Find(1).Update(&book, changes)
// err := db.From("books").Find(1).Replace(&book)
// err := db.From("books").Find(1).Delete()

// err := db.From("books").Find(1).Preload("Users", db.From("users")).All(&book)

// Find(interface{}) // is short hand, automatically where(id = ?, value)
