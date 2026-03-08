package robot_ref

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID xid.ID

func (id ID) String() string {
	return xid.ID(id).String()
}

func NewID(s string) (ID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return ID{}, err
	}
	return ID(id), nil
}

type Robot struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Playbook    string
	Tools       []string
	Metadata    map[string]any

	AuthorID account.AccountID
}

type Robots []*Robot

func Map(in *ent.Robot) *Robot {
	return &Robot{
		ID:        ID(in.ID),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,

		Name:        in.Name,
		Description: in.Description,
		Playbook:    in.Playbook,
		Tools:       in.Tools,
		Metadata:    in.Metadata,

		AuthorID: account.AccountID(in.AuthorID),
	}
}
