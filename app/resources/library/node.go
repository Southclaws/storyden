package library

import (
	"time"

	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
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
	IndexedAt opt.Optional[time.Time]

	Name            string
	Assets          []*asset.Asset
	WebLink         opt.Optional[link_ref.LinkRef]
	Content         opt.Optional[datagraph.Content]
	Description     opt.Optional[string]
	PrimaryImage    opt.Optional[asset.Asset]
	Owner           profile.Ref
	Parent          opt.Optional[Node]
	Properties      opt.Optional[PropertyTable]
	ChildProperties opt.Optional[PropertySchema]
	HideChildTree   bool
	Tags            tag_ref.Tags
	Collections     collection_item_status.Status // NOTE: Not done yet
	Visibility      visibility.Visibility
	SortKey         lexorank.Key
	RelevanceScore  opt.Optional[float64]
	Metadata        map[string]any

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
func (c *Node) GetAuthor() xid.ID             { return xid.ID(c.Owner.ID) }
func (c *Node) GetTags() []string {
	tags := make([]string, len(c.Tags))
	for i, tag := range c.Tags {
		tags[i] = tag.Name.String()
	}
	return tags
}
