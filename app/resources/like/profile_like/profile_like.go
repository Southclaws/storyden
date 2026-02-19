// Package profile_like defines a Like within the context of a profile. Because
// the call site will already have the profile, it doesn't need the owner field.
package profile_like

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

// Like on the profile side does not contain the owner, just the liked item.
type Like struct {
	ID      xid.ID
	Created time.Time
	Item    datagraph.Item
}

func Map(in *ent.LikePost, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Like, error) {
	postEdge, err := in.Edges.PostOrErr()
	if err != nil {
		return nil, err
	}

	item, err := post.Map(postEdge, roleHydratorFn)
	if err != nil {
		return nil, err
	}

	return &Like{
		ID:      in.ID,
		Created: in.CreatedAt,
		Item:    item,
	}, nil
}
