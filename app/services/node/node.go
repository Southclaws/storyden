package node

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/datagraph/node_children"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/hydrator"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		p Partial,
	) (*datagraph.Node, error)

	Get(ctx context.Context, slug datagraph.NodeSlug) (*datagraph.Node, error)
	Update(ctx context.Context, slug datagraph.NodeSlug, p Partial) (*datagraph.Node, error)
	Delete(ctx context.Context, slug datagraph.NodeSlug, d DeleteOptions) (*datagraph.Node, error)
}

type Partial struct {
	Name         opt.Optional[string]
	Slug         opt.Optional[string]
	URL          opt.Optional[string]
	Description  opt.Optional[string]
	Content      opt.Optional[string]
	Parent       opt.Optional[datagraph.NodeSlug]
	Visibility   opt.Optional[post.Visibility]
	Properties   opt.Optional[any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
}

type DeleteOptions struct {
	NewParent opt.Optional[datagraph.NodeSlug]
}

func (p Partial) Opts() (opts []node.Option) {
	p.Name.Call(func(value string) { opts = append(opts, node.WithName(value)) })
	p.Slug.Call(func(value string) { opts = append(opts, node.WithSlug(value)) })
	p.Description.Call(func(value string) { opts = append(opts, node.WithDescription(value)) })
	p.Content.Call(func(value string) { opts = append(opts, node.WithContent(value)) })
	p.Properties.Call(func(value any) { opts = append(opts, node.WithProperties(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, node.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, node.WithAssetsRemoved(value)) })
	return
}

type service struct {
	ar       account.Repository
	nr       node.Repository
	nc       node_children.Repository
	hydrator hydrator.Service
}

func New(
	ar account.Repository,
	nr node.Repository,
	nc node_children.Repository,
	hydrator hydrator.Service,
) Manager {
	return &service{
		ar:       ar,
		nr:       nr,
		nc:       nc,
		hydrator: hydrator,
	}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	slug string,
	desc string,
	p Partial,
) (*datagraph.Node, error) {
	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nr.Create(ctx, owner, name, slug, desc, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return n, nil
}

func (s *service) Get(ctx context.Context, slug datagraph.NodeSlug) (*datagraph.Node, error) {
	n, err := s.nr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return n, nil
}

func (s *service) Update(ctx context.Context, slug datagraph.NodeSlug, p Partial) (*datagraph.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !n.Owner.Admin {
		if n.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err = s.nr.Update(ctx, n.ID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return n, nil
}

func (s *service) Delete(ctx context.Context, slug datagraph.NodeSlug, d DeleteOptions) (*datagraph.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !n.Owner.Admin {
		if n.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	destination, err := opt.MapErr(d.NewParent, func(target datagraph.NodeSlug) (datagraph.Node, error) {
		destination, err := s.nc.Move(ctx, slug, target)
		if err != nil {
			return datagraph.Node{}, fault.Wrap(err, fctx.With(ctx))
		}

		return *destination, fault.Wrap(err, fctx.With(ctx))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.nr.Delete(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return destination.Ptr(), nil
}

func (s *service) hydrateLink(ctx context.Context, partial Partial) (opts []node.Option) {
	text, textOK := partial.Content.Get()

	if !textOK && !partial.URL.Ok() {
		return
	}

	return s.hydrator.HydrateNode(ctx, text, partial.URL)
}

func (s *service) applyOpts(ctx context.Context, p Partial) ([]node.Option, error) {
	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (*account.Account, error) {
		return s.ar.GetByID(ctx, aid)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := p.Opts()

	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.nr.Get(ctx, parentSlug)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node.WithParent(parent.ID))
	}

	if acc, ok := acc.Get(); ok {
		p.Visibility.Call(func(value post.Visibility) {
			// Only admins can immediately post to the public feed.
			if value == post.VisibilityPublished && !acc.Admin {
				return
			}

			opts = append(opts, node.WithVisibility(value))
		})
	}

	opts = append(opts, s.hydrateLink(ctx, p)...)

	return opts, nil
}
