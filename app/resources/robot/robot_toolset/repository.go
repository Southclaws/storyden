package robot_toolset

import (
	"context"
	"slices"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot "github.com/Southclaws/storyden/internal/ent/robot"
	ent_toolset "github.com/Southclaws/storyden/internal/ent/robottoolset"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]*Toolset, error) {
	rows, err := r.db.RobotToolset.Query().
		WithAuthor().
		Order(ent_toolset.ByName(), ent_toolset.ByID()).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	toolsets, err := dt.MapErr(rows, Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return toolsets, nil
}

func (r *Repository) Get(ctx context.Context, id ID) (*Toolset, error) {
	row, err := r.db.RobotToolset.Query().
		Where(ent_toolset.IDEQ(xid.ID(id))).
		WithAuthor().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return Map(row)
}

func (r *Repository) Create(
	ctx context.Context,
	authorID account.AccountID,
	name, description, instruction string,
	tools []string,
	metadata map[string]any,
) (*Toolset, error) {
	row, err := r.db.RobotToolset.Create().
		SetAuthorID(xid.ID(authorID)).
		SetName(name).
		SetDescription(description).
		SetInstruction(instruction).
		SetTools(tools).
		SetMetadata(metadata).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return r.Get(ctx, ID(row.ID))
}

type Update struct {
	Name        *string
	Description *string
	Instruction *string
	Tools       *[]string
	Metadata    *map[string]any
}

func (r *Repository) Update(ctx context.Context, id ID, patch Update) (*Toolset, error) {
	update := r.db.RobotToolset.UpdateOneID(xid.ID(id))
	if patch.Name != nil {
		update.SetName(*patch.Name)
	}
	if patch.Description != nil {
		update.SetDescription(*patch.Description)
	}
	if patch.Instruction != nil {
		update.SetInstruction(*patch.Instruction)
	}
	if patch.Tools != nil {
		update.SetTools(*patch.Tools)
	}
	if patch.Metadata != nil {
		update.SetMetadata(*patch.Metadata)
	}
	if err := update.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		} else if ent.IsConstraintError(err) {
			err = fault.Wrap(err, ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return r.Get(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id ID) error {
	usage, err := r.UsageCount(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if usage > 0 {
		return fault.New("toolset is still assigned to one or more robots", fctx.With(ctx), ftag.With(ftag.AlreadyExists))
	}

	if err := r.db.RobotToolset.DeleteOneID(xid.ID(id)).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (r *Repository) UsageCount(ctx context.Context, id ID) (int, error) {
	robots, err := r.db.Robot.Query().
		Select(ent_robot.FieldToolsets).
		All(ctx)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}
	count := 0
	for _, robot := range robots {
		if slices.Contains(robot.Toolsets, id.String()) {
			count++
		}
	}
	return count, nil
}

func (r *Repository) IsAuthor(ctx context.Context, id ID, accountID account.AccountID) (bool, error) {
	exists, err := r.db.RobotToolset.Query().
		Where(
			ent_toolset.IDEQ(xid.ID(id)),
			ent_toolset.AuthorIDEQ(xid.ID(accountID)),
		).
		Exist(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return exists, nil
}
