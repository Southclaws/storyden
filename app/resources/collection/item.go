package collection

import (
	"sort"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type CollectionItem struct {
	Added          time.Time
	MembershipType MembershipType
	Author         profile.Ref
	Item           datagraph.Item
	RelevanceScore opt.Optional[float64]
}

type CollectionItemStatus struct {
	Collection Collection
	Item       opt.Optional[CollectionItem]
}

type CollectionItems []*CollectionItem

func (a CollectionItems) Len() int           { return len(a) }
func (a CollectionItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CollectionItems) Less(i, j int) bool { return a[i].Added.After(a[j].Added) }

func MapWithItems(c *ent.Collection) (*CollectionWithItems, error) {
	col, err := Map(nil)(c)
	if err != nil {
		return nil, err
	}

	posts, err := dt.MapErr(c.Edges.CollectionPosts, MapPost)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	nodes, err := dt.MapErr(c.Edges.CollectionNodes, MapNode)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	items := CollectionItems(append(posts, nodes...))

	sort.Sort(items)

	colWithItems := &CollectionWithItems{
		Collection: *col,
		Items:      items,
	}

	return colWithItems, nil
}

func MapPost(n *ent.CollectionPost) (*CollectionItem, error) {
	p := n.Edges.Post

	accEdge, err := p.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mt, err := NewMembershipType(n.MembershipType)
	if err != nil {
		return nil, err
	}

	pro, err := profile.MapRef(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	item, err := post.Map(p)
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

func MapNode(n *ent.CollectionNode) (*CollectionItem, error) {
	p := n.Edges.Node

	accEdge, err := p.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mt, err := NewMembershipType(n.MembershipType)
	if err != nil {
		return nil, err
	}

	pro, err := profile.MapRef(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	item, err := library.MapNode(true, nil)(p)
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
