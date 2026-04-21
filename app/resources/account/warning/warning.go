package warning

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type Warning struct {
	ID       ID
	IssuedAt time.Time
	IssuedBy opt.Optional[account.Account]
	Reason   string
}

type Warnings []*Warning

func Map(in *ent.Warning) (*Warning, error) {
	issuedBy := opt.NewEmpty[account.Account]()
	if in.Edges.Author != nil {
		author, err := account.MapRef(in.Edges.Author)
		if err != nil {
			return nil, err
		}
		issuedBy = opt.New(*author)
	}

	return &Warning{
		ID:       in.ID,
		IssuedAt: in.CreatedAt,
		IssuedBy: issuedBy,
		Reason:   in.Reason,
	}, nil
}
