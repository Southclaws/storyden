package role

import (
	"sort"
	"time"

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
	SortKey     float64
	Metadata    map[string]any
	CreatedAt   time.Time
}

type Roles []*Role

func (a Roles) Len() int           { return len(a) }
func (a Roles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Roles) Less(i, j int) bool { return a[i].SortKey < a[j].SortKey }

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

func Map(r *ent.Role) (*Role, error) {
	perms, err := rbac.NewPermissions(r.Permissions)
	if err != nil {
		return nil, err
	}

	return &Role{
		ID:          RoleID(r.ID),
		Name:        r.Name,
		Colour:      r.Colour,
		Permissions: *perms,
		SortKey:     defaultSortKey(RoleID(r.ID), r.SortKey),
		Metadata:    r.Metadata,
		CreatedAt:   r.CreatedAt,
	}, nil
}

func MapList(in []*ent.Role) (Roles, error) {
	mapped, err := dt.MapErr(in, Map)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	sort.Sort(Roles(mapped))

	return mapped, nil
}

func defaultSortKey(id RoleID, current float64) float64 {
	switch id {
	case DefaultRoleGuestID:
		return DefaultRoleGuest.SortKey
	case DefaultRoleMemberID:
		return DefaultRoleMember.SortKey
	case DefaultRoleAdminID:
		return DefaultRoleAdmin.SortKey
	default:
		return current
	}
}
