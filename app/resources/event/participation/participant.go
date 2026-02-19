package participation

import (
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type EventParticipant struct {
	Role    Role
	Status  Status
	Account profile.Ref
}

type EventParticipants []*EventParticipant

func (p EventParticipants) IsHost(id account.AccountID) bool {
	_, isHost := lo.Find(p, func(p *EventParticipant) bool {
		return p.Account.ID == id && p.Role == RoleHost
	})
	return isHost
}

func Map(in *ent.EventParticipant, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*EventParticipant, error) {
	profileMapper := profile.RefMapper(roleHydratorFn)

	accountEdge, err := in.Edges.AccountOrErr()
	if err != nil {
		return nil, err
	}

	role, err := NewRole(in.Role)
	if err != nil {
		return nil, err
	}

	status, err := NewStatus(in.Status)
	if err != nil {
		return nil, err
	}

	acc, err := profileMapper(accountEdge)
	if err != nil {
		return nil, err
	}

	return &EventParticipant{
		Role:    role,
		Status:  status,
		Account: *acc,
	}, nil
}
