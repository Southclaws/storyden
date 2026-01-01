package audit

//go:generate go run github.com/Southclaws/enumerator

type eventTypeEnum string

const (
	eventTypeThreadDeleted        eventTypeEnum = "thread_deleted"
	eventTypeThreadReplyDeleted   eventTypeEnum = "thread_reply_deleted"
	eventTypeAccountSuspended     eventTypeEnum = "account_suspended"
	eventTypeAccountUnsuspended   eventTypeEnum = "account_unsuspended"
	eventTypeAccountContentPurged eventTypeEnum = "account_content_purged"
)
