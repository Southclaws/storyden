package mark

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
)

var ErrInvalidSlug = fault.New("slug is not formed well", ftag.With(ftag.InvalidArgument))

// Slug is a type that represents a URL slug. Not the same as a Mark, which is a
// more flexible identifier a slug is simply the URL-friendly version of a name.
type Slug struct {
	slug string
}

func (s Slug) String() string {
	return s.slug
}

func NewSlug(s string) (*Slug, error) {
	if !IsSlug(s) {
		return nil, fault.Wrap(ErrInvalidSlug, fmsg.WithDesc("invalid slug", "The specified slug is not valid, it must be a URL-friendly string without spaces."))
	}

	return &Slug{
		slug: s,
	}, nil
}

func NewSlugFromName(s string) Slug {
	return Slug{slug: Slugify(s)}
}
