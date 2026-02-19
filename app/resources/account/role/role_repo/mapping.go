package role_repo

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
)

type Hydrator interface {
	Hydrate(accountID xid.ID) (held.Roles, error)
}

type staticHydrator struct {
	lookup map[xid.ID]held.Roles
}

func BuildMultiHydrator(lookup map[xid.ID]held.Roles) Hydrator {
	if lookup == nil {
		lookup = map[xid.ID]held.Roles{}
	}

	return &staticHydrator{
		lookup: lookup,
	}
}

func BuildSingleHydrator(accountID xid.ID, roles held.Roles) Hydrator {
	return BuildMultiHydrator(map[xid.ID]held.Roles{
		accountID: roles,
	})
}

func (h *staticHydrator) Hydrate(accountID xid.ID) (held.Roles, error) {
	roles, ok := h.lookup[accountID]
	if !ok {
		return held.Roles{}, nil
	}

	return roles, nil
}
