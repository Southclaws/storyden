package message

import (
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

// -
// Thread events and commands
// -

type EventThreadPublished struct {
	ID post.ID
}

type EventThreadUnpublished struct {
	ID post.ID
}

type EventThreadUpdated struct {
	ID post.ID
}

type EventThreadDeleted struct {
	ID post.ID
}

type EventThreadReplyCreated struct {
	ThreadID        post.ID
	ReplyID         post.ID
	ThreadAuthorID  account_ref.ID
	ReplyAuthorID   account_ref.ID
	ReplyToAuthorID opt.Optional[account_ref.ID]
	ReplyToTargetID opt.Optional[post.ID]
}

type EventThreadReplyDeleted struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventThreadReplyUpdated struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventThreadReplyPublished struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventThreadReplyUnpublished struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventPostLiked struct {
	PostID     post.ID
	RootPostID post.ID
}

type EventPostUnliked struct {
	PostID     post.ID
	RootPostID post.ID
}

type EventPostReacted struct {
	PostID     post.ID
	RootPostID post.ID
}

type EventPostUnreacted struct {
	PostID     post.ID
	RootPostID post.ID
}

// -
// Category events and commands
// -

type EventCategoryUpdated struct {
	ID   xid.ID
	Slug string
}

type EventCategoryDeleted struct {
	ID   xid.ID
	Slug string
}

type EventMemberMentioned struct {
	By     account_ref.ID
	Source datagraph.Ref
	Item   datagraph.Ref
}

type CommandThreadIndex struct {
	ID post.ID
}

type CommandThreadDeindex struct {
	ID post.ID
}

type CommandReplyIndex struct {
	ID post.ID
}

type CommandReplyDeindex struct {
	ID post.ID
}

// -
// Library node events and commands
// -

type EventNodeCreated struct {
	ID   library.NodeID
	Slug string
}

type EventNodeUpdated struct {
	ID   library.NodeID
	Slug string
}

type EventNodeDeleted struct {
	ID   library.NodeID
	Slug string
}

type EventNodePublished struct {
	ID   library.NodeID
	Slug string
}

type EventNodeSubmittedForReview struct {
	ID   library.NodeID
	Slug string
}

type EventNodeUnpublished struct {
	ID   library.NodeID
	Slug string
}

type CommandNodeIndex struct {
	ID library.NodeID
}

type CommandNodeDeindex struct {
	ID library.NodeID
}

// -
// Account and profile events and commands
// -

type EventAccountCreated struct {
	ID account_ref.ID
}

type EventAccountUpdated struct {
	ID account_ref.ID
}

type EventAccountSuspended struct {
	ID account_ref.ID
}

type EventAccountUnsuspended struct {
	ID account_ref.ID
}

type CommandProfileIndex struct {
	ID account_ref.ID
}

// -
// Notifications
// -

type CommandSendNotification struct {
	Event    notification.Event
	Item     *datagraph.Ref
	TargetID account_ref.ID
	SourceID opt.Optional[account_ref.ID]
}

type CommandSendEmail struct {
	Message mailer.Message
}

type CommandSendBeacon struct {
	Item    datagraph.Ref
	Subject opt.Optional[account_ref.ID]
}

// -
// Reports
// -

type EventReportCreated struct {
	ID         report.ID
	Target     *datagraph.Ref
	ReportedBy opt.Optional[account_ref.ID]
}

type EventReportUpdated struct {
	ID         report.ID
	Target     *datagraph.Ref
	ReportedBy opt.Optional[account_ref.ID]
	HandledBy  opt.Optional[account_ref.ID]
	Status     report.Status
}

// -
// Scraping commands
// -

type CommandScrapeLink struct {
	URL  url.URL
	Item *datagraph.Ref
}

// -
// Scheduled event events
// -

type EventActivityCreated struct {
	ID event_ref.EventID
}

type EventActivityUpdated struct {
	ID event_ref.EventID
}

type EventActivityDeleted struct {
	ID event_ref.EventID
}

type EventActivityPublished struct {
	ID event_ref.EventID
}

// -
// Settings events
// -

type EventSettingsUpdated struct {
	Settings *settings.Settings
}
