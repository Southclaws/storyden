package invitation

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

type Invitation struct {
	ID        xid.ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
	Message   opt.Optional[string]
	Creator   account.Account
}

func Map(in *ent.Invitation) (*Invitation, error) {
	creatorEdge, err := in.Edges.CreatorOrErr()
	if err != nil {
		return nil, err
	}

	acc, err := account.MapRef(creatorEdge)
	if err != nil {
		return nil, err
	}

	return &Invitation{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		DeletedAt: opt.NewPtr(in.DeletedAt),
		Message:   opt.NewPtr(in.Message),
		Creator:   *acc,
	}, nil
}
