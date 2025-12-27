package reply

import (
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

const RepliesPerPage = 50

type Reply struct {
	post.Post

	RootPostID      post.ID
	RootThreadMark  string
	RootThreadTitle string
	RootAuthor      profile.Ref
	Slug            string // The root slug with the post ID as a #fragment
	ReplyTo         opt.Optional[Reply]
}

type ReplyRef struct {
	ID         post.ID
	RootPostID post.ID
}

func (*Reply) GetResourceName() string { return "post" }

func (r *Reply) GetName() string {
	if xid.ID(r.RootPostID).IsZero() {
		return r.RootThreadTitle
	}

	if r.RootThreadTitle == "" {
		return ""
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

func (r *Reply) GetCreated() time.Time { return r.CreatedAt }
func (r *Reply) GetUpdated() time.Time { return r.UpdatedAt }
func (r *Reply) GetAuthor() xid.ID     { return xid.ID(r.Author.ID) }

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

	replyTo, err := func(m *ent.Post) (opt.Optional[Reply], error) {
		if m.Edges.ReplyTo == nil {
			return opt.NewEmpty[Reply](), nil
		}

		r, err := Map(m.Edges.ReplyTo)
		if err != nil {
			return nil, err
		}

		return opt.New(*r), nil
	}(m)
	if err != nil {
		return nil, err
	}

	reply := &Reply{
		Post: post.Post{
			ID: post.ID(m.ID),

			Content:    content,
			Author:     *pro,
			Visibility: visibility.NewVisibilityFromEnt(m.Visibility),
			Assets:     dt.Map(m.Edges.Assets, asset.Map),
			Meta:       m.Metadata,

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

		reply.RootPostID = rootPostID
		reply.RootThreadMark = rootThreadMark
		reply.RootThreadTitle = rootThreadTitle

		reply.Slug = slug

		if m.Edges.Root.Edges.Author != nil {
			p, err := profile.MapRef(m.Edges.Root.Edges.Author)
			if err != nil {
				return nil, err
			}
			reply.RootAuthor = *p
		}
	}

	return reply, nil
}

func Mapper(
	am account.Lookup,
	ls post.PostLikesMap,
	rl reaction.Lookup,
) func(m *ent.Post) (*Reply, error) {
	mapReplyTo := func(m *ent.Post) (opt.Optional[Reply], error) {
		if m.Edges.ReplyTo == nil {
			return opt.NewEmpty[Reply](), nil
		}

		r, err := Mapper(am, ls, rl)(m.Edges.ReplyTo)
		if err != nil {
			return nil, err
		}

		return opt.New(*r), nil
	}

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

		replyTo, err := mapReplyTo(m)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		reacts := rl[xid.ID(m.ID)]

		reply := &Reply{
			Post: post.Post{
				ID: post.ID(m.ID),

				Content:     content,
				Author:      *pro,
				Visibility:  visibility.NewVisibilityFromEnt(m.Visibility),
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

func MapRef(m *ent.Post) *ReplyRef {
	root := func() xid.ID {
		if m.RootPostID == nil {
			return m.ID
		}
		return *m.RootPostID
	}()

	return &ReplyRef{
		ID:         post.ID(m.ID),
		RootPostID: post.ID(root),
	}
}

func ItemRef(r *ent.Post) (datagraph.Item, error) {
	content, err := datagraph.NewRichText(r.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	var rootPostID post.ID
	if r.RootPostID != nil {
		rootPostID = post.ID(*r.RootPostID)
	}

	rootSlug := opt.NewPtrMap(r.Edges.Root, func(p ent.Post) string {
		return p.Slug
	}).Or(r.RootPostID.String())

	return &Reply{
		Post: post.Post{
			ID:         post.ID(r.ID),
			Content:    content,
			Visibility: visibility.NewVisibilityFromEnt(r.Visibility),
			Meta:       r.Metadata,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
			DeletedAt:  opt.NewPtr(r.DeletedAt),
		},
		RootPostID: rootPostID,
		Slug:       fmt.Sprintf("%s#%s", rootSlug, r.ID),
	}, nil
}
