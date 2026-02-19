package role_badge

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
)

type Badge interface {
	UpdateBadge(ctx context.Context, accountID xid.ID, roleID role.RoleID, badge bool) error
}
