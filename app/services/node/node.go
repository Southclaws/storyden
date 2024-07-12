package library

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/gosimple/slug"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/hydrator"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
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
	Content      opt.Optional[content.Rich]
	Parent       opt.Optional[datagraph.NodeSlug]
	Visibility   opt.Optional[visibility.Visibility]
	Metadata     opt.Optional[map[string]any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
	AssetSources opt.Optional[[]string]
}

type DeleteOptions struct {
	NewParent opt.Optional[datagraph.NodeSlug]
}

func (p Partial) Opts() (opts []library.Option) {
	p.Name.Call(func(value string) { opts = append(opts, library.WithName(value)) })
	p.Slug.Call(func(value string) { opts = append(opts, library.WithSlug(value)) })
	p.Content.Call(func(value content.Rich) { opts = append(opts, library.WithContent(value)) })
	p.Metadata.Call(func(value map[string]any) { opts = append(opts, library.WithMetadata(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, library.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, library.WithAssetsRemoved(value)) })
	return
}

type service struct {
	ar       account.Repository
	nr       library.Repository
	nc       node_children.Repository
	hydrator hydrator.Service
	fs       fetcher.Service
}

func New(
	ar account.Repository,
	nr library.Repository,
	nc node_children.Repository,
	hydrator hydrator.Service,
	fs fetcher.Service,
) Manager {
	return &service{
		ar:       ar,
		nr:       nr,
		nc:       nc,
		hydrator: hydrator,
		fs:       fs,
	}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	p Partial,
) (*datagraph.Node, error) {
	if v, ok := p.Visibility.Get(); ok {
		if v == visibility.VisibilityPublished {
			acc, err := s.ar.GetByID(ctx, owner)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			if !acc.Admin {
				return nil, fault.Wrap(errNotAuthorised,
					fctx.With(ctx),
					fmsg.WithDesc("non admin cannot publish nodes", "You do not have permission to publish, please submit as draft, review or unlisted."),
				)
			}
		}
	}

	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v, ok := p.AssetSources.Get(); ok {
		for _, source := range v {
			a, err := s.fs.Copy(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, library.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	nodeSlug := p.Slug.Or(slug.Make(name))

	n, err := s.nr.Create(ctx, owner, name, nodeSlug, opts...)
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

	if v, ok := p.AssetSources.Get(); ok {
		for _, source := range v {
			a, err := s.fs.Copy(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, library.WithAssets([]asset.AssetID{a.ID}))
		}
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

func (s *service) hydrateLink(ctx context.Context, partial Partial) (opts []library.Option) {
	text, textOK := partial.Content.Get()

	if !textOK && !partial.URL.Ok() {
		return
	}

	return s.hydrator.HydrateNode(ctx, text, partial.URL)
}

func (s *service) applyOpts(ctx context.Context, p Partial) ([]library.Option, error) {
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

		opts = append(opts, library.WithParent(parent.ID))
	}

	if acc, ok := acc.Get(); ok {
		p.Visibility.Call(func(value visibility.Visibility) {
			// Only admins can immediately post to the public feed.
			if value == visibility.VisibilityPublished && !acc.Admin {
				return
			}

			opts = append(opts, library.WithVisibility(value))
		})
	}

	opts = append(opts, s.hydrateLink(ctx, p)...)

	return opts, nil
}
