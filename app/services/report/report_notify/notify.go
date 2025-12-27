package report_notify

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/report/report_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Invoke(func(
		ctx context.Context,
		lc fx.Lifecycle,
		bus *pubsub.Bus,
		notifier *notify.Notifier,
		accountQuerier *account_querier.Querier,
		reportQuerier *report_querier.Querier,
	) {
		consumer := func(hctx context.Context) error {
			// Report submitted
			// Notify only members with MANAGE_REPORTS or ADMINISTRATOR perms.
			if _, err := pubsub.Subscribe(ctx, bus, "report_notify.report_created", func(ctx context.Context, evt *message.EventReportCreated) error {
				return sendReportSubmitted(ctx, notifier, accountQuerier, evt)
			}); err != nil {
				return err
			}

			// Report updated
			// Depending on the source of the update:
			// - author: notify handlers (admins/mods)
			// - handler: notify author
			if _, err := pubsub.Subscribe(ctx, bus, "report_notify.report_updated", func(ctx context.Context, evt *message.EventReportUpdated) error {
				return sendReportUpdated(ctx, notifier, accountQuerier, reportQuerier, evt)
			}); err != nil {
				return err
			}

			return nil
		}

		lc.Append(fx.StartHook(consumer))
	})
}

func sendReportSubmitted(
	ctx context.Context,
	notifier *notify.Notifier,
	accountQuerier *account_querier.Querier,
	evt *message.EventReportCreated,
) error {
	accs, err := accountQuerier.ListByHeldPermission(ctx, rbac.PermissionAdministrator, rbac.PermissionManageReports)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	for _, acc := range accs {
		if err := notifier.Send(ctx, acc.ID, evt.ReportedBy, notification.EventReportSubmitted, evt.Target); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func sendReportUpdated(
	ctx context.Context,
	notifier *notify.Notifier,
	accountQuerier *account_querier.Querier,
	reportQuerier *report_querier.Querier,
	evt *message.EventReportUpdated,
) error {
	source, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	rep, err := reportQuerier.Get(ctx, evt.ID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	reportedBy, ok := rep.ReportedBy.Get()
	if !ok {
		return nil
	}

	if reportedBy.ID != source {
		// Moderator updated the report; notify author.

		if err := notifier.Send(
			ctx,
			reportedBy.ID,
			opt.NewEmpty[account.AccountID](),
			notification.EventReportUpdated,
			nil,
		); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		return nil
	}

	// Author updated the report; notify handlers.

	accs, err := accountQuerier.ListByHeldPermission(ctx, rbac.PermissionAdministrator, rbac.PermissionManageReports)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	for _, acc := range accs {
		if err := notifier.Send(ctx, acc.ID, opt.New(source), notification.EventReportUpdated, nil); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}
