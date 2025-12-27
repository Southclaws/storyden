package thread

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

type Thread struct {
	post.Post

	Title       string
	Slug        string
	Short       string
	Pinned      bool
	LastReplyAt opt.Optional[time.Time]

	ReadStatus  opt.Optional[post.ReadStatus]
	ReplyStatus post.ReplyStatus
	Replies     pagination.Result[*reply.Reply]
	Category    opt.Optional[category.Category]
	Tags        tag_ref.Tags
	Related     datagraph.ItemList
}

type ThreadRef struct {
	post.Post
	Title       string
	Slug        string
	Short       string
	Pinned      bool
	LastReplyAt opt.Optional[time.Time]
}

func (*Thread) GetResourceName() string { return "thread" }

func (t *Thread) GetKind() datagraph.Kind { return datagraph.KindThread }
func (t *Thread) GetName() string         { return t.Title }
func (t *Thread) GetSlug() string         { return t.Slug }
func (t *Thread) GetDesc() string         { return t.Short }
func (t *Thread) GetCreated() time.Time   { return t.CreatedAt }
func (t *Thread) GetUpdated() time.Time   { return t.UpdatedAt }
func (t *Thread) GetAuthor() xid.ID       { return xid.ID(t.Author.ID) }
func (t *Thread) GetCategory() xid.ID {
	if cat, ok := t.Category.Get(); ok {
		return xid.ID(cat.ID)
	}
	return xid.NilID()
}

func (t *Thread) GetTags() []string {
	tags := make([]string, len(t.Tags))
	for i, tag := range t.Tags {
		tags[i] = tag.Name.String()
	}
	return tags
}

func Map(m *ent.Post) (*Thread, error) {
	category := opt.Map(opt.NewPtr(m.Edges.Category), func(in ent.Category) category.Category {
		return *category.FromModel(&in)
	})

	link := opt.Map(opt.NewPtr(m.Edges.Link), func(in ent.Link) link_ref.LinkRef {
		return *link_ref.Map(&in)
	})

	content, err := datagraph.NewRichText(m.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.MapRef(m.Edges.Author)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	tags := dt.Map(m.Edges.Tags, tag_ref.Map(nil))

	return &Thread{
		Post: post.Post{
			ID: post.ID(m.ID),

			Content:    content,
			Author:     *pro,
			Assets:     dt.Map(m.Edges.Assets, asset.Map),
			WebLink:    link,
			Meta:       m.Metadata,
			Visibility: visibility.NewVisibilityFromEnt(m.Visibility),

			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: opt.NewPtr(m.DeletedAt),
		},

		Title:       m.Title,
		Slug:        m.Slug,
		Short:       m.Short,
		Pinned:      m.Pinned,
		LastReplyAt: opt.New(m.LastReplyAt),

		Category: category,
		Tags:     tags,
	}, nil
}

func Mapper(
	am account.Lookup,
	rr post.ReadStateMap,
	ls post.PostLikesMap,
	cs collection_item_status.CollectionStatusMap,
	rs post.PostRepliesMap,
	rl reaction.Lookup,
) func(m *ent.Post) (*Thread, error) {
	return func(m *ent.Post) (*Thread, error) {
		category := opt.Map(opt.NewPtr(m.Edges.Category), func(in ent.Category) category.Category {
			return *category.FromModel(&in)
		})

		link := opt.Map(opt.NewPtr(m.Edges.Link), func(in ent.Link) link_ref.LinkRef {
			return *link_ref.Map(&in)
		})

		content, err := datagraph.NewRichText(m.Body)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		var pro *profile.Ref
		authorEdge := am[m.AccountPosts]
		if authorEdge != nil {
			pro, err = profile.MapRef(authorEdge)
			if err != nil {
				return nil, fault.Wrap(err)
			}
		} else {
			pro, err = profile.MapRef(m.Edges.Author)
			if err != nil {
				return nil, fault.Wrap(err)
			}
		}

		reacts := rl[xid.ID(m.ID)]

		return &Thread{
			Post: post.Post{
				ID: post.ID(m.ID),

				Content:     content,
				Author:      *pro,
				Likes:       ls.Status(m.ID),
				Collections: cs.Status(m.ID),
				Reacts:      reacts,
				WebLink:     link,
				Meta:        m.Metadata,
				Visibility:  visibility.NewVisibilityFromEnt(m.Visibility),

				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
				DeletedAt: opt.NewPtr(m.DeletedAt),
			},

			Title:  m.Title,
			Slug:   m.Slug,
			Short:  m.Short,
			Pinned: m.Pinned,
			// Only populate the last-reply-at if there are replies.
			LastReplyAt: opt.NewSafe(m.LastReplyAt, rs.Status(m.ID).Count > 0),

			ReadStatus:  rr.Status(m.ID),
			ReplyStatus: rs.Status(m.ID),
			Category:    category,
		}, nil
	}
}

func MapRef(m *ent.Post) *ThreadRef {
	return &ThreadRef{
		Post: post.Post{
			ID: post.ID(m.ID),

			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: opt.NewPtr(m.DeletedAt),
		},

		Title:       m.Title,
		Slug:        m.Slug,
		Short:       m.Short,
		Pinned:      m.Pinned,
		LastReplyAt: opt.New(m.LastReplyAt),
	}
}

func ItemRef(t *ent.Post) (datagraph.Item, error) {
	content, err := datagraph.NewRichText(t.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Thread{
		Post: post.Post{
			ID:      post.ID(t.ID),
			Content: content,
			Meta:    t.Metadata,
			Author: profile.Ref{
				ID: account.AccountID(t.AccountPosts),
			},
			Visibility: visibility.NewVisibilityFromEnt(t.Visibility),
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
			DeletedAt:  opt.NewPtr(t.DeletedAt),
		},

		Title: t.Title,
		Slug:  t.Slug,
		Category: opt.New(category.Category{
			ID: category.CategoryID(t.CategoryID),
		}),
		Short: t.Short,
		Tags:  dt.Map(t.Edges.Tags, tag_ref.Map(nil)),
	}, nil
}
