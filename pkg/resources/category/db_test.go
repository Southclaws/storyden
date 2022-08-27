package category_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/utils/integration"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/seed"
)

func TestCreateCategory(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo category.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			c, err := repo.CreateCategory(ctx, "test", "desc", "ffffff", 1, false)
			r.NoError(err)
			r.NotNil(c)

			a.Equal("test", c.Name)
			a.Equal("desc", c.Description)
			a.Equal("ffffff", c.Colour)
			a.Equal(1, c.Sort)
			a.Equal(false, c.Admin)
			a.Len(c.Recent, 0)
			a.Equal(0, c.PostCount)
		}),
	)
}

func TestGetCategories(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo category.Repository,
			d *model.Client,
		) {
			a := assert.New(t)
			r := require.New(t)

			// No results means a non-nil, empty list.
			none, err := repo.GetCategories(ctx, false)
			r.NoError(err)
			a.NotNil(none)

			// Create two categories
			c1, err := d.Category.Create().SetName("cat5").SetSort(6).Save(ctx)
			r.NoError(err)
			c2, err := d.Category.Create().SetName("cat6").SetSort(7).Save(ctx)
			r.NoError(err)

			// Create four posts, the first two are tagged with tag1 the third one is
			// tagged with tag2 and the fourth is tagged with both.
			p0, err := d.Post.Create().
				SetBody("").
				SetShort("").
				SetFirst(true).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()).
				SetAuthorID(uuid.UUID(seed.Account_002.ID)).
				SetCategory(c1).
				Save(ctx)
			r.NoError(err)
			p1, err := d.Post.Create().
				SetBody("").
				SetShort("").
				SetFirst(true).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()).
				SetAuthorID(uuid.UUID(seed.Account_002.ID)).
				SetCategory(c1).
				Save(ctx)
			r.NoError(err)
			p2, err := d.Post.Create().
				SetBody("").
				SetShort("").
				SetFirst(true).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()).
				SetAuthorID(uuid.UUID(seed.Account_002.ID)).
				SetCategory(c2).
				Save(ctx)
			r.NoError(err)

			// Searching for the prefix "ta" should get all our tags
			categories, err := repo.GetCategories(ctx, false)
			r.NoError(err)
			r.NotNil(categories)

			// Seed categories plus newly created
			r.Len(categories, 6)

			cat1 := categories[0]
			cat2 := categories[1]
			cat5 := categories[4]
			cat6 := categories[5]

			a.Equal(seed.Category_01_General.Name, cat1.Name)
			a.Equal(seed.Category_02_Photos.Name, cat2.Name)
			a.Equal("cat5", cat5.Name)
			a.Equal("cat6", cat6.Name)

			a.Equal(2, cat5.PostCount)
			a.Equal(1, cat6.PostCount)

			a.Len(cat5.Recent, 2)
			a.Len(cat6.Recent, 1)

			// recent posts list is correct
			c1posts := ids(cat5.Recent)
			a.Contains(c1posts, p0.ID)
			a.Contains(c1posts, p1.ID)
			a.NotContains(c1posts, p2.ID)

			c2posts := ids(cat6.Recent)
			a.Contains(c2posts, p2.ID)
			a.NotContains(c2posts, p0.ID)
			a.NotContains(c2posts, p1.ID)
		}),
	)
}

func ids(c []category.PostMeta) []uuid.UUID {
	return lo.Map(c, func(p category.PostMeta, _ int) uuid.UUID { return p.PostID })
}
