package audit_logger

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/audit/audit_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(s *Service) {}),
	)
}

type Service struct {
	writer *audit_writer.Writer
	bus    *pubsub.Bus
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	writer *audit_writer.Writer,
	bus *pubsub.Bus,
) *Service {
	s := &Service{
		writer: writer,
		bus:    bus,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.thread_deleted", s.onThreadDeleted); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.thread_reply_deleted", s.onThreadReplyDeleted); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.account_suspended", s.onAccountSuspended); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.account_unsuspended", s.onAccountUnsuspended); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.moderation_note_created", s.onModerationNoteCreated); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.moderation_note_deleted", s.onModerationNoteDeleted); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.account_warned", s.onAccountWarned); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.account_warning_updated", s.onAccountWarningUpdated); err != nil {
			return err
		}
		if _, err := pubsub.Subscribe(ctx, bus, "audit_logger.account_warning_deleted", s.onAccountWarningDeleted); err != nil {
			return err
		}

		return nil
	}))

	return s
}

func (s *Service) onThreadDeleted(ctx context.Context, event *rpc.EventThreadDeleted) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeThreadDeleted,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.ID),
			Kind: datagraph.KindThread,
		}),
		nil,
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onThreadReplyDeleted(ctx context.Context, event *rpc.EventThreadReplyDeleted) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeThreadReplyDeleted,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.ReplyID),
			Kind: datagraph.KindReply,
		}),
		nil,
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onAccountSuspended(ctx context.Context, event *rpc.EventAccountSuspended) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeAccountSuspended,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.ID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.ID.String(),
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onAccountUnsuspended(ctx context.Context, event *rpc.EventAccountUnsuspended) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeAccountUnsuspended,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.ID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.ID.String(),
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onModerationNoteCreated(ctx context.Context, event *rpc.EventModerationNoteCreated) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeModerationNoteCreated,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.AccountID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.AccountID.String(),
			"note_id":    event.NoteID.String(),
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onModerationNoteDeleted(ctx context.Context, event *rpc.EventModerationNoteDeleted) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeModerationNoteDeleted,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.AccountID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.AccountID.String(),
			"note_id":    event.NoteID.String(),
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onAccountWarned(ctx context.Context, event *rpc.EventAccountWarned) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeAccountWarned,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.AccountID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.AccountID.String(),
			"warning_id": event.WarningID,
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onAccountWarningUpdated(ctx context.Context, event *rpc.EventAccountWarningUpdated) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeAccountWarningUpdated,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.AccountID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id":      event.AccountID.String(),
			"warning_id":      event.WarningID,
			"previous_reason": event.PreviousReason,
			"reason":          event.Reason,
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Service) onAccountWarningDeleted(ctx context.Context, event *rpc.EventAccountWarningDeleted) error {
	enactedBy := session.GetOptAccountID(ctx)

	_, err := s.writer.Create(
		ctx,
		audit.EventTypeAccountWarningDeleted,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(event.AccountID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": event.AccountID.String(),
			"warning_id": event.WarningID,
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
