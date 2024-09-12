package mq

import (
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
)

type IndexNode struct {
	ID library.NodeID
}

type IndexPost struct {
	ID post.ID
}

type IndexProfile struct {
	ID account.AccountID
}

type SummariseNode struct {
	ID library.NodeID
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

type Notification struct {
	Event    notification.Event
	Item     *datagraph.Ref
	TargetID account.AccountID
}

type Mention struct {
	Source datagraph.Ref
	Item   datagraph.Ref
}
