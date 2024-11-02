package collection

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type (
	Mark     struct{ mark.Mark }
	QueryKey struct{ mark.Queryable }
)

func NewMark(id xid.ID, slug string) Mark { return Mark{*mark.NewMark(id, slug)} }
func NewKey(in string) QueryKey           { return QueryKey{mark.NewQueryKey(in)} }
func NewID(in xid.ID) QueryKey            { return QueryKey{mark.NewQueryKeyID(in)} }

func (m QueryKey) Predicate() (p predicate.Collection) {
	m.Apply(
		func(i xid.ID) { p = collection.ID(i) },
		func(s string) { p = collection.Slug(s) })
	return
}
