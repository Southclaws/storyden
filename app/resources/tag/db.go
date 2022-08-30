package tag

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/errtag"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/tag"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) GetTags(ctx context.Context, query string) ([]Tag, error) {
	tags := []Tag{}

	err := d.db.Tag.Query().Modify(func(s *sql.Selector) {
		s.
			Select(
				sql.As(sql.Table("t").C("id"), "id"),
				sql.As("name", "name"),
				sql.As(sql.Count("*"), "posts"),
			).
			From(sql.Table(tag.PostsTable)).
			Where(sql.HasPrefix(sql.Table("t").C("name"), query)).
			Join(sql.Table(tag.Table).As("t")).On(sql.Table("t").C(tag.FieldID), "tag_id").
			Join(sql.Table(post.Table).As("p")).On(sql.Table("p").C(post.FieldID), "post_id").
			GroupBy(sql.Table("t").C("id")).
			OrderBy(sql.Desc("posts"))
	}).Scan(ctx, &tags)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return tags, nil
}
