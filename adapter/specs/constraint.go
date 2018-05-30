package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/errors"
)

// UniqueConstraint tests unique constraint specifications.
func UniqueConstraint(t *testing.T, repo grimoire.Repo) {
	extra1 := Extra{}
	extra2 := Extra{}
	repo.From(extras).Set("slug", "slug1").MustInsert(&extra1)
	repo.From(extras).Set("slug", "slug2").MustInsert(&extra2)

	t.Run("UniqueConstraint", func(t *testing.T) {
		// inserting
		err := repo.From(extras).Set("slug", extra1.Slug).Insert(nil)
		assertConstraint(t, err, errors.UniqueConstraint, "slug")

		// updating
		err = repo.From(extras).Find(extra2.ID).Set("slug", extra1.Slug).Update(nil)
		assertConstraint(t, err, errors.UniqueConstraint, "slug")
	})
}

// ForeignKeyConstraint tests foreign key constraint specifications.
func ForeignKeyConstraint(t *testing.T, repo grimoire.Repo) {
	extra := Extra{}
	repo.From(extras).MustSave(&extra)

	t.Run("ForeignKeyConstraint", func(t *testing.T) {
		// inserting
		err := repo.From(extras).Set("user_id", 1000).Insert(nil)
		assertConstraint(t, err, errors.ForeignKeyConstraint, "user_id")

		// updating
		err = repo.From(extras).Find(extra.ID).Set("user_id", 1000).Update(nil)
		assertConstraint(t, err, errors.ForeignKeyConstraint, "user_id")
	})
}

// CheckConstraint tests foreign key constraint specifications.
func CheckConstraint(t *testing.T, repo grimoire.Repo) {
	extra := Extra{}
	repo.From(extras).MustSave(&extra)

	t.Run("CheckConstraint", func(t *testing.T) {
		// inserting
		err := repo.From(extras).Set("score", 150).Insert(nil)
		assertConstraint(t, err, errors.CheckConstraint, "score")

		// updating
		err = repo.From(extras).Find(extra.ID).Set("score", 150).Update(nil)
		assertConstraint(t, err, errors.CheckConstraint, "score")
	})
}
