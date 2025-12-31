package audit_querier

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	ent_auditlog "github.com/Southclaws/storyden/internal/ent/auditlog"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

type Filter struct {
	Types     opt.Optional[[]audit.EventType]
	TimeRange opt.Optional[TimeRange]
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (q *Querier) List(
	ctx context.Context,
	page pagination.Parameters,
	filter Filter,
) (*pagination.Result[*audit.AuditLog], error) {
	query := q.db.AuditLog.Query().
		WithEnactedBy()

	filter.Types.Call(func(types []audit.EventType) {
		typeStrings := dt.Map(types, func(t audit.EventType) string {
			return t.String()
		})
		query.Where(ent_auditlog.TypeIn(typeStrings...))
	})

	filter.TimeRange.Call(func(tr TimeRange) {
		query.Where(
			ent_auditlog.CreatedAtGTE(tr.Start),
			ent_auditlog.CreatedAtLTE(tr.End),
		)
	})

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	results, err := query.
		Order(ent_auditlog.ByCreatedAt(sql.OrderDesc())).
		Limit(page.Limit()).
		Offset(page.Offset()).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	logs, err := dt.MapErr(results, audit.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(page, total, logs)

	return &result, nil
}

func (q *Querier) Get(ctx context.Context, id audit.AuditLogID) (*audit.AuditLog, error) {
	al, err := q.db.AuditLog.Query().
		Where(ent_auditlog.IDEQ(xid.ID(id))).
		WithEnactedBy().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return audit.Map(al)
}
