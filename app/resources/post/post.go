package post

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

// ID wraps the underlying xid type for all kinds of Storyden Post data type.
type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }

type Post struct {
	ID ID

	Content content.Rich
	Author  profile.Public
	Reacts  []*react.React
	Assets  []*asset.Asset
	Links   datagraph.Links
	Meta    map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

func (p *Post) GetID() xid.ID             { return xid.ID(p.ID) }
func (p *Post) GetKind() datagraph.Kind   { return datagraph.KindPost }
func (p *Post) GetContent() content.Rich  { return p.Content }
func (p *Post) GetProps() map[string]any  { return p.Meta }
func (p *Post) GetAssets() []*asset.Asset { return p.Assets }

func Map(in *ent.Post) (*Post, error) {
	return nil, nil
}
