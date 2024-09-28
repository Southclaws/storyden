package post

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/like"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

// ID wraps the underlying xid type for all kinds of Storyden Post data type.
type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }

type Post struct {
	ID   ID
	Root ID // Identical to ID if this is the root.

	Title   string
	Slug    string
	Content datagraph.Content
	Author  profile.Public
	Likes   like.Status
	Reacts  []*reaction.React
	Assets  []*asset.Asset
	WebLink opt.Optional[link_ref.LinkRef]
	Meta    map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

func (p *Post) GetID() xid.ID                 { return xid.ID(p.ID) }
func (p *Post) GetKind() datagraph.Kind       { return datagraph.KindPost }
func (p *Post) GetName() string               { return p.Title }
func (p *Post) GetSlug() string               { return p.Slug }
func (p *Post) GetContent() datagraph.Content { return p.Content }
func (p *Post) GetDesc() string               { return p.Content.Short() }
func (p *Post) GetProps() map[string]any      { return p.Meta }
func (p *Post) GetAssets() []*asset.Asset     { return p.Assets }

func Map(in *ent.Post) (*Post, error) {
	rootID, title, slug := func() (ID, string, string) {
		if in.First {
			return ID(in.ID), in.Title, in.Slug
		}

		return ID(in.Edges.Root.ID), in.Edges.Root.Title, in.Edges.Root.Slug
	}()

	authorEdge, err := in.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.ProfileFromModel(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	content, err := datagraph.NewRichText(in.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	// This edge is optional anyway, so if not loaded, nothing bad happens.
	link := opt.Map(opt.NewPtr(in.Edges.Link), func(in ent.Link) link_ref.LinkRef {
		return *link_ref.Map(&in)
	})

	// These edges are arrays so if not loaded, nothing bad happens.
	reacts, err := reaction.MapList(in.Edges.Reacts)
	if err != nil {
		return nil, err
	}

	assets := dt.Map(in.Edges.Assets, asset.FromModel)

	return &Post{
		ID:   ID(in.ID),
		Root: rootID,

		Title:   title,
		Slug:    slug,
		Content: content,
		Author:  *pro,
		Reacts:  reacts,
		Assets:  assets,
		WebLink: link,
		Meta:    in.Metadata,

		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		DeletedAt: opt.NewPtr(in.DeletedAt),
	}, nil
}
