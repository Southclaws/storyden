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

type IndexNode struct {
	ID library.NodeID
}

type DeleteNode struct {
	ID library.NodeID
}

type IndexThread struct {
	ID post.ID
}

type DeleteThread struct {
	ID post.ID
}

type IndexReply struct {
	ID post.ID
}

type IndexProfile struct {
	ID account.AccountID
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
}

type Mention struct {
	Source datagraph.Ref
	Item   datagraph.Ref
}

type ReactToPost struct {
	PostID post.ID
}

type CreateEvent struct {
	ID event_ref.EventID
}
