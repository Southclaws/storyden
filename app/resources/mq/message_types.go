package mq

import (
	"github.com/Southclaws/storyden/app/resources/account"
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
