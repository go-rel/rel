package specs

import (
	"testing"

	"github.com/Fs02/rel"
)

func createExtra(repo rel.Repository, slug string) Extra {
	var user User
	repo.MustInsert(&user)

	extra := Extra{Slug: &slug, UserID: user.ID}
	repo.MustInsert(&extra)
	return extra
}

// UniqueConstraint tests unique constraint specifications.
func UniqueConstraint(t *testing.T, repo rel.Repository) {
	var (
		extra1 = createExtra(repo, "unique-slug1")
		extra2 = createExtra(repo, "unique-slug2")
	)

	t.Run("UniqueConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{Slug: extra1.Slug})
		assertConstraint(t, err, rel.UniqueConstraint, "slug")

		// updating
		err = repo.Update(&Extra{ID: extra2.ID, Slug: extra1.Slug})
		assertConstraint(t, err, rel.UniqueConstraint, "slug")
	})
}

// ForeignKeyConstraint tests foreign key constraint specifications.
func ForeignKeyConstraint(t *testing.T, repo rel.Repository) {
	var (
		extra = createExtra(repo, "fk-slug")
	)

	t.Run("ForeignKeyConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{UserID: 1000})
		assertConstraint(t, err, rel.ForeignKeyConstraint, "user_id")

		// updating
		extra.UserID = 1000
		err = repo.Update(&extra)
		assertConstraint(t, err, rel.ForeignKeyConstraint, "user_id")
	})
}

// CheckConstraint tests foreign key constraint specifications.
func CheckConstraint(t *testing.T, repo rel.Repository) {
	var (
		extra = createExtra(repo, "check-slug")
	)

	t.Run("CheckConstraint", func(t *testing.T) {
		// inserting
		err := repo.Insert(&Extra{Score: 150})
		assertConstraint(t, err, rel.CheckConstraint, "score")

		// updating
		extra.Score = 150
		err = repo.Update(&extra)
		assertConstraint(t, err, rel.CheckConstraint, "score")
	})
}
