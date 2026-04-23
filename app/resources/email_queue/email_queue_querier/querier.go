package email_queue_querier

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/email_queue"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	ent_emailqueue "github.com/Southclaws/storyden/internal/ent/emailqueue"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

type Filter struct {
	// TODO: Add date range parameters for filtering.
	Statuses opt.Optional[[]email_queue.Status]
}

func (q *Querier) List(
	ctx context.Context,
	page pagination.Parameters,
	filter Filter,
) (*pagination.Result[*email_queue.Email], error) {
	query := q.db.EmailQueue.Query()

	filter.Statuses.Call(func(statuses []email_queue.Status) {
		statusValues := dt.Map(statuses, func(s email_queue.Status) ent_emailqueue.Status {
			return ent_emailqueue.Status(s.String())
		})
		query.Where(ent_emailqueue.StatusIn(statusValues...))
	})

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	results, err := query.
		Order(ent_emailqueue.ByCreatedAt(sql.OrderDesc())).
		Limit(page.Limit()).
		Offset(page.Offset()).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	logs, err := dt.MapErr(results, email_queue.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(page, total, logs)
	return &result, nil
}
