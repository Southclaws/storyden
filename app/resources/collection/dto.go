package collection

import (
	"sort"
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
	Items       CollectionItems
}

func (*Collection) GetResourceName() string { return "collection" }

type CollectionItem struct {
	Added          time.Time
	MembershipType MembershipType
	Author         datagraph.Profile
	Item           datagraph.Indexable
}

type CollectionItems []*CollectionItem

func (a CollectionItems) Len() int           { return len(a) }
func (a CollectionItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CollectionItems) Less(i, j int) bool { return a[i].Added.After(a[j].Added) }

func FromModel(c *ent.Collection) (*Collection, error) {
	accEdge, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := datagraph.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts, err := dt.MapErr(c.Edges.CollectionPosts, MapCollectionPost)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	nodes, err := dt.MapErr(c.Edges.CollectionNodes, MapCollectionNode)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	items := CollectionItems(append(posts, nodes...))

	sort.Sort(items)

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

func MapCollectionPost(n *ent.CollectionPost) (*CollectionItem, error) {
	p := n.Edges.Post

	accEdge, err := p.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mt, err := NewMembershipType(n.MembershipType)
	if err != nil {
		return nil, err
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
		Added:          n.CreatedAt,
		MembershipType: mt,
		Author:         *pro,
		Item:           item,
	}, nil
}

func MapCollectionNode(n *ent.CollectionNode) (*CollectionItem, error) {
	p := n.Edges.Node

	accEdge, err := p.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mt, err := NewMembershipType(n.MembershipType)
	if err != nil {
		return nil, err
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
		Added:          n.CreatedAt,
		MembershipType: mt,
		Author:         *pro,
		Item:           item,
	}, nil
}
