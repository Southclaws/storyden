package role_writer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

var ErrWritePermissionsNotAllowed = fault.New("write permissions not allowed on guest role")

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

func (w *Writer) Create(ctx context.Context, name string, colour string, perms rbac.PermissionList) (*role.Role, error) {
	ps := dt.Map(perms, func(p rbac.Permission) string { return p.String() })

	r, err := w.db.Role.Create().
		SetName(name).
		SetColour(colour).
		SetPermissions(ps).
		SetSortKey(0.0).
		Save(ctx)
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
	for _, opt := range opts {
		opt(mutate)
	}

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return role.Map(r)
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
	for _, opt := range opts {
		opt(mutate)
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
