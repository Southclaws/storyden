package audit

import (
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type AuditLogID xid.ID

func (a AuditLogID) String() string { return xid.ID(a).String() }

type AuditLog struct {
	ID        AuditLogID
	CreatedAt time.Time
	EnactedBy opt.Optional[account.Account]
	Target    opt.Optional[datagraph.Ref]
	Type      EventType
	Metadata  map[string]any
}

func Map(in *ent.AuditLog) (*AuditLog, error) {
	eventType, err := NewEventType(in.Type)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	var enactedBy opt.Optional[account.Account]
	if eb := in.Edges.EnactedBy; eb != nil {
		mapped, err := account.MapRef(eb)
		if err != nil {
			return nil, err
		}
		enactedBy = opt.New(*mapped)
	}

	var target opt.Optional[datagraph.Ref]
	if tid := in.TargetID; tid != nil && in.TargetKind != nil {
		kind, err := datagraph.NewKind(*in.TargetKind)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		target = opt.New(datagraph.Ref{
			ID:   xid.ID(*tid),
			Kind: kind,
		})
	}

	return &AuditLog{
		ID:        AuditLogID(in.ID),
		CreatedAt: in.CreatedAt,
		EnactedBy: enactedBy,
		Target:    target,
		Type:      eventType,
		Metadata:  in.Metadata,
	}, nil
}
