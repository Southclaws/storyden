package role_repo

import (
	"context"
	"encoding/json"
	"log/slog"
	"sort"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

const (
	roleListKey  = "role:list"
	roleCacheTTL = time.Hour * 24 * 365
)

type cachedRole struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Colour      string         `json:"colour"`
	Permissions []string       `json:"permissions"`
	SortKey     float64        `json:"sort_key"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

func (r cachedRole) toDomain() (*role.Role, error) {
	id, err := xid.FromString(r.ID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	perms, err := rbac.NewPermissions(r.Permissions)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &role.Role{
		ID:          role.RoleID(id),
		Name:        r.Name,
		Colour:      r.Colour,
		Permissions: *perms,
		SortKey:     r.SortKey,
		Metadata:    r.Metadata,
	}, nil
}

func fromDomainRole(r *role.Role) cachedRole {
	return cachedRole{
		ID:          r.ID.String(),
		Name:        r.Name,
		Colour:      r.Colour,
		Permissions: dt.Map(r.Permissions.List(), func(p rbac.Permission) string { return p.String() }),
		SortKey:     r.SortKey,
		Metadata:    r.Metadata,
	}
}

func (r *Repository) getRoleListCache(ctx context.Context) (role.Roles, bool) {
	raw, err := r.store.Get(ctx, roleListKey)
	if err != nil {
		return nil, false
	}

	var payload []cachedRole
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, false
	}

	roles := make(role.Roles, 0, len(payload))
	for _, in := range payload {
		parsed, err := in.toDomain()
		if err != nil {
			return nil, false
		}

		roles = append(roles, parsed)
	}

	sort.Sort(roles)

	return roles, true
}

func (r *Repository) setRoleListCache(ctx context.Context, roles role.Roles) {
	if err := r.setRoleListCacheStrict(ctx, roles); err != nil {
		slog.Error("failed to store role list cache payload",
			slog.String("error", err.Error()),
		)
	}
}

func (r *Repository) setRoleListCacheStrict(ctx context.Context, roles role.Roles) error {
	payload, err := json.Marshal(dt.Map(roles, fromDomainRole))
	if err != nil {
		return fault.Wrap(err)
	}

	if err := r.store.Set(ctx, roleListKey, string(payload), roleCacheTTL); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func (r *Repository) refreshRoleListCache(ctx context.Context) {
	if err := r.refreshRoleListCacheStrict(ctx); err != nil {
		slog.Error("failed to refresh role list cache payload",
			slog.String("error", err.Error()),
		)
	}
}

func (r *Repository) refreshRoleListCacheStrict(ctx context.Context) error {
	roles, err := r.listFromDB(ctx)
	if err != nil {
		return fault.Wrap(err)
	}

	if err := r.setRoleListCacheStrict(ctx, roles); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func (r *Repository) getRoleCache(ctx context.Context, id role.RoleID) (*role.Role, bool) {
	roles, ok := r.getRoleListCache(ctx)
	if !ok {
		return nil, false
	}

	for _, rl := range roles {
		if rl.ID == id {
			return rl, true
		}
	}

	return nil, false
}

func (r *Repository) setRoleCache(ctx context.Context, in *role.Role) {
	if err := r.setRoleCacheStrict(ctx, in); err != nil {
		slog.Error("failed to update role list cache payload",
			slog.String("role_id", in.ID.String()),
			slog.String("error", err.Error()),
		)
	}
}

func (r *Repository) setRoleCacheStrict(ctx context.Context, in *role.Role) error {
	roles, ok := r.getRoleListCache(ctx)
	if !ok {
		loaded, err := r.listFromDB(ctx)
		if err != nil {
			return fault.Wrap(err)
		}
		roles = loaded
	}

	updated := false
	for i := range roles {
		if roles[i].ID == in.ID {
			copy := *in
			roles[i] = &copy
			updated = true
			break
		}
	}

	if !updated {
		copy := *in
		roles = append(roles, &copy)
	}

	sort.Sort(roles)

	if err := r.setRoleListCacheStrict(ctx, roles); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func (r *Repository) removeRoleCacheStrict(ctx context.Context, id role.RoleID) error {
	roles, ok := r.getRoleListCache(ctx)
	if !ok {
		loaded, err := r.listFromDB(ctx)
		if err != nil {
			return fault.Wrap(err)
		}
		roles = loaded
	}

	filtered := make(role.Roles, 0, len(roles))
	for _, rl := range roles {
		if rl.ID == id {
			continue
		}

		filtered = append(filtered, rl)
	}

	if err := r.setRoleListCacheStrict(ctx, filtered); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func (r *Repository) listFromDB(ctx context.Context) (role.Roles, error) {
	rows, err := r.db.Role.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	list, err := role.MapList(rows)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	byID := map[role.RoleID]struct{}{}
	for _, rl := range list {
		byID[rl.ID] = struct{}{}
	}

	if _, ok := byID[role.DefaultRoleGuestID]; !ok {
		copy := role.DefaultRoleGuest
		list = append(list, &copy)
	}
	if _, ok := byID[role.DefaultRoleMemberID]; !ok {
		copy := role.DefaultRoleMember
		list = append(list, &copy)
	}
	if _, ok := byID[role.DefaultRoleAdminID]; !ok {
		copy := role.DefaultRoleAdmin
		list = append(list, &copy)
	}

	sort.Sort(list)

	return list, nil
}
