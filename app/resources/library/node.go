package library

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
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
	ID        NodeID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name           string
	Slug           string
	Assets         []*asset.Asset
	Links          datagraph.Links
	Content        opt.Optional[content.Rich]
	Description    opt.Optional[string]
	Owner          profile.Public
	Parent         opt.Optional[Node]
	Visibility     visibility.Visibility
	RelevanceScore opt.Optional[float64]
	Metadata       map[string]any

	Nodes []*Node
}

func (*Node) GetResourceName() string { return "node" }

func (c *Node) GetID() xid.ID           { return xid.ID(c.ID) }
func (c *Node) GetKind() datagraph.Kind { return datagraph.KindNode }
func (c *Node) GetName() string         { return c.Name }
func (c *Node) GetSlug() string         { return c.Slug }
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
func (c *Node) GetText() string { return c.Content.OrZero().HTML() }
func (c *Node) GetProps() any   { return nil }
