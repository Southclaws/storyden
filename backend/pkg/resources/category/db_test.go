package category

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

func TestCreateCategory(t *testing.T) {
	// a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	d := db.TestDB(t)
	repo := New(d)

	c, err := repo.CreateCategory(ctx, "test", "desc", "fffff", 1, false)
	r.NoError(err)
	r.NotNil(c)
}

func TestGetCategories(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	d := db.TestDB(t)
	user.Seed(user.New(d))
	repo := New(d)

	// No results means a non-nil, empty list.
	none, err := repo.GetCategories(ctx, false)
	r.NoError(err)
	a.NotNil(none)
	a.Empty(none)

	// Create two categories
	c1, err := d.Category.Create().SetName("cat1").SetSort(0).Save(ctx)
	r.NoError(err)
	c2, err := d.Category.Create().SetName("cat2").SetSort(1).Save(ctx)
	r.NoError(err)

	// Create four posts, the first two are tagged with tag1 the third one is
	// tagged with tag2 and the fourth is tagged with both.
	p0, err := d.Post.Create().
		SetBody("").
		SetShort("").
		SetFirst(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAuthorID(uuid.UUID(user.SeedUser_02_User.ID)).
		SetCategory(c1).
		Save(ctx)
	r.NoError(err)
	p1, err := d.Post.Create().
		SetBody("").
		SetShort("").
		SetFirst(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAuthorID(uuid.UUID(user.SeedUser_02_User.ID)).
		SetCategory(c1).
		Save(ctx)
	r.NoError(err)
	p2, err := d.Post.Create().
		SetBody("").
		SetShort("").
		SetFirst(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAuthorID(uuid.UUID(user.SeedUser_02_User.ID)).
		SetCategory(c2).
		Save(ctx)
	r.NoError(err)

	// Searching for the prefix "ta" should get all our tags
	categories, err := repo.GetCategories(ctx, false)
	r.NoError(err)
	r.NotNil(categories)

	// All 2 of them
	r.Len(categories, 2)

	cat1 := categories[0]
	cat2 := categories[1]

	a.Equal(cat1.Name, "cat1")
	a.Equal(cat2.Name, "cat2")

	a.Equal(cat1.PostCount, 2)
	a.Equal(cat2.PostCount, 1)

	a.Len(cat1.Recent, 2)
	a.Len(cat2.Recent, 1)

	// recent posts list is correct
	c1posts := ids(cat1.Recent)
	a.Contains(c1posts, p0.ID)
	a.Contains(c1posts, p1.ID)
	a.NotContains(c1posts, p2.ID)

	c2posts := ids(cat2.Recent)
	a.Contains(c2posts, p2.ID)
	a.NotContains(c2posts, p0.ID)
	a.NotContains(c2posts, p1.ID)
}

func ids(c []PostMeta) []uuid.UUID {
	return lo.Map(c, func(p PostMeta, _ int) uuid.UUID { return p.PostID })
}
