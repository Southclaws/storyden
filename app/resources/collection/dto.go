package collection

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
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
	Items       []*thread.Thread
}

func (*Collection) GetResourceName() string { return "collection" }

type Item struct {
	ID     post.ID
	Slug   string
	Author string
	Title  string
	Short  string
}

func FromModel(c *ent.Collection) (*Collection, error) {
	acc, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts, err := c.Edges.PostsOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	items := dt.Map(posts, thread.FromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Collection{
		ID:          CollectionID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Owner:       *account.FromModel(*acc),
		Name:        c.Name,
		Description: c.Description,
		Items:       items,
	}, nil
}
