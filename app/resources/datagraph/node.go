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

	Name        string
	Slug        string
	Assets      []*asset.Asset
	Links       Links
	Description string
	Content     opt.Optional[string]
	Owner       profile.Profile
	Parent      opt.Optional[Node]
	Visibility  post.Visibility
	Properties  any

	Nodes []*Node
}

func (*Node) GetResourceName() string { return "node" }

func (c *Node) GetID() xid.ID   { return xid.ID(c.ID) }
func (c *Node) GetKind() Kind   { return KindNode }
func (c *Node) GetName() string { return c.Name }
func (c *Node) GetSlug() string { return c.Slug }
func (c *Node) GetDesc() string { return c.Description }
func (c *Node) GetText() string { return c.Content.String() }
func (c *Node) GetProps() any   { return nil }
