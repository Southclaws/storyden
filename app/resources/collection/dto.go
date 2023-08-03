package collection

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type CollectionID xid.ID

func (i CollectionID) String() string { return xid.ID(i).String() }

type Collection struct {
	ID          CollectionID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       account.Account
	Name        string
	Description string
	Items       []*Item
}

func (*Collection) GetResourceName() string { return "collection" }

func FromModel(c *ent.Collection) (*Collection, error) {
	acc, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts := opt.NewIf(c.Edges.Posts, func(p []*ent.Post) bool { return p != nil })

	items := opt.Map[[]*ent.Post, []*Item](posts, func(p []*ent.Post) []*Item {
		return dt.Map(p, ItemFromModel)
	})

	return &Collection{
		ID:          CollectionID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Owner:       *account.FromModel(*acc),
		Name:        c.Name,
		Description: c.Description,
		Items:       items.Or([]*Item{}),
	}, nil
}

type Item struct {
	ID        post.ID
	CreatedAt time.Time
	UpdatedAt time.Time
	Slug      string
	Author    account.Account
	Title     string
	Short     string
}

func ItemFromModel(p *ent.Post) *Item {
	return &Item{
		ID:        post.ID(p.ID),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Slug:      p.Slug,
		Author:    *account.FromModel(*p.Edges.Author),
		Title:     p.Title,
		Short:     p.Short,
	}
}
