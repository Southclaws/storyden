package reply

import (
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type Reply struct {
	post.Post

	RootPostID      post.ID
	RootThreadMark  string
	RootThreadTitle string
	RootAuthor      profile.Ref
	Slug            string // The root slug with the post ID as a #fragment
	ReplyTo         opt.Optional[post.ID]
}

func (*Reply) GetResourceName() string { return "post" }

func (r *Reply) GetName() string {
	if xid.ID(r.RootPostID).IsZero() {
		return r.RootThreadTitle
	}

	return fmt.Sprintf("reply to: %s", r.RootThreadTitle)
}

func (r *Reply) GetKind() datagraph.Kind { return datagraph.KindReply }

func (r *Reply) GetSlug() string {
	if xid.ID(r.RootPostID).IsZero() {
		return r.ID.String()
	}
	return r.RootThreadMark
}
func (r *Reply) GetDesc() string { return r.Content.Short() }

func (p Reply) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, p.Content.Short())
}

func replyTo(m *ent.Post) opt.Optional[post.ID] {
	if m.Edges.ReplyTo != nil {
		return opt.New(post.ID(m.Edges.ReplyTo.ID))
	}

	return opt.NewEmpty[post.ID]()
}

func (r *Reply) GetCreated() time.Time { return r.CreatedAt }
func (r *Reply) GetUpdated() time.Time { return r.UpdatedAt }

func Map(m *ent.Post) (*Reply, error) {
	authorEdge, err := m.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.MapRef(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	content, err := datagraph.NewRichText(m.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	replyTo := replyTo(m)

	rootAuthor, err := profile.MapRef(m.Edges.Root.Edges.Author)
	if err != nil {
		return nil, err
	}

	var rootPostID post.ID
	if m.RootPostID != nil {
		rootPostID = post.ID(*m.RootPostID)
	}

	return &Reply{
		Post: post.Post{
			ID: post.ID(m.ID),

			Content: content,
			Author:  *pro,
			Assets:  dt.Map(m.Edges.Assets, asset.Map),
			Meta:    m.Metadata,

			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: opt.NewPtr(m.DeletedAt),
		},
		ReplyTo:         replyTo,
		RootAuthor:      *rootAuthor,
		RootPostID:      rootPostID,
		RootThreadMark:  m.Edges.Root.Slug,
		RootThreadTitle: m.Edges.Root.Title,
	}, nil
}

func Mapper(
	am account.Lookup,
	ls post.PostLikesMap,
	rl reaction.Lookup,
) func(m *ent.Post) (*Reply, error) {
	return func(m *ent.Post) (*Reply, error) {
		authorEdge := am[m.AccountPosts]
		pro, err := profile.MapRef(authorEdge)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		content, err := datagraph.NewRichText(m.Body)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		replyTo := replyTo(m)

		reacts := rl[xid.ID(m.ID)]

		reply := &Reply{
			Post: post.Post{
				ID: post.ID(m.ID),

				Content:     content,
				Author:      *pro,
				Likes:       ls.Status(m.ID),
				Collections: collection_item_status.Status{
					// NOTE: Members cannot yet add replies to collections.
				},
				Reacts: reacts,
				Assets: dt.Map(m.Edges.Assets, asset.Map),
				Meta:   m.Metadata,

				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
				DeletedAt: opt.NewPtr(m.DeletedAt),
			},
			ReplyTo: replyTo,
		}

		if m.Edges.Root != nil {
			var rootPostID post.ID
			if m.RootPostID != nil {
				rootPostID = post.ID(*m.RootPostID)
			}
			rootThreadMark := opt.NewPtr(m.Edges.Root).OrZero().Slug
			rootThreadTitle := opt.NewPtr(m.Edges.Root).OrZero().Title

			slug := fmt.Sprintf("%s#%s", rootThreadMark, m.ID)

			rootAuthor, err := profile.MapRef(m.Edges.Root.Edges.Author)
			if err != nil {
				return nil, fault.Wrap(err)
			}

			reply.RootPostID = rootPostID
			reply.RootThreadMark = rootThreadMark
			reply.RootThreadTitle = rootThreadTitle
			reply.RootAuthor = *rootAuthor
			reply.Slug = slug
		}

		return reply, nil
	}
}
