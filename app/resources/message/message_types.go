package message

import (
	"net/url"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/report"
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
	ThreadID       post.ID
	ReplyID        post.ID
	ThreadAuthorID account.AccountID
	ReplyAuthorID  account.AccountID
}

type EventThreadReplyDeleted struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventThreadReplyUpdated struct {
	ThreadID post.ID
	ReplyID  post.ID
}

type EventPostLiked struct {
	PostID post.ID
}

type EventPostUnliked struct {
	PostID post.ID
}

type EventPostReacted struct {
	PostID post.ID
}

type EventMemberMentioned struct {
	By     account.AccountID
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
	ID library.NodeID
}

type EventNodeUpdated struct {
	ID library.NodeID
}

type EventNodeDeleted struct {
	ID library.NodeID
}

type EventNodePublished struct {
	ID library.NodeID
}

type EventNodeSubmittedForReview struct {
	ID library.NodeID
}

type EventNodeUnpublished struct {
	ID library.NodeID
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
	ID account.AccountID
}

type EventAccountUpdated struct {
	ID account.AccountID
}

type CommandProfileIndex struct {
	ID account.AccountID
}

// -
// Notifications
// -

type CommandSendNotification struct {
	Event    notification.Event
	Item     *datagraph.Ref
	TargetID account.AccountID
	SourceID opt.Optional[account.AccountID]
}

type CommandSendEmail struct {
	Message mailer.Message
}

type CommandSendBeacon struct {
	Item    datagraph.Ref
	Subject opt.Optional[account.AccountID]
}

// -
// Reports
// -

type EventReportCreated struct {
	ID         report.ID
	Target     *datagraph.Ref
	ReportedBy account.AccountID
}

type EventReportUpdated struct {
	ID         report.ID
	Target     *datagraph.Ref
	ReportedBy account.AccountID
	HandledBy  opt.Optional[account.AccountID]
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
