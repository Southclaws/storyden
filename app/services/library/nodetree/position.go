package nodetree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/internal/deletable"
)

var ErrNoParent = fault.New("node has no parent", ftag.With(ftag.InvalidArgument))

type Position struct {
	nodeChildren *node_children.Writer
	nodeQuerier  *node_querier.Querier
	nodeWriter   *node_writer.Writer
	graph        Graph
	accountQuery *account_querier.Querier
}

func NewPositionService(
	nodeChildren *node_children.Writer,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	graph Graph,
	accountQuery *account_querier.Querier,
) *Position {
	return &Position{
		nodeChildren: nodeChildren,
		nodeQuerier:  nodeQuerier,
		nodeWriter:   nodeWriter,
		graph:        graph,
		accountQuery: accountQuery,
	}
}

type Options struct {
	Parent deletable.Value[library.QueryKey]
	Before opt.Optional[library.NodeID]
	After  opt.Optional[library.NodeID]
	Index  opt.Optional[int]
}

func (p *Position) Move(ctx context.Context, nm library.QueryKey, opts Options) (*library.Node, error) {
	// Move the node to the new position before dealing with sort order.
	parent, sever := opts.Parent.Get()

	thisnode, err := p.nodeQuerier.Get(ctx, nm)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// If the parent ID is explicitly set to null, we need to sever the node.
	if sever {
		// HACK: Because the Sever API is poorly composable, it requires the ID
		// of the parent node to be passed in, so we must query it wastefully...

		var err error

		parentNode, ok := thisnode.Parent.Get()
		if !ok {
			return nil, fault.Wrap(ErrNoParent, fctx.With(ctx))
		}

		thisnode, err = p.graph.Sever(ctx, nm, library.QueryKey{parentNode.Mark.Queryable()})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else if parentNode, ok := parent.Get(); ok {
		// If the parent ID is explicitly set to a value, move the node.

		var err error

		thisnode, err = p.graph.Move(ctx, nm, library.QueryKey{parentNode.Queryable})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// Now handle re-ordering of the node using the before/after/index params.

	if beforeID, ok := opts.Before.Get(); ok {
		n, err := p.nodeChildren.MoveBefore(ctx, thisnode, beforeID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return n, nil
	}

	if afterID, ok := opts.After.Get(); ok {
		n, err := p.nodeChildren.MoveAfter(ctx, thisnode, afterID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return n, nil
	}

	if index, ok := opts.Index.Get(); ok {
		n, err := p.nodeChildren.MoveIndex(ctx, thisnode, index)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return n, nil
	}

	return thisnode, nil
}
