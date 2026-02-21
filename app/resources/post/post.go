package post

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/like"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

// ID wraps the underlying xid type for all kinds of Storyden Post data type.
type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }

func (u ID) MarshalJSON() ([]byte, error) {
	return xid.ID(u).MarshalJSON()
}

func (u *ID) UnmarshalJSON(data []byte) error {
	var id xid.ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	*u = ID(id)
	return nil
}

type Post struct {
	ID   ID
	Root ID // Identical to ID if this is the root.

	Title       string
	Slug        string
	Content     datagraph.Content
	Author      profile.Ref
	Likes       like.Status
	Collections collection_item_status.Status
	Reacts      []*reaction.React
	Assets      []*asset.Asset
	WebLink     opt.Optional[link_ref.LinkRef]
	Meta        map[string]any
	Visibility  visibility.Visibility

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
	IndexedAt opt.Optional[time.Time]
}

type PostRef struct {
	ID   ID
	Root ID
}

func (p *PostRef) IsThread() bool {
	return p.ID == p.Root
}

func (p *Post) GetID() xid.ID                 { return xid.ID(p.ID) }
func (p *Post) GetKind() datagraph.Kind       { return datagraph.KindPost }
func (p *Post) GetName() string               { return p.Title }
func (p *Post) GetSlug() string               { return p.Slug }
func (p *Post) GetContent() datagraph.Content { return p.Content }
func (p *Post) GetDesc() string               { return p.Content.Short() }
func (p *Post) GetProps() map[string]any      { return p.Meta }
func (p *Post) GetAssets() []*asset.Asset     { return p.Assets }
func (p *Post) GetCreated() time.Time         { return p.CreatedAt }
func (p *Post) GetUpdated() time.Time         { return p.UpdatedAt }

func Map(in *ent.Post) (*Post, error) {
	rootID, title, slug := func() (ID, string, string) {
		if in.RootPostID == nil {
			return ID(in.ID), in.Title, in.Slug
		}

		return ID(in.Edges.Root.ID), in.Edges.Root.Title, in.Edges.Root.Slug
	}()

	authorEdge, err := in.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.MapRef(authorEdge)
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

	assets := dt.Map(in.Edges.Assets, asset.Map)

	vis, err := visibility.NewVisibility(string(in.Visibility))
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Post{
		ID:   ID(in.ID),
		Root: rootID,

		Title:      title,
		Slug:       slug,
		Content:    content,
		Author:     *pro,
		Reacts:     reacts,
		Assets:     assets,
		WebLink:    link,
		Meta:       in.Metadata,
		Visibility: vis,

		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		DeletedAt: opt.NewPtr(in.DeletedAt),
		IndexedAt: opt.NewPtr(in.IndexedAt),
	}, nil
}

func MapRef(in *ent.Post) *PostRef {
	root := func() ID {
		if in.RootPostID == nil {
			return ID(in.ID)
		}
		return ID(*in.RootPostID)
	}()

	return &PostRef{
		ID:   ID(in.ID),
		Root: root,
	}
}
