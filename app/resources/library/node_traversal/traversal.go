package node_traversal

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type Repository interface {
	Root(ctx context.Context, opts ...Filter) ([]*library.Node, error)
	Subtree(ctx context.Context, id opt.Optional[library.NodeID], flatten bool, opts ...Filter) ([]*library.Node, error)
}

type filters struct {
	rootAccountHandleFilter *string
	// NOTE: This should really just be the roles list, not a full account obj.
	requestingAccount opt.Optional[account.AccountWithEdges]
	visibility        []visibility.Visibility
	depth             *uint
}

type Filter func(*filters)

// WithRootOwner filters top level nodes only by the account handle. When this
// is used, it will only retrieve fully private trees from the top level. If the
// specified handle either only owns nodes that are children or has no nodes at
// all, the result will be an empty list. Generally this is only used for use
// cases where you need to show users their own private root level page trees.
func WithRootOwner(v string) Filter {
	return func(f *filters) {
		f.rootAccountHandleFilter = &v
	}
}

// WithVisibility applies permission-based filtering for the given visibilities
// against the requesting account (if any) to ensure that visibility rules are
// implemented correctly. Owners can view their own drafts, library managers can
// view all review items, etc.
func WithVisibility(acc opt.Optional[account.AccountWithEdges], v ...visibility.Visibility) Filter {
	return func(f *filters) {
		f.requestingAccount = acc
		f.visibility = v
	}
}

func WithDepth(v uint) Filter {
	return func(f *filters) {
		f.depth = &v
	}
}
