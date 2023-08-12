package collection

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type CollectionID xid.ID

func (i CollectionID) String() string { return xid.ID(i).String() }

type Collection struct {
	ID          CollectionID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       profile.Profile
	Name        string
	Description string
	Items       []*Item
}

func (*Collection) GetResourceName() string { return "collection" }

func FromModel(c *ent.Collection) (*Collection, error) {
	accEdge, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.FromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts := opt.NewIf(c.Edges.Posts, func(p []*ent.Post) bool { return p != nil })

	items, err := opt.MapErr[[]*ent.Post, []*Item](posts, func(p []*ent.Post) ([]*Item, error) {
		ps, err := dt.MapErr(p, ItemFromModel)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		return ps, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Collection{
		ID:          CollectionID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Owner:       *pro,
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
	Author    profile.Profile
	Title     string
	Short     string
}

func ItemFromModel(p *ent.Post) (*Item, error) {
	accEdge, err := p.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.FromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Item{
		ID:        post.ID(p.ID),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Slug:      p.Slug,
		Author:    *pro,
		Title:     p.Title,
		Short:     p.Short,
	}, nil
}
