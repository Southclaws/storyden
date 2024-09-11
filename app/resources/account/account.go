package account

import (
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/rs/xid"
)

var errSuspended = fault.Wrap(fault.New("suspended"), ftag.With(ftag.PermissionDenied))

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID             AccountID
	Handle         string
	Name           string
	Bio            datagraph.Content
	Admin          bool
	Followers      int
	Following      int
	LikeScore      int
	Auths          []string
	EmailAddresses []*EmailAddress
	VerifiedStatus VerifiedStatus
	ExternalLinks  []ExternalLink
	Metadata       map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

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

// Name is the role/resource name.
const Name = "Account"

func (a *Account) GetRole() string {
	if a.Admin {
		return "everyone"
	}

	return "owner"
}

func (*Account) GetResourceName() string { return Name }
