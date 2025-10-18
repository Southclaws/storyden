package notification

//go:generate go run github.com/Southclaws/enumerator

type eventEnum string

// TODO: Maybe rename these, "Event" is duplicated on events management.
const (
	eventThreadReply          eventEnum = "thread_reply"
	eventPostLike             eventEnum = "post_like"
	eventFollow               eventEnum = "follow"
	eventProfileMention       eventEnum = "profile_mention"
	eventEventHostAdded       eventEnum = `event_host_added`
	eventMemberAttendingEvent eventEnum = `member_attending_event`
	eventMemberDeclinedEvent  eventEnum = `member_declined_event`
	eventAttendeeRemoved      eventEnum = `attendee_removed`
	eventReportSubmitted      eventEnum = "report_submitted"
	eventReportUpdated        eventEnum = "report_updated"
)
