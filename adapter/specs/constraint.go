package specs

import (
	"testing"

	"github.com/go-rel/rel"
)

func createExtra(repo rel.Repository, slug string) Extra {
	var user User
	repo.MustInsert(ctx, &user)

	extra := Extra{Slug: &slug, UserID: user.ID}
	repo.MustInsert(ctx, &extra)
	return extra
}

// UniqueConstraintOnInsert tests unique constraint specifications on insert.
func UniqueConstraintOnInsert(t *testing.T, repo rel.Repository) {
	var (
		existing = createExtra(repo, "unique-insert")
		err      = repo.Insert(ctx, &Extra{Slug: existing.Slug})
	)

	assertConstraint(t, err, rel.UniqueConstraint, "slug")
}

// UniqueConstraintOnUpdate tests unique constraint specifications on insert.
func UniqueConstraintOnUpdate(t *testing.T, repo rel.Repository) {
	var (
		record   = createExtra(repo, "unique-record")
		existing = createExtra(repo, "unique-update-existing")
		err      = repo.Update(ctx, &Extra{ID: record.ID, Slug: existing.Slug})
	)

	assertConstraint(t, err, rel.UniqueConstraint, "slug")
}

// ForeignKeyConstraintOnInsert tests foreign key constraint specifications on insert.
func ForeignKeyConstraintOnInsert(t *testing.T, repo rel.Repository) {
	var (
		err = repo.Insert(ctx, &Extra{UserID: 1000})
	)

	assertConstraint(t, err, rel.ForeignKeyConstraint, "user_id")
}

// ForeignKeyConstraintOnUpdate tests foreign key constraint specifications on update.
func ForeignKeyConstraintOnUpdate(t *testing.T, repo rel.Repository) {
	var (
		record = createExtra(repo, "fk-slug")
	)

	record.UserID = 1000
	err := repo.Update(ctx, &record)
	assertConstraint(t, err, rel.ForeignKeyConstraint, "user_id")
}

// CheckConstraintOnInsert tests foreign key constraint specifications on insert.
func CheckConstraintOnInsert(t *testing.T, repo rel.Repository) {
	var (
		err = repo.Insert(ctx, &Extra{Score: 150})
	)

	assertConstraint(t, err, rel.CheckConstraint, "score")
}

// CheckConstraintOnUpdate tests foreign key constraint specifications.
func CheckConstraintOnUpdate(t *testing.T, repo rel.Repository) {
	var (
		record = createExtra(repo, "check-slug")
	)

	// updating
	record.Score = 150
	err := repo.Update(ctx, &record)
	assertConstraint(t, err, rel.CheckConstraint, "score")
}
