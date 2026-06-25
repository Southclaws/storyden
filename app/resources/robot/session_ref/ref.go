package session_ref

import (
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

type ID xid.ID

func (id ID) String() string {
	return xid.ID(id).String()
}

type Ref struct {
	ID        ID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Human     account.Account
}

func Map(in *ent.RobotSession) (*Ref, error) {
	acc, err := account.MapRef(in.Edges.User)
	if err != nil {
		return nil, err
	}

	return &Ref{
		ID:        ID(in.ID),
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		Human:     *acc,
	}, nil
}
