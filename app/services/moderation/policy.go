package moderation

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/app/services/report/system_report"
)

type Manager struct {
	registry       *checker.Registry
	systemReporter *system_report.Manager
}

func New(
	registry *checker.Registry,
	systemReporter *system_report.Manager,
) *Manager {
	return &Manager{
		registry:       registry,
		systemReporter: systemReporter,
	}
}

func (m *Manager) CheckContent(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*checker.Result, error) {
	checkers := m.registry.GetEnabled()

	for _, c := range checkers {
		result, err := c.Check(ctx, targetID, targetKind, name, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if result.RequiresReview {
			_, err := m.systemReporter.Submit(ctx, targetID, targetKind, result.Reason)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			return result, nil
		}
	}

	return &checker.Result{
		RequiresReview: false,
		Reason:         "",
	}, nil
}
