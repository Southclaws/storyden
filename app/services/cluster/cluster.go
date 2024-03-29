package cluster

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster_children"
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
	) (*datagraph.Cluster, error)

	Get(ctx context.Context, slug datagraph.ClusterSlug) (*datagraph.Cluster, error)
	Update(ctx context.Context, slug datagraph.ClusterSlug, p Partial) (*datagraph.Cluster, error)
	Delete(ctx context.Context, slug datagraph.ClusterSlug, d DeleteOptions) (*datagraph.Cluster, error)
}

type Partial struct {
	Name         opt.Optional[string]
	Slug         opt.Optional[string]
	URL          opt.Optional[string]
	Description  opt.Optional[string]
	Content      opt.Optional[string]
	Parent       opt.Optional[datagraph.ClusterSlug]
	Visibility   opt.Optional[post.Visibility]
	Properties   opt.Optional[any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
}

type DeleteOptions struct {
	MoveTo   opt.Optional[datagraph.ClusterSlug]
	Clusters bool
	Items    bool
}

func (p Partial) Opts() (opts []cluster.Option) {
	p.Name.Call(func(value string) { opts = append(opts, cluster.WithName(value)) })
	p.Slug.Call(func(value string) { opts = append(opts, cluster.WithSlug(value)) })
	p.Description.Call(func(value string) { opts = append(opts, cluster.WithDescription(value)) })
	p.Content.Call(func(value string) { opts = append(opts, cluster.WithContent(value)) })
	p.Properties.Call(func(value any) { opts = append(opts, cluster.WithProperties(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, cluster.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, cluster.WithAssetsRemoved(value)) })
	return
}

type service struct {
	ar       account.Repository
	cr       cluster.Repository
	cc       cluster_children.Repository
	hydrator hydrator.Service
}

func New(
	ar account.Repository,
	cr cluster.Repository,
	cc cluster_children.Repository,
	hydrator hydrator.Service,
) Manager {
	return &service{
		ar:       ar,
		cr:       cr,
		cc:       cc,
		hydrator: hydrator,
	}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	slug string,
	desc string,
	p Partial,
) (*datagraph.Cluster, error) {
	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := s.cr.Create(ctx, owner, name, slug, desc, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return clus, nil
}

func (s *service) Get(ctx context.Context, slug datagraph.ClusterSlug) (*datagraph.Cluster, error) {
	clus, err := s.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return clus, nil
}

func (s *service) Update(ctx context.Context, slug datagraph.ClusterSlug, p Partial) (*datagraph.Cluster, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := s.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !clus.Owner.Admin {
		if clus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err = s.cr.Update(ctx, clus.ID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return clus, nil
}

func (s *service) Delete(ctx context.Context, slug datagraph.ClusterSlug, d DeleteOptions) (*datagraph.Cluster, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := s.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !clus.Owner.Admin {
		if clus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	destination, err := opt.MapErr(d.MoveTo, func(target datagraph.ClusterSlug) (datagraph.Cluster, error) {
		opts := []cluster_children.Option{}

		if d.Clusters {
			opts = append(opts, cluster_children.MoveClusters())
		}

		if d.Items {
			opts = append(opts, cluster_children.MoveItems())
		}

		destination, err := s.cc.Move(ctx, slug, target, opts...)
		if err != nil {
			return datagraph.Cluster{}, fault.Wrap(err, fctx.With(ctx))
		}

		return *destination, fault.Wrap(err, fctx.With(ctx))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.cr.Delete(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return destination.Ptr(), nil
}

func (s *service) hydrateLink(ctx context.Context, partial Partial) (opts []cluster.Option) {
	text, textOK := partial.Content.Get()

	if !textOK && !partial.URL.Ok() {
		return
	}

	return s.hydrator.HydrateCluster(ctx, text, partial.URL)
}

func (s *service) applyOpts(ctx context.Context, p Partial) ([]cluster.Option, error) {
	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (*account.Account, error) {
		return s.ar.GetByID(ctx, aid)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := p.Opts()

	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.cr.Get(ctx, parentSlug)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, cluster.WithParent(parent.ID))
	}

	if acc, ok := acc.Get(); ok {
		p.Visibility.Call(func(value post.Visibility) {
			// Only admins can immediately post to the public feed.
			if value == post.VisibilityPublished && !acc.Admin {
				return
			}

			opts = append(opts, cluster.WithVisibility(value))
		})
	}

	opts = append(opts, s.hydrateLink(ctx, p)...)

	return opts, nil
}
