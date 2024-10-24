package library

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type (
	NodeID   xid.ID
	NodeSlug string
)

func NodeIDFromString(id string) (NodeID, error) {
	parsed, err := xid.FromString(id)
	if err != nil {
		return NodeID(xid.NilID()), err
	}

	return NodeID(parsed), nil
}

func (i NodeID) String() string { return xid.ID(i).String() }

type Node struct {
	Mark      Mark
	CreatedAt time.Time
	UpdatedAt time.Time

	Name           string
	Assets         []*asset.Asset
	WebLink        opt.Optional[link_ref.LinkRef]
	Content        opt.Optional[datagraph.Content]
	Description    opt.Optional[string]
	PrimaryImage   opt.Optional[asset.Asset]
	Owner          profile.Public
	Parent         opt.Optional[Node]
	Tags           tag_ref.Tags
	Visibility     visibility.Visibility
	RelevanceScore opt.Optional[float64]
	Metadata       map[string]any

	Nodes []*Node
}

func (*Node) GetResourceName() string { return "node" }

func (c *Node) GetID() xid.ID           { return c.Mark.ID() }
func (c *Node) GetKind() datagraph.Kind { return datagraph.KindNode }
func (c *Node) GetName() string         { return c.Name }
func (c *Node) GetSlug() string         { return c.Mark.Slug() }
func (c *Node) GetDesc() string {
	if d, ok := c.Description.Get(); ok && d != "" {
		return d
	}

	cd, ok := c.Content.Get()
	if ok && cd.Short() != "" {
		return cd.Short()
	}

	return ""
}
func (c *Node) GetContent() datagraph.Content { return c.Content.OrZero() }
func (c *Node) GetProps() map[string]any      { return c.Metadata }
func (c *Node) GetAssets() []*asset.Asset     { return c.Assets }
func (c *Node) GetCreated() time.Time         { return c.CreatedAt }
func (c *Node) GetUpdated() time.Time         { return c.UpdatedAt }
