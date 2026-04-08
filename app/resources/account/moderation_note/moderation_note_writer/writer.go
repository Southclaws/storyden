package moderation_note_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/moderation_note"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/internal/ent"
	entmoderationnote "github.com/Southclaws/storyden/internal/ent/moderationnote"
)

type Writer struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *Writer {
	return &Writer{db: db, roleQuerier: roleQuerier}
}

func (w *Writer) Create(ctx context.Context, accountID xid.ID, authorID xid.ID, content string) (*moderation_note.Note, error) {
	note, err := w.db.ModerationNote.Create().
		SetAccountID(accountID).
		SetAuthorID(authorID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	stored, err := w.db.ModerationNote.Query().
		Where(entmoderationnote.IDEQ(note.ID)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if stored.Edges.Author != nil {
		if err := w.roleQuerier.HydrateRoleEdges(ctx, stored.Edges.Author); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return moderation_note.Map(stored)
}

func (w *Writer) Delete(ctx context.Context, accountID xid.ID, noteID moderation_note.ID) error {
	err := w.db.ModerationNote.
		DeleteOneID(noteID).
		Where(entmoderationnote.AccountIDEQ(accountID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
