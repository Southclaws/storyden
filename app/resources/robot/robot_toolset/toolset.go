package robot_toolset

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID xid.ID

func (id ID) String() string { return xid.ID(id).String() }

func ParseID(value string) (ID, error) {
	id, err := xid.FromString(value)
	return ID(id), err
}

type Toolset struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Instruction string
	Tools       []string
	Metadata    map[string]any
	Author      account.Account
}

func Map(in *ent.RobotToolset) (*Toolset, error) {
	author, err := account.MapRef(in.Edges.Author)
	if err != nil {
		return nil, err
	}

	return &Toolset{
		ID:          ID(in.ID),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Description: in.Description,
		Instruction: in.Instruction,
		Tools:       append([]string(nil), in.Tools...),
		Metadata:    in.Metadata,
		Author:      *author,
	}, nil
}
