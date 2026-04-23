package message

import (
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

// -
// Indexing commands
// -

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
// Library node commands
// -

type CommandNodeIndex struct {
	ID library.NodeID
}

type CommandNodeDeindex struct {
	ID library.NodeID
}

// -
// Account and profile commands
// -

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
	ID      xid.ID
	Message mailer.Message
}

type CommandAttemptQueuedEmail struct {
	ID xid.ID
}

type CommandSendBeacon struct {
	Item    datagraph.Ref
	Subject opt.Optional[account_ref.ID]
}

// -
// Scraping commands
// -

type CommandScrapeLink struct {
	URL  url.URL
	Item *datagraph.Ref
}
