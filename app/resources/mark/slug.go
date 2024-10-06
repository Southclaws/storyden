package mark

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/gosimple/slug"
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
	if !slug.IsSlug(s) {
		return nil, ErrInvalidSlug
	}

	return &Slug{
		slug: s,
	}, nil
}

func NewSlugFromName(s string) Slug {
	return Slug{slug: slug.Make(s)}
}
