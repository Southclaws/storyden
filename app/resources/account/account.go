package account

import (
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

var errSuspended = fault.Wrap(fault.New("suspended"), ftag.With(ftag.PermissionDenied))

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

func (u AccountID) MarshalJSON() ([]byte, error) {
	return xid.ID(u).MarshalJSON()
}

func (u *AccountID) UnmarshalJSON(data []byte) error {
	var id xid.ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	*u = AccountID(id)
	return nil
}

type Account struct {
	ID        AccountID
	CreatedAt time.Time
	UpdatedAt time.Time

	Handle   string
	Name     string
	Bio      datagraph.Content
	Kind     AccountKind
	Admin    bool
	Metadata map[string]any

	DeletedAt opt.Optional[time.Time]
	IndexedAt opt.Optional[time.Time]
}

type AccountWithEdges struct {
	Account
	Roles          held.Roles
	Auths          []string
	EmailAddresses []*EmailAddress
	VerifiedStatus VerifiedStatus
	InvitedBy      opt.Optional[Account]
	ExternalLinks  []ExternalLink
}

type Accounts []*Account

type Lookup map[xid.ID]*ent.Account

func NewAccountLookup(in []*ent.Account) Lookup {
	return lo.KeyBy(in, func(a *ent.Account) xid.ID {
		return a.ID
	})
}

type ExternalLink struct {
	Text string
	URL  url.URL
}

func MapExternalLink(e schema.ExternalLink) (ExternalLink, error) {
	u, err := url.Parse(e.URL)
	if err != nil {
		return ExternalLink{}, err
	}

	return ExternalLink{
		Text: e.Text,
		URL:  *u,
	}, nil
}

func (a *Account) IsSuspended() bool {
	return a.DeletedAt.Ok()
}

func (a *Account) RejectSuspended() error {
	if a.IsSuspended() {
		return fault.Wrap(errSuspended, ftag.With(ftag.PermissionDenied))
	}

	return nil
}
