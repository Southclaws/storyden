package moderation

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/app/services/report/system_report"
)

var errModerationPolicy = fault.New("moderation policy")

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

		switch result.Action {
		case checker.ActionReject:
			return nil, fault.Wrap(
				errModerationPolicy,
				fctx.With(ctx),
				ftag.With(ftag.InvalidArgument),
				fmsg.WithDesc("rejected", result.Reason),
			)

		case checker.ActionReport:
			_, err := m.systemReporter.Submit(ctx, targetID, targetKind, result.Reason)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			return result, nil

		case checker.ActionAllow:
			continue
		}
	}

	return &checker.Result{
		Action: checker.ActionAllow,
	}, nil
}
