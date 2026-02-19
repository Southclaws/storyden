package collection

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type Collection struct {
	Mark      Mark
	CreatedAt time.Time
	UpdatedAt time.Time
	IndexedAt opt.Optional[time.Time]

	Name        string
	Owner       profile.Ref
	Description opt.Optional[string]
	Cover       opt.Optional[asset.Asset]

	ItemCount      uint
	HasQueriedItem bool
}

type CollectionWithItems struct {
	Collection
	Items CollectionItems
}

func Map(queriedItems []xid.ID, roleHydratorFn func(accID xid.ID) (held.Roles, error)) func(c *ent.Collection) (*Collection, error) {
	profileMapper := profile.RefMapper(roleHydratorFn)

	return func(c *ent.Collection) (*Collection, error) {
		accEdge, err := c.Edges.OwnerOrErr()
		if err != nil {
			return nil, fault.Wrap(err)
		}

		postsEdge := c.Edges.CollectionPosts

		nodesEdge := c.Edges.CollectionNodes

		pro, err := profileMapper(accEdge)
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
			Mark:           NewMark(c.ID, c.Slug),
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
			IndexedAt:      opt.NewPtr(c.IndexedAt),
			Owner:          *pro,
			Name:           c.Name,
			Description:    opt.NewPtr(c.Description),
			ItemCount:      uint(len(postsEdge) + len(nodesEdge)),
			HasQueriedItem: hasQueriedItem,
		}, nil
	}
}

func MapList(queriedItems []xid.ID, c []*ent.Collection, roleHydratorFn func(accID xid.ID) (held.Roles, error)) ([]*Collection, error) {
	return dt.MapErr(c, Map(queriedItems, roleHydratorFn))
}
