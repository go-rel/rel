package grimoire

import (
	"testing"

	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/query"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var repo = Repo{}

func createCursor(row int) *testCursor {
	cur := &testCursor{}

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id"}, nil).Once()

	if row > 0 {
		cur.On("Next").Return(true).Times(row)
		cur.SetScan(row, 10)
	}

	cur.On("Next").Return(false).Once()

	return cur
}

func TestNew(t *testing.T) {
	adapter := &testAdapter{}
	repo := New(adapter)

	assert.NotNil(t, repo)
	assert.Equal(t, adapter, repo.Adapter())
}

func TestRepo_SetLogger(t *testing.T) {
	repo := Repo{}
	assert.Nil(t, repo.logger)
	repo.SetLogger(DefaultLogger)
	assert.NotNil(t, repo.logger)
}

func TestRepo_Aggregate(t *testing.T) {
	adapter := &testAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users")
	mode := "COUNT"
	field := "*"

	var out struct {
		Count int
	}

	adapter.On("Aggregate", query, &out, mode, field).Return(nil)

	err := repo.Aggregate(User{}, "COUNT", "*", &out, query)
	assert.Nil(t, err)

	assert.NotPanics(t, func() {
		repo.MustAggregate(User{}, "COUNT", "*", &out, query)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_Count(t *testing.T) {
	adapter := &testAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users")

	adapter.On("Aggregate", query, mock.Anything, "COUNT", "*").Return(nil)

	count, err := repo.Count(User{}, query)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	assert.NotPanics(t, func() {
		assert.Equal(t, 0, repo.MustCount(User{}, query))
	})

	adapter.AssertExpectations(t)
}

func TestRepo_One(t *testing.T) {
	var (
		user    User
		doc     = newDocument(&user)
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		query   = query.From("users").Limit(1)
		cur     = createCursor(1)
	)

	doc.(*document).reflect()

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.One(&user, query))
	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_One_queryError(t *testing.T) {
	var (
		user    User
		doc     = newDocument(&user)
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		cur     = &testCursor{}
		query   = query.From("users").Limit(1)
	)

	doc.(*document).reflect()

	adapter.On("Query", query).Return(cur, errors.NewUnexpected("error")).Once()

	assert.NotNil(t, repo.One(&user, query))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_One_notFound(t *testing.T) {
	var (
		user    User
		doc     = newDocument(&user)
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		cur     = createCursor(0)
		query   = query.From("users").Limit(1)
	)

	doc.(*document).reflect()

	adapter.On("Query", query).Return(cur, nil).Once()

	err := repo.One(&user, query)
	assert.Equal(t, errors.New("no result found", "", errors.NotFound), err)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_All(t *testing.T) {
	var (
		users   []User
		collec  = newCollection(&users)
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		query   = query.From("users").Limit(1)
		cur     = createCursor(2)
	)

	collec.(*collection).reflect()

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.All(&users, query))
	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Insert(t *testing.T) {
	var (
		user      User
		doc       = newDocument(&user)
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
		cur     = createCursor(1)
	)

	doc.(*document).reflect()

	adapter.On("Insert", query.From("users"), changes).Return(1, nil).Once()
	adapter.On("Query", query.From("users").Where(where.Eq("id", 1)).Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Insert(&user, cbuilders...))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Insert_error(t *testing.T) {
	var (
		user      User
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
	)

	adapter.On("Insert", query.From("users"), changes).Return(0, errors.NewUnexpected("error")).Once()

	assert.NotNil(t, repo.Insert(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustInsert(&user, cbuilders...) })

	adapter.AssertExpectations(t)
}

func TestRepo_InsertAll(t *testing.T) {
	var (
		users   []User
		collec  = newCollection(&users)
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		changes = []change.Changes{
			change.Build(change.Set("name", "name1")),
			change.Build(change.Set("name", "name2")),
		}
		cur = createCursor(2)
	)

	collec.(*collection).reflect()

	adapter.On("InsertAll", query.From("users"), changes).Return([]interface{}{1, 2}, nil).Once()
	adapter.On("Query", query.From("users").Where(where.In("id", 1, 2))).Return(cur, nil).Once()

	assert.Nil(t, repo.InsertAll(&users, changes))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Update(t *testing.T) {
	var (
		user      = User{ID: 1}
		doc       = newDocument(&user)
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
		queries = query.From("users").Where(where.Eq("id", user.ID))
		cur     = createCursor(1)
	)

	doc.(*document).reflect()

	adapter.On("Update", queries, changes).Return(nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Update(&user, cbuilders...))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
}

func TestRepo_Update_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
	)

	assert.Nil(t, repo.Update(nil))
	assert.NotPanics(t, func() { repo.MustUpdate(nil) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_unchanged(t *testing.T) {
	var (
		user    = User{ID: 1}
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
	)

	assert.Nil(t, repo.Update(&user))
	assert.NotPanics(t, func() { repo.MustUpdate(&user) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_error(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
		queries = query.From("users").Where(where.Eq("id", user.ID))
	)

	adapter.On("Update", queries, changes).Return(errors.NewUnexpected("error")).Once()

	assert.NotNil(t, repo.Update(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustUpdate(&user, cbuilders...) })
	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_update(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = newDocument(transaction)
		buyerDoc    = newDocument(&transaction.Buyer)
		changes     = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users", where.Eq("id", 1))
		buyer, _ = changes.GetAssoc("Buyer")
		cur      = createCursor(1)
	)

	buyerDoc.(*document).reflect()

	adapter.On("Update", q, buyer[0]).Return(nil).Once()
	adapter.On("Query", q.Limit(1)).Return(cur, nil).Once()

	err := repo.upsertBelongsTo(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_updateError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = newDocument(transaction)
		changes     = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users", where.Eq("id", 1))
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.On("Update", q, buyer[0]).Return(errors.NewUnexpected("update error")).Once()

	err := repo.upsertBelongsTo(doc, &changes)
	assert.Equal(t, errors.NewUnexpected("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = newDocument(transaction)
		changes     = change.Build(
			change.Map{
				"Buyer": change.Map{
					"id":   2,
					"name": "buyer1",
					"age":  20,
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.upsertBelongsTo(doc, &changes)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_insertNew(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{}
		doc         = newDocument(transaction)
		buyerDoc    = newDocument(&transaction.Buyer)
		changes     = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users")
		buyer, _ = changes.GetAssoc("Buyer")
		cur      = createCursor(1)
	)

	buyerDoc.(*document).reflect()
	buyerDoc.(*document).initAssociations()

	adapter.On("Insert", q, buyer[0]).Return(1, nil).Once()
	adapter.On("Query", q.Where(where.Eq("id", 1)).Limit(1)).Return(cur, nil).Once()

	err := repo.upsertBelongsTo(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	ref, ok := changes.Get("user_id")
	assert.True(t, ok)
	assert.Equal(t, change.Set("user_id", 10), ref)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_insertNewError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{}
		doc         = newDocument(transaction)
		changes     = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users")
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.On("Insert", q, buyer[0]).Return(0, errors.NewUnexpected("insert error")).Once()

	err := repo.upsertBelongsTo(doc, &changes)
	assert.Equal(t, errors.NewUnexpected("insert error"), err)

	_, ok := changes.Get("user_id")
	assert.False(t, ok)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_notChanged(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = &Transaction{}
		doc         = newDocument(transaction)
		changes     = change.Build()
	)

	err := repo.upsertBelongsTo(doc, &changes)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_update(t *testing.T) {
	var (
		adapter    = &testAdapter{}
		repo       = Repo{adapter: adapter}
		user       = &User{ID: 1, Address: Address{ID: 2}}
		doc        = newDocument(user)
		addressDoc = newDocument(&user.Address)
		changes    = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q            = query.Build("addresses").Where(where.Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("Address")
		cur          = createCursor(1)
	)

	addressDoc.(*document).reflect()

	adapter.On("Update", q, addresses[0]).Return(nil).Once()
	adapter.On("Query", q.Limit(1)).Return(cur, nil).Once()

	err := repo.upsertHasOne(doc, &changes, nil)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertHasOne_updateError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{ID: 1, Address: Address{ID: 2}}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q            = query.Build("addresses").Where(where.Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("Address")
	)

	adapter.On("Update", q, addresses[0]).Return(errors.NewUnexpected("update error")).Once()

	err := repo.upsertHasOne(doc, &changes, nil)
	assert.Equal(t, errors.NewUnexpected("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{ID: 1, Address: Address{ID: 2}}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"id":     1,
					"street": "street1",
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.upsertHasOne(doc, &changes, nil)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_insertNew(t *testing.T) {
	var (
		adapter    = &testAdapter{}
		repo       = Repo{adapter: adapter}
		user       = &User{}
		doc        = newDocument(user)
		addressDoc = newDocument(&user.Address)
		changes    = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q       = query.Build("addresses")
		address = change.Build(change.Set("street", "street1"))
		cur     = createCursor(1)
	)

	addressDoc.(*document).reflect()
	addressDoc.(*document).initAssociations()

	// foreign value set after associations infered
	user.ID = 1
	address.SetValue("user_id", user.ID)

	adapter.On("Insert", q, address).Return(2, nil).Once()
	adapter.On("Query", q.Where(where.Eq("id", 2)).Limit(1)).Return(cur, nil).Once()

	err := repo.upsertHasOne(doc, &changes, nil)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertHasOne_insertNewError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q       = query.Build("addresses")
		address = change.Build(change.Set("street", "street1"))
	)

	// foreign value set after associations infered
	user.ID = 1
	address.SetValue("user_id", user.ID)

	adapter.On("Insert", q, address).Return(nil, errors.NewUnexpected("insert error")).Once()

	err := repo.upsertHasOne(doc, &changes, nil)
	assert.Equal(t, errors.NewUnexpected("insert error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasMany_insert(t *testing.T) {
	var (
		adapter           = &testAdapter{}
		repo              = Repo{adapter: adapter}
		user              = &User{ID: 1}
		doc               = newDocument(user)
		transactionCollec = newCollection(&user.Transactions)
		changes           = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item1",
					},
					{
						"item": "item2",
					},
				},
			},
		)
		q               = query.Build("transactions")
		transactions, _ = changes.GetAssoc("Transactions")
		cur             = createCursor(2)
	)

	transactionCollec.(*collection).reflect()

	adapter.On("InsertAll", q, transactions).Return(nil).Return([]interface{}{2, 3}, nil).Once()
	adapter.On("Query", q.Where(where.In("id", 2, 3))).Return(cur, nil).Once()

	err := repo.upsertHasMany(doc, &changes, user.ID, true)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(transactions))

	for i := range transactions {
		id, ok := transactions[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertHasMany_insertError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{ID: 1}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item1",
					},
					{
						"item": "item2",
					},
				},
			},
		)
		q               = query.Build("transactions")
		transactions, _ = changes.GetAssoc("Transactions")
		rerr            = errors.NewUnexpected("insert all error")
	)

	adapter.On("InsertAll", q, transactions).Return(nil).Return([]interface{}{}, rerr).Once()

	err := repo.upsertHasMany(doc, &changes, user.ID, true)
	assert.Equal(t, rerr, err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasMany_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
				{
					ID:   2,
					Item: "item2",
				},
			},
		}
		doc               = newDocument(user)
		transactionCollec = newCollection(&user.Transactions)
		changes           = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item3",
					},
					{
						"item": "item4",
					},
					{
						"item": "item5",
					},
				},
			},
		)
		q               = query.Build("transactions")
		transactions, _ = changes.GetAssoc("Transactions")
		cur             = createCursor(3)
	)

	transactionCollec.(*collection).reflect()

	adapter.On("Delete", q.Where(where.Eq("user_id", 1).AndIn("id", 1, 2))).Return(nil).Once()
	adapter.On("InsertAll", q, transactions).Return(nil).Return([]interface{}{3, 4, 5}, nil).Once()
	adapter.On("Query", q.Where(where.In("id", 3, 4, 5))).Return(cur, nil).Once()

	err := repo.upsertHasMany(doc, &changes, user.ID, false)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(transactions))

	for i := range transactions {
		id, ok := transactions[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertHasMany_updateEmptyAssoc(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{
			ID:           1,
			Transactions: []Transaction{},
		}
		doc               = newDocument(user)
		transactionCollec = newCollection(&user.Transactions)
		changes           = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item3",
					},
					{
						"item": "item4",
					},
					{
						"item": "item5",
					},
				},
			},
		)
		q               = query.Build("transactions")
		transactions, _ = changes.GetAssoc("Transactions")
		cur             = createCursor(3)
	)

	transactionCollec.(*collection).reflect()

	adapter.On("InsertAll", q, transactions).Return(nil).Return([]interface{}{3, 4, 5}, nil).Once()
	adapter.On("Query", q.Where(where.In("id", 3, 4, 5))).Return(cur, nil).Once()

	err := repo.upsertHasMany(doc, &changes, user.ID, false)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(transactions))

	for i := range transactions {
		id, ok := transactions[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_upsertHasMany_updateDeleteAllError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
				{
					ID:   2,
					Item: "item2",
				},
			},
		}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item3",
					},
				},
			},
		)
		q    = query.Build("transactions")
		rerr = errors.NewUnexpected("delete all error")
	)

	adapter.On("Delete", q.Where(where.Eq("user_id", 1).AndIn("id", 1, 2))).Return(rerr).Once()

	err := repo.upsertHasMany(doc, &changes, user.ID, false)
	assert.Equal(t, rerr, err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasMany_updateAssocNotLoaded(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{ID: 1}
		doc     = newDocument(user)
		changes = change.Build(
			change.Map{
				"Transactions": []change.Map{
					{
						"item": "item3",
					},
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.upsertHasMany(doc, &changes, user.ID, false)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_Delete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = &User{ID: 1}
	)

	adapter.On("Delete", query.From("users").Where(where.Eq("id", user.ID))).Return(nil).Once()

	assert.Nil(t, repo.Delete(user))

	adapter.AssertExpectations(t)
}

// func TestRepo_Delete_slice(t *testing.T) {
// 	var (
// 		adapter = &testAdapter{}
// 		repo    = Repo{adapter: adapter}
// 		users   = []User{
// 			{ID: 1},
// 			{ID: 2},
// 		}
// 	)

// 	adapter.
// 		On("Delete", query.From("users").Where(where.In("id", 1, 2))).Return(nil)

// 	assert.Nil(t, repo.Delete(users))
// 	assert.NotPanics(t, func() { repo.MustDelete(users) })
// 	adapter.AssertExpectations(t)
// }

// func TestRepo_Delete_emptySlice(t *testing.T) {
// 	var (
// 		adapter = &testAdapter{}
// 		repo    = Repo{adapter: adapter}
// 		users   = []User{}
// 	)

// 	assert.Nil(t, repo.Delete(users))
// 	assert.NotPanics(t, func() { repo.MustDelete(users) })
// 	adapter.AssertExpectations(t)
// }

func TestRepo_DeleteAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		queries = query.From("logs").Where(where.Eq("user_id", 1))
	)

	adapter.On("Delete", query.From("logs").Where(where.Eq("user_id", 1))).Return(nil).Once()

	assert.Nil(t, repo.DeleteAll(queries))

	adapter.AssertExpectations(t)
}

func TestRepo_Preload_hasOne(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = User{ID: 10}
		address = Address{ID: 100, UserID: &user.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.SetScan(2, address.ID, *address.UserID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&user, "Address"))
	assert.Equal(t, address, user.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_sliceHasOne(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		users     = []User{{ID: 10}, {ID: 20}}
		addresses = []Address{
			{ID: 100, UserID: &users[0].ID},
			{ID: 200, UserID: &users[1].ID},
		}
		cur = &testCursor{}
	)

	// one of these, because of map ordering
	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.SetScan(2, addresses[0].ID, *addresses[0].UserID)
	cur.SetScan(2, addresses[1].ID, *addresses[1].UserID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&users, "Address"))
	assert.Equal(t, addresses[0], users[0].Address)
	assert.Equal(t, addresses[1], users[1].Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_nestedHasOne(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		transaction = Transaction{
			Buyer: User{ID: 10},
		}
		address = Address{ID: 100, UserID: &transaction.Buyer.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.SetScan(2, address.ID, *address.UserID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transaction, "Buyer.Address"))
	assert.Equal(t, address, transaction.Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_sliceNestedHasOne(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = Repo{adapter: adapter}
		transactions = []Transaction{
			{Buyer: User{ID: 10}},
			{Buyer: User{ID: 20}},
		}
		addresses = []Address{
			{ID: 100, UserID: &transactions[0].Buyer.ID},
			{ID: 200, UserID: &transactions[1].Buyer.ID},
		}
		cur = &testCursor{}
	)

	// one of these, because of map ordering
	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", query.From("addresses").Where(where.In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.SetScan(2, addresses[0].ID, *addresses[0].UserID)
	cur.SetScan(2, addresses[1].ID, *addresses[1].UserID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transactions, "Buyer.Address"))
	assert.Equal(t, addresses[0], transactions[0].Buyer.Address)
	assert.Equal(t, addresses[1], transactions[1].Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_hasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = Repo{adapter: adapter}
		user         = User{ID: 10}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.SetScan(2, transactions[0].ID, transactions[0].BuyerID)
	cur.SetScan(2, transactions[1].ID, transactions[1].BuyerID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&user, "Transactions"))
	assert.Equal(t, transactions, user.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_sliceHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = Repo{adapter: adapter}
		users        = []User{{ID: 10}, {ID: 20}}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 20},
			{ID: 20, BuyerID: 20},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(4)
	cur.SetScan(2, transactions[0].ID, transactions[0].BuyerID)
	cur.SetScan(2, transactions[1].ID, transactions[1].BuyerID)
	cur.SetScan(2, transactions[2].ID, transactions[2].BuyerID)
	cur.SetScan(2, transactions[3].ID, transactions[3].BuyerID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&users, "Transactions"))
	assert.Equal(t, transactions[:2], users[0].Transactions)
	assert.Equal(t, transactions[2:], users[1].Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_nestedHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = Repo{adapter: adapter}
		address      = Address{User: &User{ID: 10}}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
		}

		cur = &testCursor{}
	)

	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.SetScan(2, transactions[0].ID, transactions[0].BuyerID)
	cur.SetScan(2, transactions[1].ID, transactions[1].BuyerID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&address, "User.Transactions"))
	assert.Equal(t, transactions, address.User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_nestedNullHasMany(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		address = Address{User: nil}
	)

	assert.Nil(t, repo.Preload(&address, "User.Transactions"))

	adapter.AssertExpectations(t)
}

func TestRepo_Preload_nestedSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		addresses = []Address{
			{User: &User{ID: 10}},
			{User: &User{ID: 20}},
		}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 20},
			{ID: 20, BuyerID: 20},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(4)
	cur.SetScan(2, transactions[0].ID, transactions[0].BuyerID)
	cur.SetScan(2, transactions[1].ID, transactions[1].BuyerID)
	cur.SetScan(2, transactions[2].ID, transactions[2].BuyerID)
	cur.SetScan(2, transactions[3].ID, transactions[3].BuyerID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&addresses, "User.Transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[1].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_nestedNullSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = Repo{adapter: adapter}
		addresses = []Address{
			{User: &User{ID: 10}},
			{User: nil},
			{User: &User{ID: 15}},
		}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 15},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 10, 15))).Return(cur, nil).Maybe()
	adapter.On("Query", query.From("transactions").Where(where.In("user_id", 15, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(3)
	cur.SetScan(2, transactions[0].ID, transactions[0].BuyerID)
	cur.SetScan(2, transactions[1].ID, transactions[1].BuyerID)
	cur.SetScan(2, transactions[2].ID, transactions[2].BuyerID)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&addresses, "User.Transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, []Transaction(nil), addresses[1].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[2].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_belongsTo(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = Repo{adapter: adapter}
		user        = User{ID: 10, Name: "Del Piero"}
		transaction = Transaction{BuyerID: 10}
		cur         = &testCursor{}
	)

	adapter.On("Query", query.From("users").Where(where.In("id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.SetScan(2, user.ID, user.Name)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transaction, "Buyer"))
	assert.Equal(t, user, transaction.Buyer)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Preload_belongsToPtrKey(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = Repo{adapter: adapter}
		user    = User{ID: 10, Name: "Del Piero"}
		address = Address{UserID: &user.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", query.From("users").Where(where.In("id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.SetScan(2, user.ID, user.Name)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&address, "User"))
	assert.Equal(t, user, *address.User)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepo_Transaction(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).On("Commit").Return(nil).Once()

	repo := Repo{adapter: adapter}

	err := repo.Transaction(func(repo Repo) error {
		assert.True(t, repo.inTransaction)
		return nil
	})

	assert.False(t, repo.inTransaction)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_beginError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(errors.NewUnexpected("error")).Once()

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_commitError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Commit").Return(errors.NewUnexpected("error")).Once()

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_returnErrorAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return errors.NewUnexpected("error")
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_panicWithKnownErrorAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		panic(errors.New("error", "", errors.NotFound))
	})

	assert.Equal(t, errors.New("error", "", errors.NotFound), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_panicAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Panics(t, func() {
		Repo{adapter: adapter}.Transaction(func(r Repo) error {
			// doing good things
			panic(errors.NewUnexpected("error"))
		})
	})

	adapter.AssertExpectations(t)
}
