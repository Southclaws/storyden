package audit

//go:generate go run github.com/Southclaws/enumerator

type eventTypeEnum string

const (
	eventTypeThreadDeleted         eventTypeEnum = "thread_deleted"
	eventTypeThreadReplyDeleted    eventTypeEnum = "thread_reply_deleted"
	eventTypeAccountSuspended      eventTypeEnum = "account_suspended"
	eventTypeAccountUnsuspended    eventTypeEnum = "account_unsuspended"
	eventTypeAccountContentPurged  eventTypeEnum = "account_content_purged"
	eventTypeModerationNoteCreated eventTypeEnum = "moderation_note_created"
	eventTypeModerationNoteDeleted eventTypeEnum = "moderation_note_deleted"
	eventTypeAccountWarned                   eventTypeEnum = "account_warned"
	eventTypeAccountWarningUpdated           eventTypeEnum = "account_warning_updated"
	eventTypeAccountWarningDeleted           eventTypeEnum = "account_warning_deleted"
	eventTypeAccountPasswordResetTokenIssued eventTypeEnum = "account_password_reset_token_issued"
	eventTypeAccountPasswordResetEmailSent   eventTypeEnum = "account_password_reset_email_sent"
)
