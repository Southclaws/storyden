package role_repo

import (
	"context"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

var ErrWritePermissionsNotAllowed = fault.New("write permissions not allowed on guest role")

type Writer interface {
	Create(ctx context.Context, name string, colour string, perms rbac.PermissionList, opts ...Mutation) (*role.Role, error)
	Update(ctx context.Context, id role.RoleID, opts ...Mutation) (*role.Role, error)
	Delete(ctx context.Context, id role.RoleID) error
	UpdateSortOrder(ctx context.Context, ids []role.RoleID) error
}

type Mutation func(*ent.RoleMutation)

func WithName(name string) Mutation {
	return func(m *ent.RoleMutation) {
		m.SetName(name)
	}
}

func WithColour(colour string) Mutation {
	return func(m *ent.RoleMutation) {
		m.SetColour(colour)
	}
}

func WithPermissions(perms rbac.PermissionList) Mutation {
	ps := dt.Map(perms, func(p rbac.Permission) string { return p.String() })
	return func(m *ent.RoleMutation) {
		m.SetPermissions(ps)
	}
}

func WithMeta(meta map[string]any) Mutation {
	return func(m *ent.RoleMutation) {
		m.SetMetadata(meta)
	}
}

func (w *Repository) Create(ctx context.Context, name string, colour string, perms rbac.PermissionList, opts ...Mutation) (*role.Role, error) {
	ps := dt.Map(perms, func(p rbac.Permission) string { return p.String() })
	nextSortKey, err := w.nextCustomSortKey(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	create := w.db.Role.Create().
		SetName(name).
		SetColour(colour).
		SetPermissions(ps).
		SetSortKey(nextSortKey)

	mutation := create.Mutation()
	for _, opt := range opts {
		opt(mutation)
	}

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.storeRole(ctx, rl, true); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.create.store_role", err); recoveryErr != nil {
			return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	if err := w.syncCustomRoles(ctx); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.create.sync_custom_roles", err); recoveryErr != nil {
			return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	return rl, nil
}

func (w *Repository) Update(ctx context.Context, id role.RoleID, opts ...Mutation) (*role.Role, error) {
	if id == role.DefaultRoleMemberID {
		rl, err := w.updateDefaultRole(ctx, opts...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := w.storeRole(ctx, rl, true); err != nil {
			if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.update_default_member.store_role", err); recoveryErr != nil {
				return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
			}
		}

		return rl, nil
	}

	if id == role.DefaultRoleGuestID {
		rl, err := w.updateGuestRole(ctx, opts...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := w.storeRole(ctx, rl, true); err != nil {
			if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.update_default_guest.store_role", err); recoveryErr != nil {
				return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
			}
		}

		return rl, nil
	}

	update := w.db.Role.UpdateOneID(xid.ID(id))
	mutation := update.Mutation()

	for _, opt := range opts {
		opt(mutation)
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.storeRole(ctx, rl, true); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.update.store_role", err); recoveryErr != nil {
			return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	if err := w.syncCustomRoles(ctx); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.update.sync_custom_roles", err); recoveryErr != nil {
			return nil, fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	return rl, nil
}

func (w *Repository) updateDefaultRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := w.lookupRole(ctx, role.DefaultRoleMemberID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := w.db.Role.Create()
		mutate := create.Mutation()

		// The default Member role has a hard-coded ID.
		mutate.SetID(xid.ID(role.DefaultRoleMemberID))
		mutate.SetName("Member")
		mutate.SetSortKey(-1)

		for _, opt := range opts {
			opt(mutate)
		}

		r, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return role.Map(r)
	}

	update := rl.Update()
	mutate := update.Mutation()
	mutate.SetSortKey(role.DefaultRoleMember.SortKey)
	for _, opt := range opts {
		opt(mutate)
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return role.Map(r)
}

func (w *Repository) updateGuestRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := w.lookupRole(ctx, role.DefaultRoleGuestID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := w.db.Role.Create()
		mutate := create.Mutation()

		// The default Guest role has a hard-coded ID.
		mutate.SetID(xid.ID(role.DefaultRoleGuestID))
		mutate.SetName("Guest")
		mutate.SetSortKey(-2)

		for _, opt := range opts {
			opt(mutate)
		}

		if perms, ok := mutate.Permissions(); ok {
			// Do not allow write permissions to be added.
			list, err := rbac.NewPermissions(perms)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			if list.HasAnyWrite() {
				return nil, fault.Wrap(ErrWritePermissionsNotAllowed, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
		}

		r, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return role.Map(r)
	}

	update := rl.Update()
	mutate := update.Mutation()
	mutate.SetSortKey(role.DefaultRoleGuest.SortKey)
	for _, opt := range opts {
		opt(mutate)
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return role.Map(r)
}

func (w *Repository) lookupRole(ctx context.Context, id role.RoleID) (*ent.Role, bool, error) {
	r, err := w.db.Role.Query().Where(ent_role.ID(xid.ID(id))).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return r, true, nil
}

func (w *Repository) Delete(ctx context.Context, id role.RoleID) error {
	err := w.db.Role.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.deleteRole(ctx, id); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.delete.delete_role", err); recoveryErr != nil {
			return fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	if err := w.syncCustomRoles(ctx); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.delete.sync_custom_roles", err); recoveryErr != nil {
			return fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	return nil
}

func (w *Repository) UpdateSortOrder(ctx context.Context, ids []role.RoleID) error {
	for _, id := range ids {
		if id == role.DefaultRoleGuestID || id == role.DefaultRoleMemberID || id == role.DefaultRoleAdminID {
			return fault.New("default roles cannot be reordered", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		_ = tx.Rollback()
	}()

	customRoles, err := tx.Role.Query().Where(ent_role.IDNotIn(
		xid.ID(role.DefaultRoleGuestID),
		xid.ID(role.DefaultRoleMemberID),
		xid.ID(role.DefaultRoleAdminID),
	)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if len(customRoles) != len(ids) {
		return fault.New(
			"role reorder list must include all custom roles exactly once",
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.With("role_ids must contain every custom role ID once"),
		)
	}

	existing := make(map[role.RoleID]struct{}, len(customRoles))
	for _, r := range customRoles {
		existing[role.RoleID(r.ID)] = struct{}{}
	}

	seen := make(map[role.RoleID]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := existing[id]; !ok {
			return fault.New("unknown custom role ID in role reorder list", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if _, ok := seen[id]; ok {
			return fault.New("duplicate custom role ID in role reorder list", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		seen[id] = struct{}{}
	}

	for i, id := range ids {
		_, err := tx.Role.UpdateOneID(xid.ID(id)).SetSortKey(float64(i)).Save(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := tx.Commit(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.syncCustomRoles(ctx); err != nil {
		if recoveryErr := w.recoverFromCacheWriteFailure(ctx, "writer.update_sort_order.sync_custom_roles", err); recoveryErr != nil {
			return fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	return nil
}

func (w *Repository) nextCustomSortKey(ctx context.Context) (float64, error) {
	customRoles, err := w.db.Role.Query().Where(ent_role.IDNotIn(
		xid.ID(role.DefaultRoleGuestID),
		xid.ID(role.DefaultRoleMemberID),
		xid.ID(role.DefaultRoleAdminID),
	)).All(ctx)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	if len(customRoles) == 0 {
		return 0, nil
	}

	sort.Slice(customRoles, func(i, j int) bool {
		return customRoles[i].SortKey < customRoles[j].SortKey
	})

	return customRoles[len(customRoles)-1].SortKey + 1, nil
}

func (w *Repository) syncCustomRoles(ctx context.Context) error {
	_, err := w.listFromDB(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
