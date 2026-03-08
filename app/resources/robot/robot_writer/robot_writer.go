package robot_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db      *ent.Client
	querier *robot_querier.Querier
}

func New(db *ent.Client, querier *robot_querier.Querier) *Writer {
	return &Writer{
		db:      db,
		querier: querier,
	}
}

type Option func(*ent.RobotMutation)

func WithName(v string) Option {
	return func(m *ent.RobotMutation) {
		m.SetName(v)
	}
}

func WithDescription(v string) Option {
	return func(m *ent.RobotMutation) {
		m.SetDescription(v)
	}
}

func WithPlaybook(v string) Option {
	return func(m *ent.RobotMutation) {
		m.SetPlaybook(v)
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.RobotMutation) {
		m.SetMetadata(meta)
	}
}

func (w *Writer) Create(
	ctx context.Context,
	name string,
	description string,
	playbook string,
	authorID account.AccountID,
	opts ...Option,
) (*robot.Robot, error) {
	create := w.db.Robot.Create()
	mutate := create.Mutation()

	mutate.SetName(name)
	mutate.SetDescription(description)
	mutate.SetPlaybook(playbook)
	mutate.SetAuthorID(xid.ID(authorID))

	for _, fn := range opts {
		fn(mutate)
	}

	r, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return w.querier.Get(ctx, robot_ref.ID(r.ID))
}

func (w *Writer) Update(
	ctx context.Context,
	id robot_ref.ID,
	opts ...Option,
) (*robot.Robot, error) {
	update := w.db.Robot.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	err := update.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return w.querier.Get(ctx, robot_ref.ID(id))
}

func (w *Writer) Delete(ctx context.Context, id robot_ref.ID) error {
	err := w.db.Robot.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}
	return nil
}
