package report_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/resources/report/report_querier"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db      *ent.Client
	querier *report_querier.Querier
}

func New(db *ent.Client, querier *report_querier.Querier) *Writer {
	return &Writer{db: db, querier: querier}
}

func (w *Writer) Create(
	ctx context.Context,
	targetID xid.ID,
	targetKind datagraph.Kind,
	reportedBy opt.Optional[account.AccountID],
	comment opt.Optional[string],
) (*report.Report, error) {
	create := w.db.Report.Create().
		SetTargetID(targetID).
		SetTargetKind(targetKind.String()).
		SetStatus(report.StatusSubmitted.String())

	reportedBy.Call(func(value account.AccountID) { create.SetReportedByID(xid.ID(value)) })
	comment.Call(func(value string) { create.SetComment(value) })

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, report.ID(r.ID))
}

func (w *Writer) Update(
	ctx context.Context,
	id report.ID,
	status opt.Optional[report.Status],
	handledBy opt.Optional[account.AccountID],
) (*report.Report, error) {
	update := w.db.Report.UpdateOneID(xid.ID(id))

	status.Call(func(value report.Status) { update.SetStatus(value.String()) })
	handledBy.Call(func(value account.AccountID) { update.SetHandledByID(xid.ID(value)) })

	r, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, report.ID(r.ID))
}
