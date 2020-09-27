package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestPreloadBelongsTo(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [preload-belongs-to]
	user := User{ID: 1, Name: "Nabe"}
	repo.ExpectPreload("buyer").Result(user)
	/// [preload-belongs-to]

	assert.Nil(t, PreloadBelongsTo(ctx, repo))
	repo.AssertExpectations(t)
}

func TestPreloadHasOne(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [preload-has-one]
	address := Address{ID: 1, City: "Nazarick"}
	repo.ExpectPreload("address").Result(address)
	/// [preload-has-one]

	assert.Nil(t, PreloadHasOne(ctx, repo))
	repo.AssertExpectations(t)
}

func TestPreloadHasMany(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [preload-has-many]
	transactions := []Transaction{
		{ID: 1, Item: "Avarice and Generosity", Status: "paid"},
	}
	repo.ExpectPreload("transactions").Result(transactions)
	/// [preload-has-many]

	assert.Nil(t, PreloadHasMany(ctx, repo))
	repo.AssertExpectations(t)
}

func TestPreloadHasManyFilter(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [preload-has-many-filter]
	transactions := []Transaction{
		{ID: 1, Item: "Avarice and Generosity", Status: "paid"},
	}
	repo.ExpectPreload("transactions", where.Eq("status", "paid")).Result(transactions)
	/// [preload-has-many-filter]

	assert.Nil(t, PreloadHasManyFilter(ctx, repo))
	repo.AssertExpectations(t)
}

func TestPreloadNested(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [preload-nested]
	userID := 1
	addresses := []Address{{ID: 1, City: "Nazarick", UserID: &userID}}
	repo.ExpectPreload("buyer.address").Result(addresses)
	/// [preload-nested]

	assert.Nil(t, PreloadNested(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsertAssociation(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-association]
	repo.ExpectInsert().ForType("main.User")
	/// [insert-association]

	assert.Nil(t, InsertAssociation(ctx, repo))
	repo.AssertExpectations(t)
}

func TestUpdateAssociation(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [update-association]
	repo.ExpectUpdate().ForType("main.User")
	/// [update-association]

	assert.Nil(t, UpdateAssociation(ctx, repo))
	repo.AssertExpectations(t)
}

func TestUpdateAssociationWithMap(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [update-association-with-map]
	mutation := rel.Map{
		"address": rel.Map{
			"city": "bandung",
		},
	}

	// Update address record with id 1, only set city to bandung.
	repo.ExpectUpdate(mutation).ForType("main.User")
	/// [update-association-with-map]

	assert.Nil(t, UpdateAssociationWithMap(ctx, repo))
	repo.AssertExpectations(t)
}
