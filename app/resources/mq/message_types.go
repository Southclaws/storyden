package mq

import (
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
)

type IndexNode struct {
	ID datagraph.NodeID
}

type IndexPost struct {
	ID post.ID
}

type IndexProfile struct {
	ID account.AccountID
}
