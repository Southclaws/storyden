package reaction

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ReactID xid.ID

type React struct {
	ID     ReactID
	Emoji  string
	Author account.Account
	target xid.ID
}

type Reacts []*React

func (r Reacts) Map() Lookup {
	return lo.GroupBy(r, func(r *React) xid.ID { return r.target })
}

type Lookup map[xid.ID]Reacts

func Map(in *ent.React) (*React, error) {
	accountEdge, err := in.Edges.AccountOrErr()
	if err != nil {
		return nil, err
	}

	acc, err := account.MapAccount(accountEdge)
	if err != nil {
		return nil, err
	}

	return &React{
		ID:     ReactID(in.ID),
		Emoji:  in.Emoji,
		Author: *acc,
	}, nil
}

func MapList(in []*ent.React) ([]*React, error) {
	return dt.MapErr(in, Map)
}

func Mapper(am account.Lookup) func(in *ent.React) (*React, error) {
	return func(in *ent.React) (*React, error) {
		acc := am[xid.ID(in.AccountID)]

		return &React{
			ID:     ReactID(in.ID),
			Emoji:  in.Emoji,
			Author: *acc,
			target: xid.ID(in.PostID),
		}, nil
	}
}

func IsValidEmoji(e string) (string, bool) {
	if len(e) == 0 {
		return "", false
	}

	if e[1] == ':' {
		return "", false
	}

	return e, gomoji.ContainsEmoji(e)
}
