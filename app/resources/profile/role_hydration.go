package profile

import "github.com/Southclaws/storyden/internal/ent"

func RoleHydrationTargets(acc *ent.Account) []*ent.Account {
	targets := []*ent.Account{acc}

	if invitedBy := acc.Edges.InvitedBy; invitedBy != nil {
		if invitedBy.Edges.Creator != nil {
			targets = append(targets, invitedBy.Edges.Creator)
		}
	}

	return targets
}
