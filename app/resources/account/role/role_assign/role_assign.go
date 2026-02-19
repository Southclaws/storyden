package role_assign

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
)

type Assign interface {
	UpdateRoles(ctx context.Context, accountID xid.ID, roles ...Mutation) error
}

type Mutation struct {
	id       role.RoleID
	isDelete bool
}

func Add(id role.RoleID) Mutation {
	return Mutation{id: id}
}

func Remove(id role.RoleID) Mutation {
	return Mutation{id: id, isDelete: true}
}

func (m Mutation) ID() role.RoleID {
	return m.id
}

func (m Mutation) IsDelete() bool {
	return m.isDelete
}
