package post

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
)

// ID wraps the underlying xid type for all kinds of Storyden Post data type.
type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }

// Author represents a minimal version of the the account that created a post.
type Author struct {
	ID        account.AccountID
	Name      string
	Handle    string
	Admin     bool
	CreatedAt time.Time
}
