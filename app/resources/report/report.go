package report

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID xid.ID

func (i ID) String() string { return xid.ID(i).String() }

func (i ID) MarshalJSON() ([]byte, error) {
	return xid.ID(i).MarshalJSON()
}

func (i *ID) UnmarshalJSON(data []byte) error {
	var id xid.ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	*i = ID(id)
	return nil
}

type Report struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	Status         Status
	TargetItemKind datagraph.Kind
	TargetItemID   xid.ID
	TargetItem     datagraph.Item
	ReportedBy     opt.Optional[account.Account]
	HandledBy      opt.Optional[account.Account]
	Comment        opt.Optional[string]
}

type Reports []*Report

func (a Reports) Len() int           { return len(a) }
func (a Reports) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Reports) Less(i, j int) bool { return a[i].CreatedAt.After(a[j].CreatedAt) }

type ReportRef struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	Status     Status
	TargetRef  datagraph.Ref
	ReportedBy opt.Optional[account.Account]
	HandledBy  opt.Optional[account.Account]
	Comment    opt.Optional[string]
}

type ReportRefs []*ReportRef

func Map(r *ent.Report) (*ReportRef, error) {
	reportedByEdge := opt.NewPtr(r.Edges.ReportedBy)
	reportedBy, err := opt.MapErr(reportedByEdge, func(a ent.Account) (account.Account, error) {
		p, err := account.MapRef(&a)
		if err != nil {
			return account.Account{}, err
		}
		return *p, nil
	})
	if err != nil {
		return nil, err
	}

	handledByEdge := opt.NewPtr(r.Edges.HandledBy)
	handledBy, err := opt.MapErr(handledByEdge, func(a ent.Account) (account.Account, error) {
		p, err := account.MapRef(&a)
		if err != nil {
			return account.Account{}, err
		}
		return *p, nil
	})
	if err != nil {
		return nil, err
	}

	status, err := NewStatus(r.Status)
	if err != nil {
		return nil, err
	}

	targetKind, err := datagraph.NewKind(r.TargetKind)
	if err != nil {
		return nil, err
	}

	return &ReportRef{
		ID:        ID(r.ID),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		TargetRef: datagraph.Ref{
			ID:   r.TargetID,
			Kind: targetKind,
		},
		Status:     status,
		ReportedBy: reportedBy,
		HandledBy:  handledBy,
		Comment:    opt.NewPtr(r.Comment),
	}, nil
}
