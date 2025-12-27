package system_report

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/resources/report/report_writer"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Manager struct {
	reportWriter *report_writer.Writer
	bus          *pubsub.Bus
}

func New(
	reportWriter *report_writer.Writer,
	bus *pubsub.Bus,
) *Manager {
	return &Manager{
		reportWriter: reportWriter,
		bus:          bus,
	}
}

func (m *Manager) Submit(
	ctx context.Context,
	targetID xid.ID,
	targetKind datagraph.Kind,
	reason string,
) (*report.Report, error) {
	rep, err := m.reportWriter.Create(
		ctx,
		targetID,
		targetKind,
		opt.NewEmpty[account.AccountID](),
		opt.New(reason),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var targetRef *datagraph.Ref
	if rep.TargetItem != nil {
		targetRef = datagraph.NewRef(rep.TargetItem)
	}

	m.bus.Publish(ctx, &message.EventReportCreated{
		ID:         rep.ID,
		Target:     targetRef,
		ReportedBy: opt.NewEmpty[account.AccountID](),
	})

	return rep, nil
}
