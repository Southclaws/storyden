package library

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type (
	Mark     struct{ mark.Mark }
	QueryKey struct{ mark.Queryable }
)

func NewMark(id xid.ID, slug string) Mark { return Mark{*mark.NewMark(id, slug)} }
func NewKey(in string) QueryKey           { return QueryKey{mark.NewQueryKey(in)} }
func NewID(id xid.ID) QueryKey            { return QueryKey{mark.NewQueryKeyID(id)} }
func NewQueryKey(m Mark) QueryKey { return QueryKey{m.Queryable()} }

func (m QueryKey) Predicate() (p predicate.Node) {
	m.Apply(
		func(i xid.ID) { p = node.ID(i) },
		func(s string) { p = node.Slug(s) })
	return
}
