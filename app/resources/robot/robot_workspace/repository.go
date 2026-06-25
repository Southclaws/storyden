package robot_workspace

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_workspace "github.com/Southclaws/storyden/internal/ent/robotworkspace"
	ent_robot_workspace_instance "github.com/Southclaws/storyden/internal/ent/robotworkspaceinstance"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

type WorkspaceOption func(*ent.RobotWorkspaceMutation)

func WithName(v string) WorkspaceOption {
	return func(m *ent.RobotWorkspaceMutation) {
		m.SetName(v)
	}
}

func WithDescription(v string) WorkspaceOption {
	return func(m *ent.RobotWorkspaceMutation) {
		m.SetDescription(v)
	}
}

func WithConfig(v map[string]any) WorkspaceOption {
	return func(m *ent.RobotWorkspaceMutation) {
		m.SetConfig(v)
	}
}

func WithMetadata(v map[string]any) WorkspaceOption {
	return func(m *ent.RobotWorkspaceMutation) {
		m.SetMetadata(v)
	}
}

func (r *Repository) Create(
	ctx context.Context,
	name string,
	description string,
	provider robot.WorkspaceProvider,
	creatorID account.AccountID,
	opts ...WorkspaceOption,
) (*robot.Workspace, error) {
	create := r.db.RobotWorkspace.Create()
	mutate := create.Mutation()
	mutate.SetName(name)
	mutate.SetDescription(description)
	mutate.SetProvider(ent_robot_workspace.Provider(provider))
	mutate.SetCreatedBy(xid.ID(creatorID))

	for _, opt := range opts {
		opt(mutate)
	}

	created, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.Get(ctx, robot.WorkspaceID(created.ID))
}

func (r *Repository) List(ctx context.Context, params pagination.Parameters) (*pagination.Result[*robot.Workspace], error) {
	query := r.db.RobotWorkspace.Query().
		WithCreator().
		Order(
			ent_robot_workspace.ByCreatedAt(sql.OrderDesc()),
			ent_robot_workspace.ByID(sql.OrderDesc()),
		)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query.Limit(params.Limit()).Offset(params.Offset())

	rows, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	workspaces, err := dt.MapErr(rows, robot.MapWorkspace)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, workspaces)
	return &result, nil
}

func (r *Repository) Get(ctx context.Context, id robot.WorkspaceID) (*robot.Workspace, error) {
	row, err := r.db.RobotWorkspace.Query().
		Where(ent_robot_workspace.IDEQ(xid.ID(id))).
		WithCreator().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.MapWorkspace(row)
}

func (r *Repository) Update(ctx context.Context, id robot.WorkspaceID, opts ...WorkspaceOption) (*robot.Workspace, error) {
	update := r.db.RobotWorkspace.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, opt := range opts {
		opt(mutate)
	}

	if err := update.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.Get(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id robot.WorkspaceID) error {
	err := r.db.RobotWorkspace.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}
	return nil
}

func (r *Repository) CreateInstance(
	ctx context.Context,
	workspaceID robot.WorkspaceID,
	creatorID account.AccountID,
	providerState map[string]any,
	metadata map[string]any,
) (*robot.WorkspaceInstance, error) {
	workspace, err := r.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	created, err := r.db.RobotWorkspaceInstance.Create().
		SetWorkspaceID(xid.ID(workspaceID)).
		SetCreatedBy(xid.ID(creatorID)).
		SetProvider(ent_robot_workspace_instance.Provider(workspace.Provider)).
		SetProviderState(providerState).
		SetMetadata(metadata).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.GetInstance(ctx, robot.WorkspaceInstanceID(created.ID))
}

func (r *Repository) ListInstances(ctx context.Context, params pagination.Parameters) (*pagination.Result[*robot.WorkspaceInstance], error) {
	query := r.db.RobotWorkspaceInstance.Query().
		WithCreator().
		Order(
			ent_robot_workspace_instance.ByCreatedAt(sql.OrderDesc()),
			ent_robot_workspace_instance.ByID(sql.OrderDesc()),
		)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query.Limit(params.Limit()).Offset(params.Offset())

	rows, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	instances, err := dt.MapErr(rows, robot.MapWorkspaceInstance)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, instances)
	return &result, nil
}

func (r *Repository) GetInstance(ctx context.Context, id robot.WorkspaceInstanceID) (*robot.WorkspaceInstance, error) {
	row, err := r.db.RobotWorkspaceInstance.Query().
		Where(ent_robot_workspace_instance.IDEQ(xid.ID(id))).
		WithCreator().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.MapWorkspaceInstance(row)
}

func (r *Repository) UpdateInstanceProviderState(ctx context.Context, id robot.WorkspaceInstanceID, providerState map[string]any) (*robot.WorkspaceInstance, error) {
	if err := r.db.RobotWorkspaceInstance.UpdateOneID(xid.ID(id)).
		SetProviderState(providerState).
		Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.GetInstance(ctx, id)
}

func (r *Repository) DeleteInstance(ctx context.Context, id robot.WorkspaceInstanceID) error {
	err := r.db.RobotWorkspaceInstance.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}
	return nil
}
