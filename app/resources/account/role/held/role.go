package held

import (
	"sort"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
)

// held.Role represents an instance of a role associated with an account. It can
// contain additional properties specific to the relationship to the holder.
type Role struct {
	role.Role

	Assigned time.Time
	Badge    bool
	Default  bool
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

func Map(in *ent.AccountRoles) (*Role, error) {
	r, err := role.Map(in.Edges.Role)
	if err != nil {
		return nil, err
	}

	return &Role{
		Role: *r,

		// CreatedAt is the timestamp of the relationship, not the role itself.
		Assigned: in.CreatedAt,
		Badge:    opt.NewPtr(in.Badge).OrZero(),
	}, nil
}

func MapList(in []*ent.AccountRoles, admin bool) (Roles, error) {
	mapped, err := dt.MapErr(in, Map)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mapped = append(mapped, &Role{Role: role.DefaultRoleEveryone, Default: true})
	if admin {
		mapped = append(mapped, &Role{Role: role.DefaultRoleAdmin, Default: true})
	}

	sort.Sort(Roles(mapped))

	return mapped, nil
}
