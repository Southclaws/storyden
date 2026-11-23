package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
)

func (s *Manager) prepareParent(ctx context.Context, p Partial, current opt.Optional[library.Node]) (opt.Optional[library.NodeID], error) {
	childVisibility := visibility.VisibilityDraft
	currentParent := opt.NewEmpty[library.Node]()
	visibilityChanged := p.Visibility.Ok()
	parentID := opt.NewEmpty[library.NodeID]()
	if node, ok := current.Get(); ok {
		childVisibility = node.Visibility
		currentParent = node.Parent
		visibilityChanged = p.Visibility.Or(node.Visibility) != node.Visibility
	}
	childVisibility = p.Visibility.Or(childVisibility)

	parentToAuthorise := opt.NewEmpty[library.Node]()
	if visibilityChanged {
		parentToAuthorise = currentParent
	}

	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.nodeQuerier.Get(ctx, parentSlug)
		if err != nil {
			return parentID, fault.Wrap(err, fctx.With(ctx))
		}

		previousParent, hasParent := currentParent.Get()
		parentChanged := !hasParent || previousParent.GetID() != parent.GetID()
		if parentChanged || visibilityChanged {
			parentToAuthorise = opt.New(*parent)
		}

		parentID = opt.New(library.NodeID(parent.GetID()))
	}

	if parent, ok := parentToAuthorise.Get(); ok {
		if err := nodetree.AuthoriseParent(ctx, &parent, childVisibility); err != nil {
			return opt.NewEmpty[library.NodeID](), fault.Wrap(err, fctx.With(ctx))
		}
	}

	return parentID, nil
}
