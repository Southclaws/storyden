package collection

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/internal/ent"
)

type CollectionID xid.ID

func (i CollectionID) String() string { return xid.ID(i).String() }

type Collection struct {
	ID          CollectionID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       datagraph.Profile
	Name        string
	Description string
	Items       []*CollectionItem
}

func (*Collection) GetResourceName() string { return "collection" }

type CollectionItem struct {
	Author datagraph.Profile
	Item   datagraph.Indexable
}

func FromModel(c *ent.Collection) (*Collection, error) {
	accEdge, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := datagraph.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts, err := dt.MapErr(c.Edges.Posts, MapCollectionPost)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	nodes, err := dt.MapErr(c.Edges.Nodes, MapCollectionNode)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	items := append(posts, nodes...)

	return &Collection{
		ID:          CollectionID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Owner:       *pro,
		Name:        c.Name,
		Description: c.Description,
		Items:       items,
	}, nil
}

func MapCollectionPost(p *ent.Post) (*CollectionItem, error) {
	accEdge, err := p.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := datagraph.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	item, err := reply.FromModel(p)
	if err != nil {
		return nil, err
	}

	return &CollectionItem{
		Author: *pro,
		Item:   item,
	}, nil
}

func MapCollectionNode(p *ent.Node) (*CollectionItem, error) {
	accEdge, err := p.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := datagraph.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	item, err := datagraph.NodeFromModel(p)
	if err != nil {
		return nil, err
	}

	return &CollectionItem{
		Author: *pro,
		Item:   item,
	}, nil
}
