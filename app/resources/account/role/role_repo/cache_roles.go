package role_repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

func (h *Repository) roleKey(id role.RoleID) string {
	return roleCachePrefix + id.String()
}

func (h *Repository) deleteRole(ctx context.Context, id role.RoleID) error {
	if err := h.store.Delete(ctx, h.roleKey(id)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (h *Repository) deleteCustomRoleOrdering(ctx context.Context) error {
	if err := h.store.Delete(ctx, roleCustomOrderingCacheKey); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func defaultRole(id role.RoleID) (role.Role, bool) {
	switch id {
	case role.DefaultRoleGuestID:
		return role.DefaultRoleGuest, true
	case role.DefaultRoleMemberID:
		return role.DefaultRoleMember, true
	case role.DefaultRoleAdminID:
		return role.DefaultRoleAdmin, true
	default:
		return role.Role{}, false
	}
}

type cachedRole struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Colour      string         `json:"colour"`
	Permissions []string       `json:"permissions"`
	SortKey     float64        `json:"sort_key"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
	Persisted   bool           `json:"persisted"`
}

func (c cachedRole) toRole() (*role.Role, error) {
	id, err := xid.FromString(c.ID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	perms, err := rbac.NewPermissions(c.Permissions)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &role.Role{
		ID:          role.RoleID(id),
		Name:        c.Name,
		Colour:      c.Colour,
		Permissions: *perms,
		SortKey:     c.SortKey,
		Metadata:    c.Metadata,
		CreatedAt:   c.CreatedAt,
	}, nil
}

func (h *Repository) storeRole(ctx context.Context, rl *role.Role, persisted bool) error {
	data := cachedRole{
		ID:        rl.ID.String(),
		Name:      rl.Name,
		Colour:    rl.Colour,
		SortKey:   rl.SortKey,
		Metadata:  rl.Metadata,
		CreatedAt: rl.CreatedAt,
		Persisted: persisted,
	}

	for _, permission := range rl.Permissions.List() {
		data.Permissions = append(data.Permissions, permission.String())
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := h.store.Set(ctx, h.roleKey(rl.ID), string(raw), cacheTTL); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (h *Repository) storeCustomRoleOrdering(ctx context.Context, ids []role.RoleID) error {
	serialised := make([]string, 0, len(ids))
	for _, id := range ids {
		serialised = append(serialised, id.String())
	}

	raw, err := json.Marshal(serialised)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := h.store.Set(ctx, roleCustomOrderingCacheKey, string(raw), cacheTTL); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (h *Repository) cachedRole(ctx context.Context, id role.RoleID) (cachedRole, bool) {
	raw, err := h.store.Get(ctx, h.roleKey(id))
	if err != nil {
		return cachedRole{}, false
	}

	var out cachedRole
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		_ = h.deleteRole(ctx, id)
		return cachedRole{}, false
	}

	return out, true
}

func (h *Repository) cachedCustomRoleOrdering(ctx context.Context) ([]role.RoleID, bool) {
	raw, err := h.store.Get(ctx, roleCustomOrderingCacheKey)
	if err != nil {
		return nil, false
	}

	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		_ = h.deleteCustomRoleOrdering(ctx)
		return nil, false
	}

	ids := make([]role.RoleID, 0, len(out))
	for _, item := range out {
		id, err := xid.FromString(item)
		if err != nil {
			_ = h.deleteCustomRoleOrdering(ctx)
			return nil, false
		}

		ids = append(ids, role.RoleID(id))
	}

	return ids, true
}
