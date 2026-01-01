package thread_writer

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/category"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

type Option func(*ent.PostMutation)

func WithID(id post.ID) Option {
	return func(m *ent.PostMutation) {
		m.SetID(xid.ID(id))
	}
}

func WithIndexed() Option {
	return func(m *ent.PostMutation) {
		m.SetIndexedAt(time.Now())
	}
}

func WithTitle(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetTitle(v)
	}
}

func WithContent(v datagraph.Content) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v.HTML())
		pm.SetShort(v.Short())
	}
}

func WithCategory(v xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetCategoryID(v)
	}
}

func WithVisibility(v visibility.Visibility) Option {
	return func(pm *ent.PostMutation) {
		pm.SetVisibility(ent_post.Visibility(v.String()))
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}

func WithPinned(v int) Option {
	return func(m *ent.PostMutation) {
		m.SetPinned(v)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.PostMutation) {
		m.AddAssetIDs(a...)
	}
}

func WithLink(id xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetLinkID(id)
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}

func WithTagsAdd(refs ...tag_ref.ID) Option {
	ids := dt.Map(refs, func(i tag_ref.ID) xid.ID { return xid.ID(i) })
	return func(c *ent.PostMutation) {
		c.AddTagIDs(ids...)
	}
}

func WithTagsRemove(refs ...tag_ref.ID) Option {
	ids := dt.Map(refs, func(i tag_ref.ID) xid.ID { return xid.ID(i) })
	return func(c *ent.PostMutation) {
		c.RemoveTagIDs(ids...)
	}
}

func (d *Writer) Create(
	ctx context.Context,
	title string,
	authorID account.AccountID,
	opts ...Option,
) (*thread.Thread, error) {
	create := d.db.Post.Create()
	mutate := create.Mutation()
	mutate.SetUpdatedAt(time.Now())

	for _, fn := range opts {
		fn(mutate)
	}

	// If a category was specified, check if it exists first.
	if categoryID, ok := mutate.CategoryID(); ok {
		exists, err := d.db.Category.Query().Where(category.ID(xid.ID(categoryID))).Exist(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
		}
		if !exists {
			return nil, fault.Wrap(fault.New("category not found"),
				fctx.With(ctx),
				ftag.With(ftag.InvalidArgument),
				fmsg.WithDesc("category not found",
					"The specified category was not found."))
		}
	}

	mutate.SetTitle(title)
	mutate.SetLastReplyAt(time.Now())
	mutate.SetAuthorID(xid.ID(authorID))
	mutate.SetTitle(title)

	p, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// Update the slug so it includes the ID for uniqueness.

	_, err = d.db.Post.
		UpdateOneID(p.ID).
		SetSlug(fmt.Sprintf("%s-%s", p.ID, mark.Slugify(title))).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// Finally, query the created thread with related entities.

	p, err = d.db.Post.
		Query().
		Where(ent_post.IDEQ(p.ID)).
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return thread.Map(p)
}

func (d *Writer) Update(ctx context.Context, id post.ID, opts ...Option) (*thread.Thread, error) {
	update := d.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	// Only set the updated_at field if not changing the indexed_at field.
	if _, set := mutate.IndexedAt(); !set {
		mutate.SetUpdatedAt(time.Now())
	}

	err := update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return thread.Map(p)
}

func (d *Writer) Delete(ctx context.Context, id post.ID) error {
	err := d.db.Post.
		UpdateOneID(xid.ID(id)).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread root post"))
	}

	err = d.db.Post.
		Update().
		Where(ent_post.RootPostID(xid.ID(id))).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread posts"))
	}

	return nil
}
