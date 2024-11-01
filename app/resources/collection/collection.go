package collection

import (
	"sort"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type CollectionID xid.ID

func (i CollectionID) String() string { return xid.ID(i).String() }

type Collection struct {
	ID        CollectionID
	CreatedAt time.Time
	UpdatedAt time.Time
	Owner     profile.Public

	Name        string
	Description opt.Optional[string]

	ItemCount      uint
	HasQueriedItem bool
}

type CollectionWithItems struct {
	Collection
	Items CollectionItems
}

func (*Collection) GetResourceName() string { return "collection" }

type CollectionItem struct {
	Added          time.Time
	MembershipType MembershipType
	Author         profile.Public
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

func Map(queriedItems []xid.ID) func(c *ent.Collection) (*Collection, error) {
	return func(c *ent.Collection) (*Collection, error) {
		accEdge, err := c.Edges.OwnerOrErr()
		if err != nil {
			return nil, fault.Wrap(err)
		}

		postsEdge := c.Edges.CollectionPosts

		nodesEdge := c.Edges.CollectionNodes

		pro, err := profile.ProfileFromModel(accEdge)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		contains := make(map[xid.ID]struct{})
		for _, cp := range postsEdge {
			contains[cp.PostID] = struct{}{}
		}
		for _, cn := range nodesEdge {
			contains[cn.NodeID] = struct{}{}
		}

		var hasQueriedItem bool
		for _, qi := range queriedItems {
			if _, ok := contains[qi]; ok {
				hasQueriedItem = true
				break
			}
		}

		return &Collection{
			ID:             CollectionID(c.ID),
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
			Owner:          *pro,
			Name:           c.Name,
			Description:    opt.NewPtr(c.Description),
			ItemCount:      uint(len(postsEdge) + len(nodesEdge)),
			HasQueriedItem: hasQueriedItem,
		}, nil
	}
}

func MapList(queriedItems []xid.ID, c []*ent.Collection) ([]*Collection, error) {
	return dt.MapErr(c, Map(queriedItems))
}

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

	pro, err := profile.ProfileFromModel(accEdge)
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

	pro, err := profile.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	item, err := library.NodeFromModel(p)
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
