package reaction

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type ReactID xid.ID

type React struct {
	ID     ReactID
	Emoji  string
	Author profile.Ref
	target xid.ID
}

func (r *React) Target() xid.ID { return r.target }

type Reacts []*React

func (r Reacts) Map() Lookup {
	return lo.GroupBy(r, func(r *React) xid.ID { return r.target })
}

type Lookup map[xid.ID]Reacts

func Map(in *ent.React, roleLookup func(accID xid.ID) (held.Roles, error)) (*React, error) {
	accountEdge, err := in.Edges.AccountOrErr()
	if err != nil {
		return nil, err
	}

	profileMapper := profile.RefMapper(roleLookup)
	acc, err := profileMapper(accountEdge)
	if err != nil {
		return nil, err
	}

	return &React{
		ID:     ReactID(in.ID),
		Emoji:  in.Emoji,
		Author: *acc,
		target: xid.ID(in.PostID),
	}, nil
}

func MapList(in []*ent.React, roleLookup func(accID xid.ID) (held.Roles, error)) ([]*React, error) {
	return dt.MapErr(in, func(in *ent.React) (*React, error) {
		return Map(in, roleLookup)
	})
}

func Mapper(am account.Lookup, roleLookup func(accID xid.ID) (held.Roles, error)) func(in *ent.React) (*React, error) {
	profileMapper := profile.RefMapper(roleLookup)
	return func(in *ent.React) (*React, error) {
		acc := am[xid.ID(in.AccountID)]

		pro, err := profileMapper(acc)
		if err != nil {
			return nil, err
		}

		return &React{
			ID:     ReactID(in.ID),
			Emoji:  in.Emoji,
			Author: *pro,
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
