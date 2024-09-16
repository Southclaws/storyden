package role

import (
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
)

type RoleID xid.ID

func (u RoleID) String() string { return xid.ID(u).String() }

type Role struct {
	ID          RoleID
	Name        string
	Colour      string
	Permissions rbac.Permissions
}

type Roles []*Role

func (r Roles) Permissions() rbac.Permissions {
	set := map[rbac.Permission]bool{}

	for _, role := range r {
		for _, perm := range role.Permissions.List() {
			set[perm] = true
		}
	}

	flat := lo.Keys(set)

	return rbac.NewList(flat...)
}

func Map(in []*ent.Role, admin bool) (Roles, error) {
	mapped, err := dt.MapErr(in, func(r *ent.Role) (*Role, error) {
		perms, err := rbac.NewPermissions(r.Permissions)
		if err != nil {
			return nil, err
		}

		return &Role{
			ID:          RoleID(r.ID),
			Name:        r.Name,
			Colour:      r.Colour,
			Permissions: *perms,
		}, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mapped = append(mapped, &DefaultRoleEveryone)
	if admin {
		mapped = append(mapped, &DefaultRoleAdmin)
	}

	return mapped, nil
}
