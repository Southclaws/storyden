package datagraph

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
)

type (
	ClusterID   xid.ID
	ClusterSlug string
)

func (i ClusterID) String() string { return xid.ID(i).String() }

type Cluster struct {
	ID        ClusterID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Slug        string
	Assets      []*asset.Asset
	Links       Links
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

func (c *Cluster) GetID() xid.ID   { return xid.ID(c.ID) }
func (c *Cluster) GetKind() Kind   { return KindCluster }
func (c *Cluster) GetName() string { return c.Name }
func (c *Cluster) GetText() string { return c.Description }
func (c *Cluster) GetProps() any   { return nil }
