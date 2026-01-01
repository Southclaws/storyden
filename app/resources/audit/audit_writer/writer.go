package audit_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (w *Writer) Create(
	ctx context.Context,
	eventType audit.EventType,
	enactedBy opt.Optional[account.AccountID],
	target opt.Optional[datagraph.Ref],
	metadata map[string]any,
) (*audit.AuditLog, error) {
	create := w.db.AuditLog.Create()

	create.SetType(eventType.String())

	enactedBy.Call(func(id account.AccountID) {
		create.SetEnactedByID(xid.ID(id))
	})

	target.Call(func(ref datagraph.Ref) {
		create.SetTargetID(xid.ID(ref.ID))
		create.SetTargetKind(ref.Kind.String())
	})

	if metadata != nil {
		create.SetMetadata(metadata)
	}

	al, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return audit.Map(al)
}

func (w *Writer) RecordFailure(
	ctx context.Context,
	eventType audit.EventType,
	enactedBy opt.Optional[account.AccountID],
	target opt.Optional[datagraph.Ref],
	metadata map[string]any,
	failure error,
) (*audit.AuditLog, error) {
	create := w.db.AuditLog.Create()

	create.SetType(eventType.String())
	create.SetError(failure.Error())

	enactedBy.Call(func(id account.AccountID) {
		create.SetEnactedByID(xid.ID(id))
	})

	target.Call(func(ref datagraph.Ref) {
		create.SetTargetID(xid.ID(ref.ID))
		create.SetTargetKind(ref.Kind.String())
	})

	if metadata != nil {
		create.SetMetadata(metadata)
	}

	al, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return audit.Map(al)
}
