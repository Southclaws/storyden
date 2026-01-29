package robot_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot "github.com/Southclaws/storyden/internal/ent/robot"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) List(ctx context.Context, params pagination.Parameters) (*pagination.Result[*robot.Robot], error) {
	query := q.db.Robot.Query().
		WithAuthor()

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query.Limit(params.Limit()).Offset(params.Offset())

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	robots, err := dt.MapErr(r, robot.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, robots)
	return &result, nil
}

func (q *Querier) Get(ctx context.Context, id robot_ref.ID) (*robot.Robot, error) {
	r, err := q.db.Robot.Query().
		Where(ent_robot.IDEQ(xid.ID(id))).
		WithAuthor().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.Map(r)
}

func (q *Querier) GetByName(ctx context.Context, name string) (*robot.Robot, error) {
	r, err := q.db.Robot.Query().
		Where(ent_robot.NameEQ(name)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.Map(r)
}
