package mq

import (
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/rs/xid"
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
