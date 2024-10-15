package node_mutate

import (
	"context"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	library_service "github.com/Southclaws/storyden/app/services/library"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/internal/deletable"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		p Partial,
	) (*library.Node, error)

	Update(ctx context.Context, slug library.QueryKey, p Partial) (*library.Node, error)
	Delete(ctx context.Context, slug library.QueryKey, d DeleteOptions) (*library.Node, error)
}

type Partial struct {
	Name         opt.Optional[string]
	Slug         opt.Optional[mark.Slug]
	URL          opt.Optional[url.URL]
	PrimaryImage deletable.Value[asset.AssetID]
	Content      opt.Optional[datagraph.Content]
	Parent       opt.Optional[library.QueryKey]
	Visibility   opt.Optional[visibility.Visibility]
	Metadata     opt.Optional[map[string]any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
	AssetSources opt.Optional[[]string]
	ContentFill  opt.Optional[asset.ContentFillCommand]
}

type DeleteOptions struct {
	NewParent opt.Optional[library.QueryKey]
}

func (p Partial) Opts() (opts []node_writer.Option) {
	p.Name.Call(func(value string) { opts = append(opts, node_writer.WithName(value)) })
	p.Slug.Call(func(value mark.Slug) { opts = append(opts, node_writer.WithSlug(value.String())) })
	p.PrimaryImage.Call(func(value xid.ID) {
		opts = append(opts, node_writer.WithPrimaryImage(value))
	}, func() {
		opts = append(opts, node_writer.WithPrimaryImageRemoved())
	})
	p.Content.Call(func(value datagraph.Content) { opts = append(opts, node_writer.WithContent(value)) })
	p.Metadata.Call(func(value map[string]any) { opts = append(opts, node_writer.WithMetadata(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssetsRemoved(value)) })
	p.Visibility.Call(func(value visibility.Visibility) { opts = append(opts, node_writer.WithVisibility(value)) })
	return
}

type service struct {
	accountQuery      *account_querier.Querier
	nodeQuerier       *node_querier.Querier
	nodeWriter        *node_writer.Writer
	nc                node_children.Repository
	fetcher           *fetcher.Fetcher
	fs                *fetcher.Fetcher
	indexQueue        pubsub.Topic[mq.IndexNode]
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset]
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	nc node_children.Repository,
	fetcher *fetcher.Fetcher,
	fs *fetcher.Fetcher,
	indexQueue pubsub.Topic[mq.IndexNode],
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset],
) Manager {
	return &service{
		accountQuery:      accountQuery,
		nodeQuerier:       nodeQuerier,
		nodeWriter:        nodeWriter,
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

			if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionManageLibrary); err != nil {
				return nil, fault.Wrap(err,
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

			opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	nodeSlug := p.Slug.Or(mark.NewSlugFromName(name))

	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	n, err := s.nodeWriter.Create(ctx, owner, name, nodeSlug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: library.NodeID(n.Mark.ID())}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}

func (s *service) Update(ctx context.Context, qk library.QueryKey, p Partial) (*library.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := library_service.AuthoriseNodeMutation(ctx, acc, n); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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

			opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
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
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	n, err = s.nodeWriter.Update(ctx, qk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: library.NodeID(n.Mark.ID())}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}

func (s *service) Delete(ctx context.Context, qk library.QueryKey, d DeleteOptions) (*library.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := library_service.AuthoriseNodeMutation(ctx, acc, n); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	destination, err := opt.MapErr(d.NewParent, func(target library.QueryKey) (library.Node, error) {
		destination, err := s.nc.Move(ctx, qk, target)
		if err != nil {
			return library.Node{}, fault.Wrap(err, fctx.With(ctx))
		}

		return *destination, fault.Wrap(err, fctx.With(ctx))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.nodeWriter.Delete(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return destination.Ptr(), nil
}

func (s *service) applyOpts(ctx context.Context, p Partial) ([]node_writer.Option, error) {
	opts := p.Opts()

	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.nodeQuerier.Get(ctx, parentSlug)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node_writer.WithParent(library.NodeID(parent.Mark.ID())))
	}

	return opts, nil
}
