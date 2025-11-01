package event_ref

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/mark"
	ent_event "github.com/Southclaws/storyden/internal/ent/event"
	"github.com/Southclaws/storyden/internal/ent/eventparticipant"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type (
	Mark     struct{ mark.Mark }
	QueryKey struct{ mark.Queryable }
)

func NewMark(id xid.ID, slug string) Mark { return Mark{*mark.NewMark(id, slug)} }
func NewKey(in string) QueryKey           { return QueryKey{mark.NewQueryKey(in)} }
func NewID(id xid.ID) QueryKey            { return QueryKey{mark.NewQueryKeyID(id)} }
func NewQueryKey(m Mark) QueryKey         { return QueryKey{m.Queryable()} }

func (m QueryKey) Predicate() (p predicate.Event) {
	m.Apply(
		func(i xid.ID) { p = ent_event.ID(i) },
		func(s string) { p = ent_event.Slug(s) })
	return
}

func (m QueryKey) ParticipantPredicate() (p predicate.EventParticipant) {
	m.Apply(
		func(i xid.ID) { p = eventparticipant.HasEventWith(ent_event.ID(i)) },
		func(s string) { p = eventparticipant.HasEventWith(ent_event.Slug(s)) })
	return
}
