package robot_ref

import (
	"time"

	"github.com/Southclaws/dt"
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

type ToolName string

type Robot struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Playbook    string
	Tools       []ToolName
	Metadata    map[string]any

	AuthorID account.AccountID
}

type Robots []*Robot

func Map(in *ent.Robot) *Robot {
	tools := dt.Map(in.Tools, func(tool string) ToolName {
		return ToolName(tool)
	})

	return &Robot{
		ID:        ID(in.ID),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,

		Name:        in.Name,
		Description: in.Description,
		Playbook:    in.Playbook,
		Tools:       tools,
		Metadata:    in.Metadata,

		AuthorID: account.AccountID(in.AuthorID),
	}
}
