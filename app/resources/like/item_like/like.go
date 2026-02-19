// Package item_like defines a Like on an item such as a post. Because the call
// site will already have the item itself, this only needs to contain the owner.
package item_like

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

// Like on the item side does not contain the item itself, just the owner.
type Like struct {
	ID      xid.ID
	Created time.Time
	Owner   profile.Ref
}

func Map(in *ent.LikePost, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Like, error) {
	profileMapper := profile.RefMapper(roleHydratorFn)

	owner, err := profileMapper(in.Edges.Account)
	if err != nil {
		return nil, err
	}

	return &Like{
		ID:      in.ID,
		Created: in.CreatedAt,
		Owner:   *owner,
	}, nil
}
