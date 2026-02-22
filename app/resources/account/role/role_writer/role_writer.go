package role_writer

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

var (
	ErrWritePermissionsNotAllowed = fault.New("write permissions not allowed on guest role")
	ErrAdminPermissionsNotAllowed = fault.New("cannot set permissions on admin default role")
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
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

func (w *Writer) Create(ctx context.Context, name string, colour string, perms rbac.PermissionList, opts ...Mutation) (*role.Role, error) {
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

	return rl, nil
}

func (w *Writer) Update(ctx context.Context, id role.RoleID, opts ...Mutation) (*role.Role, error) {
	if id == role.DefaultRoleMemberID {
		return w.updateDefaultRole(ctx, opts...)
	}

	if id == role.DefaultRoleGuestID {
		return w.updateGuestRole(ctx, opts...)
	}

	if id == role.DefaultRoleAdminID {
		return w.updateAdminRole(ctx, opts...)
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

	return rl, nil
}

func (w *Writer) updateDefaultRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
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

func (w *Writer) updateAdminRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := w.lookupRole(ctx, role.DefaultRoleAdminID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := w.db.Role.Create()
		mutate := create.Mutation()

		// The default Admin role has a hard-coded ID.
		mutate.SetID(xid.ID(role.DefaultRoleAdminID))
		mutate.SetName("Admin")
		mutate.SetSortKey(role.DefaultRoleAdmin.SortKey)

		for _, opt := range opts {
			opt(mutate)
		}

		if adminPermissionsMutationSubmitted(mutate) {
			return nil, fault.Wrap(ErrAdminPermissionsNotAllowed, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		mutate.SetPermissions(dt.Map(role.DefaultRoleAdmin.Permissions.List(), func(p rbac.Permission) string {
			return p.String()
		}))

		r, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return role.Map(r)
	}

	update := rl.Update()
	mutate := update.Mutation()
	mutate.SetSortKey(role.DefaultRoleAdmin.SortKey)
	for _, opt := range opts {
		opt(mutate)
	}

	if adminPermissionsMutationSubmitted(mutate) {
		return nil, fault.Wrap(ErrAdminPermissionsNotAllowed, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return role.Map(r)
}

func adminPermissionsMutationSubmitted(mutate *ent.RoleMutation) bool {
	if _, ok := mutate.Permissions(); ok {
		return true
	}

	if perms, ok := mutate.AppendedPermissions(); ok && len(perms) > 0 {
		return true
	}

	return false
}

func (w *Writer) updateGuestRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
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
			list, _ := rbac.NewPermissions(perms)
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

	if perms, ok := mutate.Permissions(); ok {
		list, _ := rbac.NewPermissions(perms)
		if list.HasAnyWrite() {
			return nil, fault.Wrap(ErrWritePermissionsNotAllowed, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return role.Map(r)
}

func (w *Writer) lookupRole(ctx context.Context, id role.RoleID) (*ent.Role, bool, error) {
	r, err := w.db.Role.Query().Where(ent_role.ID(xid.ID(id))).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return r, true, nil
}

func (w *Writer) Delete(ctx context.Context, id role.RoleID) error {
	err := w.db.Role.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) UpdateSortOrder(ctx context.Context, ids []role.RoleID) error {
	for _, id := range ids {
		if id == role.DefaultRoleGuestID || id == role.DefaultRoleMemberID || id == role.DefaultRoleAdminID {
			return fault.New("default roles cannot be reordered", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	customRoles, err := w.db.Role.Query().Where(ent_role.IDNotIn(
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

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for i, id := range ids {
		_, err := tx.Role.UpdateOneID(xid.ID(id)).SetSortKey(float64(i)).Save(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := tx.Commit(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) nextCustomSortKey(ctx context.Context) (float64, error) {
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
