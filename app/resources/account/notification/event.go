package notification

//go:generate go run github.com/Southclaws/enumerator

type eventEnum string

const (
	eventThreadReply    eventEnum = "thread_reply"
	eventPostLike       eventEnum = "post_like"
	eventFollow         eventEnum = "follow"
	eventProfileMention eventEnum = "profile_mention"
)
