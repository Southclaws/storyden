package participation

import (
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/samber/lo"
)

type EventParticipant struct {
	Role    Role
	Status  Status
	Account account.Account
}

type EventParticipants []*EventParticipant

func (p EventParticipants) IsHost(id account.AccountID) bool {
	_, isHost := lo.Find(p, func(p *EventParticipant) bool {
		return p.Account.ID == id && p.Role == RoleHost
	})
	return isHost
}

func Map(in *ent.EventParticipant) (*EventParticipant, error) {
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

	acc, err := account.MapAccount(accountEdge)
	if err != nil {
		return nil, err
	}

	return &EventParticipant{
		Role:    role,
		Status:  status,
		Account: *acc,
	}, nil
}
