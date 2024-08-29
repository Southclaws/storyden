package node_mutate

import (
	"context"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/gosimple/slug"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		p Partial,
	) (*library.Node, error)

	Update(ctx context.Context, slug library.NodeSlug, p Partial) (*library.Node, error)
	Delete(ctx context.Context, slug library.NodeSlug, d DeleteOptions) (*library.Node, error)
}

type Partial struct {
	Name         opt.Optional[string]
	Slug         opt.Optional[string]
	URL          opt.Optional[url.URL]
	Content      opt.Optional[content.Rich]
	Parent       opt.Optional[library.NodeSlug]
	Visibility   opt.Optional[visibility.Visibility]
	Metadata     opt.Optional[map[string]any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
	AssetSources opt.Optional[[]string]
	ContentFill  opt.Optional[asset.ContentFillCommand]
}

type DeleteOptions struct {
	NewParent opt.Optional[library.NodeSlug]
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
	accountQuery      account_querier.Querier
	nr                library.Repository
	nc                node_children.Repository
	fetcher           *fetcher.Fetcher
	fs                *fetcher.Fetcher
	indexQueue        pubsub.Topic[mq.IndexNode]
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset]
}

func New(
	accountQuery account_querier.Querier,
	nr library.Repository,
	nc node_children.Repository,
	fetcher *fetcher.Fetcher,
	fs *fetcher.Fetcher,
	indexQueue pubsub.Topic[mq.IndexNode],
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset],
) Manager {
	return &service{
		accountQuery:      accountQuery,
		nr:                nr,
		nc:                nc,
		fetcher:           fetcher,
		fs:                fs,
		indexQueue:        indexQueue,
		assetAnalyseQueue: assetAnalyseQueue,
	}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	p Partial,
) (*library.Node, error) {
	if v, ok := p.Visibility.Get(); ok {
		if v == visibility.VisibilityPublished {
			acc, err := s.accountQuery.GetByID(ctx, owner)
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
			a, err := s.fs.CopyAsset(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, library.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	nodeSlug := p.Slug.Or(slug.Make(name))

	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, library.WithLink(xid.ID(ln.ID)))
		}
	}

	n, err := s.nr.Create(ctx, owner, name, nodeSlug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: n.ID}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}

func (s *service) Update(ctx context.Context, slug library.NodeSlug, p Partial) (*library.Node, error) {
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

	// TODO: Queue this for background processing
	if v, ok := p.AssetSources.Get(); ok {
		for _, source := range v {
			a, err := s.fs.CopyAsset(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, library.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	assetsAdd, assetsAddSet := p.AssetsAdd.Get()
	if assetsAddSet && p.ContentFill.Ok() {

		messages := dt.Map(assetsAdd, func(a asset.AssetID) mq.AnalyseAsset {
			return mq.AnalyseAsset{
				AssetID:         a,
				ContentFillRule: p.ContentFill,
			}
		})

		if err := s.assetAnalyseQueue.Publish(ctx, messages...); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, library.WithLink(xid.ID(ln.ID)))
		}
	}

	n, err = s.nr.Update(ctx, n.ID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: n.ID}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}

func (s *service) Delete(ctx context.Context, slug library.NodeSlug, d DeleteOptions) (*library.Node, error) {
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

	destination, err := opt.MapErr(d.NewParent, func(target library.NodeSlug) (library.Node, error) {
		destination, err := s.nc.Move(ctx, slug, target)
		if err != nil {
			return library.Node{}, fault.Wrap(err, fctx.With(ctx))
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

func (s *service) applyOpts(ctx context.Context, p Partial) ([]library.Option, error) {
	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (*account.Account, error) {
		return s.accountQuery.GetByID(ctx, aid)
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

	return opts, nil
}
