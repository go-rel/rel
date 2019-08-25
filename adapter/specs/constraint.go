package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
)

// UniqueConstraint tests unique constraint specifications.
func UniqueConstraint(t *testing.T, repo grimoire.Repo) {
	var (
		slug1  = "slug1"
		slug2  = "slug2"
		extra1 = Extra{Slug: &slug1}
		extra2 = Extra{Slug: &slug2}
	)

	repo.MustInsert(&extra1)
	repo.MustInsert(&extra2)

	t.Run("UniqueConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{Slug: extra1.Slug})
		assertConstraint(t, err, grimoire.UniqueConstraint, "slug")

		// updating
		err = repo.Update(&Extra{ID: extra2.ID, Slug: extra1.Slug})
		assertConstraint(t, err, grimoire.UniqueConstraint, "slug")
	})
}

// ForeignKeyConstraint tests foreign key constraint specifications.
func ForeignKeyConstraint(t *testing.T, repo grimoire.Repo) {
	var (
		extra Extra
	)

	repo.MustInsert(&extra)

	t.Run("ForeignKeyConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{UserID: 1000})
		assertConstraint(t, err, grimoire.ForeignKeyConstraint, "user_id")

		// updating
		extra.UserID = 1000
		err = repo.Update(&extra)
		assertConstraint(t, err, grimoire.ForeignKeyConstraint, "user_id")
	})
}

// CheckConstraint tests foreign key constraint specifications.
func CheckConstraint(t *testing.T, repo grimoire.Repo) {
	var (
		extra Extra
	)

	repo.MustInsert(&extra)

	t.Run("CheckConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{Score: 150})
		assertConstraint(t, err, grimoire.CheckConstraint, "score")

		// updating
		extra.Score = 150
		err = repo.Update(&extra)
		assertConstraint(t, err, grimoire.CheckConstraint, "score")
	})
}
