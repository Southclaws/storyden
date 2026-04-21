package warning_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/warning"
	"github.com/Southclaws/storyden/app/resources/account/warning/warning_repo"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	repo     *warning_repo.Repository
	bus      *pubsub.Bus
	notifier *notify.Notifier
}

func Build() fx.Option { return fx.Provide(New) }

func New(
	lc fx.Lifecycle,
	repo *warning_repo.Repository,
	bus *pubsub.Bus,
	notifier *notify.Notifier,
) *Manager {
	m := &Manager{repo: repo, bus: bus, notifier: notifier}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		_, err := pubsub.Subscribe(ctx, bus, "warning_manager.account_warned", m.onAccountWarned)
		return err
	}))

	return m
}

func (m *Manager) Issue(ctx context.Context, givenBy account.AccountID, receiveBy account.AccountID, reason string) (*warning.Warning, error) {
	record, err := m.repo.Create(ctx, givenBy, receiveBy, reason)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountWarned{
		AccountID: receiveBy,
		AuthorID:  givenBy,
		WarningID: record.ID.String(),
	})

	return record, nil
}

func (m *Manager) ListForAccount(ctx context.Context, accountID account.AccountID) (warning.Warnings, error) {
	records, err := m.repo.ListByAccountID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return records, nil
}

func (m *Manager) UpdateReason(ctx context.Context, updatedBy account.AccountID, accountID account.AccountID, warningID warning.ID, reason string) (*warning.Warning, error) {
	original, err := m.repo.Get(ctx, accountID, warningID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := m.repo.UpdateReason(ctx, accountID, warningID, reason)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if original.Reason != record.Reason {
		m.bus.Publish(ctx, &rpc.EventAccountWarningUpdated{
			AccountID:      accountID,
			AuthorID:       updatedBy,
			WarningID:      record.ID.String(),
			PreviousReason: original.Reason,
			Reason:         record.Reason,
		})
	}

	return record, nil
}

func (m *Manager) Delete(ctx context.Context, accountID account.AccountID, warningID warning.ID, deletedBy account.AccountID) error {
	if err := m.repo.Delete(ctx, accountID, warningID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountWarningDeleted{
		AccountID: accountID,
		AuthorID:  deletedBy,
		WarningID: warningID.String(),
	})

	return nil
}

func (m *Manager) onAccountWarned(ctx context.Context, event *rpc.EventAccountWarned) error {
	item := &datagraph.Ref{
		ID:   xid.ID(event.AccountID),
		Kind: datagraph.KindProfile,
	}

	source := opt.New(event.AuthorID)

	if err := m.notifier.Send(ctx, event.AccountID, source, notification.EventWarningIssued, item); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
