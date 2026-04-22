package account_moderation_note

import (
	"context"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/moderation_note"
	"github.com/Southclaws/storyden/app/resources/account/moderation_note/moderation_note_querier"
	"github.com/Southclaws/storyden/app/resources/account/moderation_note/moderation_note_writer"
	"github.com/Southclaws/storyden/app/services/account/account_manage"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	accountManage *account_manage.Manager
	querier       *moderation_note_querier.Querier
	writer        *moderation_note_writer.Writer
	bus           *pubsub.Bus
}

func New(
	accountManage *account_manage.Manager,
	querier *moderation_note_querier.Querier,
	writer *moderation_note_writer.Writer,
	bus *pubsub.Bus,
) *Manager {
	return &Manager{
		accountManage: accountManage,
		querier:       querier,
		writer:        writer,
		bus:           bus,
	}
}

func (m *Manager) ListByAccountID(ctx context.Context, targetID account.AccountID) (moderation_note.Notes, error) {
	if _, err := m.accountManage.GetByID(ctx, targetID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	notes, err := m.querier.ListByAccountID(ctx, xid.ID(targetID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return notes, nil
}

func (m *Manager) Create(ctx context.Context, targetID account.AccountID, content string) (*moderation_note.Note, error) {
	if _, err := m.accountManage.GetByID(ctx, targetID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, fault.Wrap(
			fault.New("moderation note content required", ftag.With(ftag.InvalidArgument)),
			fctx.With(ctx),
			fmsg.WithDesc("content", "Moderation note content cannot be empty."),
		)
	}

	authorID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	note, err := m.writer.Create(ctx, xid.ID(targetID), xid.ID(authorID), trimmed)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventModerationNoteCreated{
		AccountID: targetID,
		NoteID:    note.ID,
	})

	return note, nil
}

func (m *Manager) Delete(ctx context.Context, targetID account.AccountID, noteID moderation_note.ID) error {
	if _, err := m.accountManage.GetByID(ctx, targetID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.writer.Delete(ctx, xid.ID(targetID), noteID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventModerationNoteDeleted{
		AccountID: targetID,
		NoteID:    noteID,
	})

	return nil
}
