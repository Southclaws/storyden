package node_mutate

import (
	"context"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
	"github.com/Southclaws/storyden/internal/deletable"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Partial struct {
	Name             opt.Optional[string]
	Slug             opt.Optional[mark.Slug]
	URL              opt.Optional[url.URL]
	PrimaryImage     deletable.Value[asset.AssetID]
	Content          opt.Optional[datagraph.Content]
	Parent           opt.Optional[library.QueryKey]
	Tags             opt.Optional[tag_ref.Names]
	Visibility       opt.Optional[visibility.Visibility]
	Metadata         opt.Optional[map[string]any]
	AssetsAdd        opt.Optional[[]asset.AssetID]
	AssetsRemove     opt.Optional[[]asset.AssetID]
	AssetSources     opt.Optional[[]string]
	ContentFill      opt.Optional[asset.ContentFillCommand]
	ContentSummarise opt.Optional[bool]
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

type Manager struct {
	logger            *zap.Logger
	accountQuery      *account_querier.Querier
	nodeQuerier       *node_querier.Querier
	nodeWriter        *node_writer.Writer
	tagWriter         *tag_writer.Writer
	tagger            *autotagger.Tagger
	nc                node_children.Repository
	fetcher           *fetcher.Fetcher
	indexQueue        pubsub.Topic[mq.IndexNode]
	deleteQueue       pubsub.Topic[mq.DeleteNode]
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset]
}

func New(
	logger *zap.Logger,
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	tagWriter *tag_writer.Writer,
	tagger *autotagger.Tagger,
	nc node_children.Repository,
	fetcher *fetcher.Fetcher,
	indexQueue pubsub.Topic[mq.IndexNode],
	deleteQueue pubsub.Topic[mq.DeleteNode],
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset],
) *Manager {
	return &Manager{
		logger:            logger,
		accountQuery:      accountQuery,
		nodeQuerier:       nodeQuerier,
		nodeWriter:        nodeWriter,
		tagWriter:         tagWriter,
		tagger:            tagger,
		nc:                nc,
		fetcher:           fetcher,
		indexQueue:        indexQueue,
		deleteQueue:       deleteQueue,
		assetAnalyseQueue: assetAnalyseQueue,
	}
}

func (s *Manager) applyOpts(ctx context.Context, p Partial) ([]node_writer.Option, error) {
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
