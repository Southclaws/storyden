package datagraph

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/profile"
)

type (
	ItemID   xid.ID
	ItemSlug string
)

func (i ItemID) String() string { return xid.ID(i).String() }

type Item struct {
	ID        ItemID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Slug        string
	Assets      []*asset.Asset
	Links       Links
	Description string
	Content     opt.Optional[string]
	Owner       profile.Profile
	In          []*Cluster
	Properties  any
}

func (*Item) GetResourceName() string { return "Item" }

func (c *Item) GetID() xid.ID   { return xid.ID(c.ID) }
func (c *Item) GetKind() Kind   { return KindItem }
func (c *Item) GetName() string { return c.Name }
func (c *Item) GetSlug() string { return c.Slug }
func (c *Item) GetDesc() string { return c.Description }
func (c *Item) GetText() string { return c.Content.String() }
func (c *Item) GetProps() any   { return nil }
