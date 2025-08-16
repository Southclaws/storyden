package mq

import (
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
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

type IndexNode struct {
	ID library.NodeID
}

type DeleteNode struct {
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

type DownloadAsset struct {
	URL             string
	ContentFillRule opt.Optional[asset.ContentFillCommand]
}

type AnalyseAsset struct {
	AssetID         xid.ID
	ContentFillRule opt.Optional[asset.ContentFillCommand]
}

type ScrapeLink struct {
	URL  url.URL
	Item *datagraph.Ref
}

type LikePost struct {
	PostID post.ID
}

type UnlikePost struct {
	PostID post.ID
}

type Email struct {
	Message mailer.Message
}

type Notification struct {
	Event    notification.Event
	Item     *datagraph.Ref
	TargetID account.AccountID
	SourceID opt.Optional[account.AccountID]
}

type Mention struct {
	By     account.AccountID
	Source datagraph.Ref
	Item   datagraph.Ref
}

type ReactToPost struct {
	PostID post.ID
}

type CreateEvent struct {
	ID event_ref.EventID
}
