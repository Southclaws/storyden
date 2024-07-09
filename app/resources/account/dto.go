package account

import (
	"net/url"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

var errSuspended = fault.Wrap(fault.New("suspended"), ftag.With(ftag.PermissionDenied))

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID            AccountID
	Handle        string
	Name          string
	Bio           content.Rich
	Admin         bool
	Auths         []string
	ExternalLinks []ExternalLink
	Metadata      map[string]any

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

func FromModel(a *ent.Account) (*Account, error) {
	auths := dt.Map(a.Edges.Authentication, func(a *ent.Authentication) string {
		return a.Service
	})

	bio, err := content.NewRichText(a.Bio)
	if err != nil {
		return nil, err
	}

	links, err := dt.MapErr(a.Links, MapExternalLink)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Account{
		ID:            AccountID(a.ID),
		Handle:        a.Handle,
		Name:          a.Name,
		Bio:           bio,
		Admin:         a.Admin,
		Auths:         auths,
		ExternalLinks: links,
		Metadata:      a.Metadata,

		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: opt.NewPtr(a.DeletedAt),
	}, nil
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
