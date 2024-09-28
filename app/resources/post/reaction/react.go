package reaction

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ReactID xid.ID

type React struct {
	ID     ReactID
	Emoji  string
	Author account.Account
}

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

func IsValidEmoji(e string) (string, bool) {
	if len(e) == 0 {
		return "", false
	}

	if e[1] == ':' {
		return "", false
	}

	return e, gomoji.ContainsEmoji(e)
}
