package grimoire

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

// err := repo.Transaction(func(r Repo) {
// 	err := r.From("blablabla").Find("blablabla")
//	if err != nil {
//		return err
//	}
// 	r.Update("blablabla").Find("blablabla")
//	return nil
// })
