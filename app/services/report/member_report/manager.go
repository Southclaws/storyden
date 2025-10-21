package member_report

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/resources/report/report_querier"
	"github.com/Southclaws/storyden/app/resources/report/report_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
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

func (m *Manager) Submit(
	ctx context.Context,
	targetID xid.ID,
	targetKind datagraph.Kind,
	comment opt.Optional[string],
) (*report.Report, error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rep, err := m.reportWriter.Create(ctx, targetID, targetKind, opt.New(acc.ID), comment)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var targetRef *datagraph.Ref
	if rep.TargetItem != nil {
		targetRef = datagraph.NewRef(rep.TargetItem)
	}

	m.bus.Publish(ctx, &rpc.EventReportCreated{
		ID:         rep.ID,
		Target:     opt.Map(opt.NewPtr(targetRef), rpc.DatagraphRefToRPC),
		ReportedBy: opt.Map(rep.ReportedBy, func(a account.Account) account.AccountID { return a.ID }),
	})

	return rep, nil
}

func (m *Manager) List(
	ctx context.Context,
	page pagination.Parameters,
) (pagination.Result[*report.Report], error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	reports, err := m.reportQuerier.List(
		ctx,
		page,
		report_querier.WithReporter(acc.ID),
	)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return reports, nil
}

func (m *Manager) Resolve(
	ctx context.Context,
	reportID report.ID,
) (*report.Report, error) {
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	existing, err := m.reportQuerier.Get(ctx, reportID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if existing.ReportedBy.OrZero().ID != acc.ID {
		return nil, fault.Wrap(
			fault.New("cannot resolve report submitted by another user"),
			fctx.With(ctx),
			ftag.With(ftag.PermissionDenied),
		)
	}

	rep, err := m.reportWriter.Update(
		ctx,
		reportID,
		opt.New(report.StatusResolved),
		opt.NewEmpty[account.AccountID](),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var targetRef *datagraph.Ref
	if rep.TargetItem != nil {
		targetRef = datagraph.NewRef(rep.TargetItem)
	}

	m.bus.Publish(ctx, &rpc.EventReportUpdated{
		ID:     rep.ID,
		Target: opt.Map(opt.NewPtr(targetRef), rpc.DatagraphRefToRPC),
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
