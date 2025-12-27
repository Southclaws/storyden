package report_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/resources/report/report_querier"
	"github.com/Southclaws/storyden/app/resources/report/report_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Manager struct {
	reportQuerier *report_querier.Querier
	reportWriter  *report_writer.Writer
	bus           *pubsub.Bus
}

func New(
	reportQuerier *report_querier.Querier,
	reportWriter *report_writer.Writer,
	bus *pubsub.Bus,
) *Manager {
	return &Manager{
		reportQuerier: reportQuerier,
		reportWriter:  reportWriter,
		bus:           bus,
	}
}

type ListOpts struct {
	Status opt.Optional[[]report.Status]
	Kind   opt.Optional[[]datagraph.Kind]
}

func (m *Manager) List(
	ctx context.Context,
	page pagination.Parameters,
	opts ListOpts,
) (pagination.Result[*report.Report], error) {
	queryOpts := []report_querier.Query{}

	opts.Status.Call(func(statuses []report.Status) {
		queryOpts = append(queryOpts, report_querier.WithStatus(statuses...))
	})

	opts.Kind.Call(func(kinds []datagraph.Kind) {
		queryOpts = append(queryOpts, report_querier.WithKind(kinds...))
	})

	reports, err := m.reportQuerier.List(ctx, page, queryOpts...)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return reports, nil
}

func (m *Manager) Get(
	ctx context.Context,
	reportID report.ID,
) (*report.Report, error) {
	r, err := m.reportQuerier.Get(ctx, reportID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

type UpdateOpts struct {
	Status    opt.Optional[report.Status]
	HandledBy opt.Optional[account.AccountID]
}

func (m *Manager) Update(
	ctx context.Context,
	reportID report.ID,
	opts UpdateOpts,
) (*report.Report, error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	handledBy := opts.HandledBy

	if opts.Status.Ok() && !opts.HandledBy.Ok() {
		handledBy = opt.New(acc.ID)
	}

	rep, err := m.reportWriter.Update(ctx, reportID, opts.Status, handledBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var targetRef *datagraph.Ref
	if rep.TargetItem != nil {
		targetRef = datagraph.NewRef(rep.TargetItem)
	}

	m.bus.Publish(ctx, &message.EventReportUpdated{
		ID:     rep.ID,
		Target: targetRef,
		ReportedBy: opt.Map(rep.ReportedBy, func(a account.Account) account.AccountID {
			return a.ID
		}),
		HandledBy: opt.Map(rep.HandledBy, func(a account.Account) account.AccountID {
			return a.ID
		}),
		Status: rep.Status,
	})

	return rep, nil
}
