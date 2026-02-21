package nodetree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_cache"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var (
	ErrIdenticalParentChild = fault.New("cannot relate a node to itself", ftag.With(ftag.InvalidArgument))
	ErrVisibilityRules      = fault.New("requested relationship violates visibility rules", ftag.With(ftag.InvalidArgument))
)

type Graph interface {
	// Move moves a node from either orphan state or belonging to one node
	// to another node essentially setting its parent slug to some/new value.
	Move(ctx context.Context, child library.QueryKey, parent library.QueryKey) (*library.Node, error)

	// Sever orphans a node by removing it from its parent to the root level.
	Sever(ctx context.Context, child library.QueryKey, parent library.QueryKey) (*library.Node, error)
}

type service struct {
	nodeQuerier  *node_querier.Querier
	nodeWriter   *node_writer.Writer
	accountQuery *account_querier.Querier
	cache        *node_cache.Cache
	bus          *pubsub.Bus
}

func New(
	nodeChildren *node_children.Writer,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	accountQuery *account_querier.Querier,
	cache *node_cache.Cache,
	bus *pubsub.Bus,
) (Graph, *Position) {
	g := &service{
		nodeQuerier:  nodeQuerier,
		nodeWriter:   nodeWriter,
		accountQuery: accountQuery,
		cache:        cache,
		bus:          bus,
	}

	p := NewPositionService(nodeChildren, nodeQuerier, nodeWriter, g, accountQuery, cache, bus)

	return g, p
}

func (s *service) Move(ctx context.Context, child library.QueryKey, parent library.QueryKey) (*library.Node, error) {
	if child.Equal(parent.Queryable) {
		return nil, fault.Wrap(ErrIdenticalParentChild, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cnode, err := s.nodeQuerier.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pnode, err := s.nodeQuerier.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeParentChildMutation(ctx, acc, cnode, pnode); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	passesVisibilityRules := visibilityRules[pnode.Visibility][cnode.Visibility]

	if !passesVisibilityRules {
		return nil, fault.Wrap(ErrVisibilityRules, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, cnode.GetSlug()); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// If the target parent is actually a child of the target child, sever this
	// connection before adding the target child to the target parent.
	if parentParent, ok := pnode.Parent.Get(); ok {
		if parentParent.Mark.ID() == cnode.Mark.ID() {
			cnode, err = s.nodeWriter.Update(ctx, library.NewQueryKey(cnode.Mark), node_writer.WithChildNodeRemove(xid.ID(pnode.Mark.ID())))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	cnode, err = s.nodeWriter.Update(ctx, library.NewQueryKey(cnode.Mark), node_writer.WithParent(library.NodeID(pnode.Mark.ID())))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeUpdated{
		ID:   library.NodeID(cnode.Mark.ID()),
		Slug: cnode.GetSlug(),
	})

	return cnode, nil
}

func (s *service) Sever(ctx context.Context, child library.QueryKey, parent library.QueryKey) (*library.Node, error) {
	if child.Equal(parent.Queryable) {
		return nil, fault.Wrap(ErrIdenticalParentChild, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cnode, err := s.nodeQuerier.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pnode, err := s.nodeQuerier.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeParentChildMutation(ctx, acc, cnode, pnode); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, cnode.GetSlug()); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = s.nodeWriter.Update(ctx, library.NewQueryKey(pnode.Mark), node_writer.WithChildNodeRemove(xid.ID(cnode.Mark.ID())))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := s.nodeQuerier.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeUpdated{
		ID:   library.NodeID(result.Mark.ID()),
		Slug: result.GetSlug(),
	})

	return result, nil
}

// visibilityRules defines the rules for which visibility levels can be nested.
//
//	--------------------- PARENT ------------------ CHILD ---------------------
var visibilityRules = map[visibility.Visibility]map[visibility.Visibility]bool{
	visibility.VisibilityDraft: {
		visibility.VisibilityDraft:     true,  // draft nodes can only ever contain other draft nodes.
		visibility.VisibilityUnlisted:  false, //
		visibility.VisibilityReview:    false, //
		visibility.VisibilityPublished: false, //
	},
	visibility.VisibilityUnlisted: {
		visibility.VisibilityDraft:     false, // unlisted nodes can only ever contain other unlisted nodes.
		visibility.VisibilityUnlisted:  true,  //
		visibility.VisibilityReview:    false, //
		visibility.VisibilityPublished: false, //
	},
	visibility.VisibilityReview: {
		visibility.VisibilityDraft:     true,  // a submission may contain children, the author may be submitting an entire tree of information and the admin can approve the whole subtree at once.
		visibility.VisibilityUnlisted:  false, // review nodes cannot contain unlisted nodes, for the same reason as published below.
		visibility.VisibilityReview:    true,  // review nodes can contain other review nodes, such as the above review+draft example above.
		visibility.VisibilityPublished: false, // review nodes cannot contain published nodes, it should be impossible to get into this state but if it happens, the parent library being "review" state would prevent any child nodes from being viewed anyway.
	},
	visibility.VisibilityPublished: {
		visibility.VisibilityDraft:     true,  // published can contain drafts, this is how review submissions work.
		visibility.VisibilityUnlisted:  false, // published cannot contain unlisted, unlisted nodes are intended for "personal" use not sharing globally with the entire world, but they can be accessed if given a URL for example.
		visibility.VisibilityReview:    true,  // published can contain review nodes, this is how the submission review process works.
		visibility.VisibilityPublished: true,  // obviously, published can contain other published nodes.
	},
}
