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
)

var errSuspended = fault.Wrap(fault.New("suspended"), ftag.With(ftag.PermissionDenied))

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID             AccountID
	Handle         string
	Name           string
	Bio            datagraph.Content
	Kind           AccountKind
	Admin          bool
	Followers      int
	Following      int
	LikeScore      int
	Roles          held.Roles
	Auths          []string
	EmailAddresses []*EmailAddress
	VerifiedStatus VerifiedStatus
	ExternalLinks  []ExternalLink
	Metadata       map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
	IndexedAt opt.Optional[time.Time]

	InvitedByID *xid.ID
	InvitedBy   opt.Optional[Account]
}

type Accounts []*Account

func (a Accounts) Map() Lookup {
	return lo.KeyBy(a, func(a *Account) xid.ID { return xid.ID(a.ID) })
}

type Lookup map[xid.ID]*Account

type ExternalLink struct {
	Text string
	URL  url.URL
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
