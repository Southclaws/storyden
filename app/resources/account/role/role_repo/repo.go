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
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

var (
	ErrWritePermissionsNotAllowed = fault.New("write permissions not allowed on guest role")
	ErrAdminPermissionsNotAllowed = fault.New("cannot set permissions on admin default role")
)

type Repository struct {
	db    *ent.Client
	store cache.Store
}

func New(db *ent.Client, store cache.Store) *Repository {
	return &Repository{
		db:    db,
		store: store,
	}
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

func (r *Repository) Create(ctx context.Context, name string, colour string, perms rbac.PermissionList, opts ...Mutation) (*role.Role, error) {
	ps := dt.Map(perms, func(p rbac.Permission) string { return p.String() })
	nextSortKey, err := r.nextCustomSortKey(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	create := r.db.Role.Create().
		SetName(name).
		SetColour(colour).
		SetPermissions(ps).
		SetSortKey(nextSortKey)

	mutation := create.Mutation()
	for _, opt := range opts {
		opt(mutation)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleCache(ctx, rl)

	return rl, nil
}

func (r *Repository) Update(ctx context.Context, id role.RoleID, opts ...Mutation) (*role.Role, error) {
	if id == role.DefaultRoleMemberID {
		return r.updateDefaultRole(ctx, opts...)
	}

	if id == role.DefaultRoleGuestID {
		return r.updateGuestRole(ctx, opts...)
	}

	if id == role.DefaultRoleAdminID {
		return r.updateAdminRole(ctx, opts...)
	}

	update := r.db.Role.UpdateOneID(xid.ID(id))
	mutation := update.Mutation()

	for _, opt := range opts {
		opt(mutation)
	}

	row, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleCache(ctx, rl)

	return rl, nil
}

func (r *Repository) Get(ctx context.Context, id role.RoleID) (*role.Role, error) {
	mapped, err := r.GetMany(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, ok := mapped[id]
	if !ok {
		return nil, fault.Wrap(fault.New("role not found"), fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	return rl, nil
}

func (r *Repository) List(ctx context.Context) (role.Roles, error) {
	if cached, ok := r.getRoleListCache(ctx); ok {
		return cached, nil
	}

	loaded, err := r.listFromDB(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleListCache(ctx, loaded)

	return loaded, nil
}

func (r *Repository) GetMany(ctx context.Context, ids ...role.RoleID) (map[role.RoleID]*role.Role, error) {
	if len(ids) == 0 {
		return map[role.RoleID]*role.Role{}, nil
	}

	byID := make(map[role.RoleID]*role.Role, len(ids))
	seen := map[role.RoleID]struct{}{}
	missing := make([]role.RoleID, 0, len(ids))

	var cachedByID map[role.RoleID]*role.Role
	if cached, ok := r.getRoleListCache(ctx); ok {
		cachedByID = make(map[role.RoleID]*role.Role, len(cached))
		for _, rl := range cached {
			cachedByID[rl.ID] = rl
		}
	}

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		if cachedByID != nil {
			if rl, ok := cachedByID[id]; ok {
				byID[id] = rl
				continue
			}
		}

		missing = append(missing, id)
	}

	if len(missing) == 0 {
		return byID, nil
	}

	rows, err := r.db.Role.Query().
		Where(ent_role.IDIn(dt.Map(missing, func(i role.RoleID) xid.ID { return xid.ID(i) })...)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mappedRows, err := role.MapList(rows)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	dbByID := make(map[role.RoleID]*role.Role, len(mappedRows))
	for _, rl := range mappedRows {
		dbByID[rl.ID] = rl
	}

	for _, id := range missing {
		if rl, ok := dbByID[id]; ok {
			byID[id] = rl
			continue
		}

		if defaultRole := resolveDefaultRole(id); defaultRole != nil {
			copy := *defaultRole
			byID[id] = &copy
		}
	}

	return byID, nil
}

func (r *Repository) GetMemberRole(ctx context.Context) (*role.Role, error) {
	return r.getOrDefaultRole(ctx, role.DefaultRoleMemberID)
}

func (r *Repository) GetGuestRole(ctx context.Context) (*role.Role, error) {
	return r.getOrDefaultRole(ctx, role.DefaultRoleGuestID)
}

func (r *Repository) GetAdminRole(ctx context.Context) (*role.Role, error) {
	return r.getOrDefaultRole(ctx, role.DefaultRoleAdminID)
}

func (r *Repository) getOrDefaultRole(ctx context.Context, id role.RoleID) (*role.Role, error) {
	if defaultRole := resolveDefaultRole(id); defaultRole != nil {
		copy := *defaultRole

		if cached, ok := r.getRoleCache(ctx, id); ok {
			return cached, nil
		}

		loaded, err := r.listFromDB(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		r.setRoleListCache(ctx, loaded)

		for _, rl := range loaded {
			if rl.ID == id {
				return rl, nil
			}
		}

		// Should be unreachable because listFromDB always injects defaults.
		return &copy, nil
	}

	result, err := r.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) Delete(ctx context.Context, id role.RoleID) error {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.Role.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.removeRoleCacheStrict(ctx, id); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		r.refreshRoleListCache(ctx)
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) UpdateSortOrder(ctx context.Context, ids []role.RoleID) error {
	for _, id := range ids {
		if id == role.DefaultRoleGuestID || id == role.DefaultRoleMemberID || id == role.DefaultRoleAdminID {
			return fault.New("default roles cannot be reordered", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	customRoles, err := r.db.Role.Query().Where(ent_role.IDNotIn(
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
	for _, row := range customRoles {
		existing[role.RoleID(row.ID)] = struct{}{}
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

	tx, err := r.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer func() { _ = tx.Rollback() }()

	for i, id := range ids {
		sortKey := float64(len(ids) - i)
		if _, err := tx.Role.UpdateOneID(xid.ID(id)).SetSortKey(sortKey).Save(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

	}

	if err := tx.Commit(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	r.refreshRoleListCache(ctx)

	return nil
}

func (r *Repository) updateDefaultRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := r.lookupRole(ctx, role.DefaultRoleMemberID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := r.db.Role.Create()
		mutate := create.Mutation()

		mutate.SetID(xid.ID(role.DefaultRoleMemberID))
		mutate.SetName("Member")
		mutate.SetSortKey(-1)

		for _, opt := range opts {
			opt(mutate)
		}

		row, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		mapped, err := role.Map(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		r.setRoleCache(ctx, mapped)

		return mapped, nil
	}

	update := rl.Update()
	mutate := update.Mutation()
	mutate.SetSortKey(role.DefaultRoleMember.SortKey)
	for _, opt := range opts {
		opt(mutate)
	}

	row, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.Map(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleCache(ctx, mapped)

	return mapped, nil
}

func (r *Repository) updateAdminRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := r.lookupRole(ctx, role.DefaultRoleAdminID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := r.db.Role.Create()
		mutate := create.Mutation()

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

		row, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		mapped, err := role.Map(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		r.setRoleCache(ctx, mapped)

		return mapped, nil
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

	row, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.Map(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleCache(ctx, mapped)

	return mapped, nil
}

func (r *Repository) updateGuestRole(ctx context.Context, opts ...Mutation) (*role.Role, error) {
	rl, found, err := r.lookupRole(ctx, role.DefaultRoleGuestID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		create := r.db.Role.Create()
		mutate := create.Mutation()

		mutate.SetID(xid.ID(role.DefaultRoleGuestID))
		mutate.SetName("Guest")
		mutate.SetSortKey(-2)

		for _, opt := range opts {
			opt(mutate)
		}

		if perms, ok := mutate.Permissions(); ok {
			list, _ := rbac.NewPermissions(perms)
			if list.HasAnyWrite() {
				return nil, fault.Wrap(ErrWritePermissionsNotAllowed, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
		}

		row, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		mapped, err := role.Map(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		r.setRoleCache(ctx, mapped)

		return mapped, nil
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

	row, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.Map(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.setRoleCache(ctx, mapped)

	return mapped, nil
}

func (r *Repository) lookupRole(ctx context.Context, id role.RoleID) (*ent.Role, bool, error) {
	row, err := r.db.Role.Query().Where(ent_role.ID(xid.ID(id))).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return row, true, nil
}

func (r *Repository) nextCustomSortKey(ctx context.Context) (float64, error) {
	customRoles, err := r.db.Role.Query().Where(ent_role.IDNotIn(
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

func adminPermissionsMutationSubmitted(mutate *ent.RoleMutation) bool {
	if _, ok := mutate.Permissions(); ok {
		return true
	}

	if perms, ok := mutate.AppendedPermissions(); ok && len(perms) > 0 {
		return true
	}

	return false
}

func resolveDefaultRole(id role.RoleID) *role.Role {
	switch id {
	case role.DefaultRoleGuestID:
		copy := role.DefaultRoleGuest
		return &copy
	case role.DefaultRoleMemberID:
		copy := role.DefaultRoleMember
		return &copy
	case role.DefaultRoleAdminID:
		copy := role.DefaultRoleAdmin
		return &copy
	default:
		return nil
	}
}
