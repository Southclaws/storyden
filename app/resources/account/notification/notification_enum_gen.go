// Code generated by enumerator. DO NOT EDIT.

package notification

import (
	"database/sql/driver"
	"fmt"
)

type Event struct {
	v eventEnum
}

var (
	EventThreadReply          = Event{eventThreadReply}
	EventPostLike             = Event{eventPostLike}
	EventFollow               = Event{eventFollow}
	EventProfileMention       = Event{eventProfileMention}
	EventEventHostAdded       = Event{eventEventHostAdded}
	EventMemberAttendingEvent = Event{eventMemberAttendingEvent}
	EventMemberDeclinedEvent  = Event{eventMemberDeclinedEvent}
	EventAttendeeRemoved      = Event{eventAttendeeRemoved}
)

func (r Event) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprint(f, r.v)
	case 'q':
		fmt.Fprintf(f, "%q", r.String())
	default:
		fmt.Fprint(f, r.v)
	}
}
func (r Event) String() string {
	return string(r.v)
}
func (r Event) MarshalText() ([]byte, error) {
	return []byte(r.v), nil
}
func (r *Event) UnmarshalText(__iNpUt__ []byte) error {
	s, err := NewEvent(string(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func (r Event) Value() (driver.Value, error) {
	return r.v, nil
}
func (r *Event) Scan(__iNpUt__ any) error {
	s, err := NewEvent(fmt.Sprint(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func NewEvent(__iNpUt__ string) (Event, error) {
	switch __iNpUt__ {
	case string(eventThreadReply):
		return EventThreadReply, nil
	case string(eventPostLike):
		return EventPostLike, nil
	case string(eventFollow):
		return EventFollow, nil
	case string(eventProfileMention):
		return EventProfileMention, nil
	case string(eventEventHostAdded):
		return EventEventHostAdded, nil
	case string(eventMemberAttendingEvent):
		return EventMemberAttendingEvent, nil
	case string(eventMemberDeclinedEvent):
		return EventMemberDeclinedEvent, nil
	case string(eventAttendeeRemoved):
		return EventAttendeeRemoved, nil
	default:
		return Event{}, fmt.Errorf("invalid value for type 'Event': '%s'", __iNpUt__)
	}
}
