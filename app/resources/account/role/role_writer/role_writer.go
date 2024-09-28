package role_writer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
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

func (w *Writer) Create(ctx context.Context, name string, colour string, perms rbac.PermissionList) (*role.Role, error) {
	ps := dt.Map(perms, func(p rbac.Permission) string { return p.String() })

	r, err := w.db.Role.Create().
		SetName(name).
		SetColour(colour).
		SetPermissions(ps).
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

func (w *Writer) Delete(ctx context.Context, id role.RoleID) error {
	err := w.db.Role.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
