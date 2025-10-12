// Mark represents a flexible identifier which may be formed as a raw XID, an
// XID and a slug separated by a hyphen, or just a slug on its own. This type
// allows API consumers to use any form easily without needing to do queries.
//
// Terminology:
//
// - "cpvf89ifunp0qr2aqkog": an ID, in xid format.
// - "some-super-neat-post": a slug, a human-readable identifier.
// - "cpvf89ifunp0qr2aqkog-some-super-neat-post": a mark, an ID and a slug.
//
// Conventions:
//
// - The mark itself is never stored in the database.
// - ID and slug must always be stored in the database.
// - Resources must hold a Hydrated mark in their struct.
// - When a resource are read from the database, its mark is constructed.
// - A key is used on the inbound request path when uncertain about input.
// - Keys are always used to query, except when only an ID is available.
package mark

import (
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

// xidEncodedLength is the length of an xid encoded as a string
const xidEncodedLength = 20

var (
	ErrMissingID   = fault.New("read key does not contain an ID")
	ErrMissingSlug = fault.New("read key does not contain a slug")
)

// Mark represents a hydrated mark key where the ID and Slug are both present.
// This is used on the read-return-path where the data has been read from the
// store to the fact that once a resource has been stored, it will always have
// an ID and a slug assigned to it.
type Mark struct {
	id   xid.ID // the resource's unique identifier
	slug string // the resource's non-unique URL slug
	mark string // the resource's unique identifier and slug joined with '-'
}

func NewMark(id xid.ID, slug string) *Mark {
	mark := fmt.Sprintf("%s-%s", id.String(), slug)

	return &Mark{
		mark: mark,
		id:   id,
		slug: slug,
	}
}

func (m Mark) String() string {
	return m.mark
}

func (m Mark) ID() xid.ID {
	return m.id
}

func (m Mark) Slug() string {
	return m.slug
}

func (m Mark) Queryable() Queryable {
	return Queryable{
		raw:  m.mark,
		id:   opt.New(m.id),
		slug: opt.New(m.slug),
	}
}

// Queryable is used on the read-request-path where the input ID is of unknown
// form and may contain an ID, a slug, or both. This is used to parse the input
// and provide a query predicate going into the resource repository layer.
//
// This type must not be stored on actual resource structs. Use Queryable for that.
type Queryable struct {
	raw string

	// id stores the underlying, parsed ID, if there is one in the input string.
	id   opt.Optional[xid.ID]
	slug opt.Optional[string]
}

func (m Queryable) String() string {
	return m.raw
}

func (m Queryable) ID() opt.Optional[xid.ID] {
	return m.id
}

func (m Queryable) Slug() opt.Optional[string] {
	return m.slug
}

func (m Queryable) Equal(b Queryable) bool {
	if a, ok := m.id.Get(); ok {
		if b, ok := b.id.Get(); ok {
			if a != b {
				return false
			}
		}
	}

	if a, ok := m.slug.Get(); ok {
		if b, ok := b.slug.Get(); ok {
			if a != b {
				return false
			}
		}
	}

	return true
}

func (m Queryable) Mark() (*Mark, error) {
	id, hasID := m.id.Get()
	if !hasID {
		return nil, ErrMissingID
	}

	slug, hasSlug := m.slug.Get()
	if !hasSlug {
		return nil, ErrMissingSlug
	}

	return &Mark{
		mark: m.raw,
		id:   id,
		slug: slug,
	}, nil
}

func NewQueryKeyID(id xid.ID) Queryable {
	return Queryable{
		raw: id.String(),
		id:  opt.New(id),
	}
}

func NewQueryKey(raw string) Queryable {
	runeCount := len([]rune(raw))
	exactLength := runeCount == xidEncodedLength
	exactLengthOrLonger := runeCount >= xidEncodedLength
	containsHyphen := runeCount > xidEncodedLength && len(raw) > xidEncodedLength && raw[xidEncodedLength] == '-'

	probablyXID := exactLength || (exactLengthOrLonger && containsHyphen)

	// Form: <xid>
	if probablyXID {
		v, err := xid.FromString(raw[:xidEncodedLength])
		if err != nil {
			// In the unlikely (but not impossible) case that the input string
			// is exactly 20 characters but not a valid XID, treat it as a slug.
			return Queryable{
				raw:  raw,
				slug: opt.New(raw),
			}
		}

		// The input is a raw ID, no trailing slug.
		m := Queryable{
			raw: raw,
			id:  opt.New(v),
		}

		// Form: <xid>-<slug>
		if runeCount > xidEncodedLength+1 && len(raw) > xidEncodedLength && raw[xidEncodedLength] == '-' {
			if len(raw) > xidEncodedLength+1 {
				m.slug = opt.New(raw[xidEncodedLength+1:])
			}
		}

		return m
	}

	// Form: <slug>
	return Queryable{
		raw:  raw,
		slug: opt.New(raw),
	}
}

// Apply calls the appropriate function based on the Mark's contents. The Mark's
// ID takes precendence as it is the most reliable identifier.
func (m Queryable) Apply(id func(xid.ID), slug func(string)) {
	if v, ok := m.id.Get(); ok {
		id(v)
		return
	}

	if v, ok := m.slug.Get(); ok {
		slug(v)
		return
	}

	panic("mark does not contain an ID or a slug")
}

// ApplyAll calls either/both functions based on the Mark's contents. This is
// generally used during create operations where both fields may be set.
func (m Queryable) ApplyAll(id func(xid.ID), slug func(string)) {
	if v, ok := m.id.Get(); ok {
		id(v)
	}

	if v, ok := m.slug.Get(); ok {
		slug(v)
	}
}
