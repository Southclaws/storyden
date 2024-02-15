package datagraph

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
)

type Indexable interface {
	GetID() xid.ID
	GetType() string
	GetName() string
	GetText() string
	GetProps() any
}

type (
	ClusterID   xid.ID
	ClusterSlug string
	ItemID      xid.ID
	ItemSlug    string
)

func (i ClusterID) String() string { return xid.ID(i).String() }
func (i ItemID) String() string    { return xid.ID(i).String() }

type Item struct {
	ID        ItemID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Slug        string
	Assets      []*asset.Asset
	Links       link.Links
	Description string
	Content     opt.Optional[string]
	Owner       profile.Profile
	In          []*Cluster
	Properties  any
}

func (*Item) GetResourceName() string { return "Item" }

type Cluster struct {
	ID        ClusterID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Slug        string
	Assets      []*asset.Asset
	Links       link.Links
	Description string
	Content     opt.Optional[string]
	Owner       profile.Profile
	Parent      opt.Optional[*Cluster]
	Visibility  post.Visibility
	Properties  any

	Items    []*Item
	Clusters []*Cluster
}

func (*Cluster) GetResourceName() string { return "cluster" }
