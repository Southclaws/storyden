package collection_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestCollections(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			d *ent.Client,

			pr thread.Repository,

			repo collection.Repository,
		) {
			a := assert.New(t)
			r := require.New(t)

			acc := seed.Account_002_Frigg
			cat := seed.Category_01_General

			p0, err := pr.Create(ctx, "p0", "p0body", acc.ID, cat.ID, nil)
			r.NoError(err)
			p1, err := pr.Create(ctx, "p1", "p1body", acc.ID, cat.ID, nil)
			r.NoError(err)
			p2, err := pr.Create(ctx, "p2", "p2body", acc.ID, cat.ID, nil)
			r.NoError(err)

			coll, err := repo.Create(ctx, acc.ID, "test", "desc")
			r.NoError(err)
			r.NotNil(coll)

			a.Equal("test", coll.Name)
			a.Equal("desc", coll.Description)
			a.Empty(coll.Items)

			none, err := repo.List(ctx)
			r.NoError(err)
			a.NotNil(none)

			_, err = repo.Update(ctx, coll.ID, collection.WithPostAdd(p0.ID))
			r.NoError(err)
			_, err = repo.Update(ctx, coll.ID, collection.WithPostAdd(p1.ID))
			r.NoError(err)

			got, err := repo.Get(ctx, coll.ID)
			r.NoError(err)
			r.NotNil(got)

			r.Len(got.Items, 2)
			a.Equal(acc.Name, got.Items[0].Author.Name)

			got, err = repo.Update(ctx, coll.ID, collection.WithPostAdd(p2.ID))
			r.NoError(err)

			r.NotNil(got)
			r.Len(got.Items, 3)

			got, err = repo.Update(ctx, coll.ID, collection.WithPostRemove(p0.ID))
			r.NoError(err)

			r.NotNil(got)
			r.Len(got.Items, 2)

			err = repo.Delete(ctx, coll.ID)
			r.NoError(err)

			got, err = repo.Get(ctx, coll.ID)
			r.Error(err)
			r.Nil(got)
		}),
	)
}

func ids(c []category.PostMeta) []xid.ID {
	return lo.Map(c, func(p category.PostMeta, _ int) xid.ID { return p.PostID })
}
