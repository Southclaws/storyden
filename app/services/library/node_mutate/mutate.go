package node_mutate

import (
	"net/url"

	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/generative"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
	"github.com/Southclaws/storyden/internal/deletable"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Partial struct {
	Name         opt.Optional[string]
	Slug         opt.Optional[mark.Slug]
	URL          opt.Optional[url.URL]
	PrimaryImage deletable.Value[asset.AssetID]
	Content      opt.Optional[datagraph.Content]
	Parent       opt.Optional[library.QueryKey]
	Properties   opt.Optional[library.PropertyMutationList]
	Tags         opt.Optional[tag_ref.Names]
	Visibility   opt.Optional[visibility.Visibility]
	Metadata     opt.Optional[map[string]any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
	AssetSources opt.Optional[[]string]
	FillSource   opt.Optional[asset.FillSource]
	TitleFill    opt.Optional[datagraph.TitleFillCommand]
	TagFill      opt.Optional[tag.TagFillCommand]
	ContentFill  opt.Optional[asset.ContentFillCommand]
}

type Manager struct {
	logger            *zap.Logger
	accountQuery      *account_querier.Querier
	nodeQuerier       *node_querier.Querier
	nodeWriter        *node_writer.Writer
	schemaWriter      *node_properties.SchemaWriter
	propWriter        *node_properties.Writer
	tagWriter         *tag_writer.Writer
	titler            generative.Titler
	tagger            *autotagger.Tagger
	nc                node_children.Repository
	fetcher           *fetcher.Fetcher
	summariser        generative.Summariser
	indexQueue        pubsub.Topic[mq.IndexNode]
	deleteQueue       pubsub.Topic[mq.DeleteNode]
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset]
}

func New(
	logger *zap.Logger,
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	schemaWriter *node_properties.SchemaWriter,
	propWriter *node_properties.Writer,
	tagWriter *tag_writer.Writer,
	titler generative.Titler,
	tagger *autotagger.Tagger,
	nc node_children.Repository,
	fetcher *fetcher.Fetcher,
	summariser generative.Summariser,
	indexQueue pubsub.Topic[mq.IndexNode],
	deleteQueue pubsub.Topic[mq.DeleteNode],
	assetAnalyseQueue pubsub.Topic[mq.AnalyseAsset],
) *Manager {
	return &Manager{
		logger:            logger,
		accountQuery:      accountQuery,
		nodeQuerier:       nodeQuerier,
		nodeWriter:        nodeWriter,
		schemaWriter:      schemaWriter,
		propWriter:        propWriter,
		tagWriter:         tagWriter,
		titler:            titler,
		tagger:            tagger,
		nc:                nc,
		fetcher:           fetcher,
		summariser:        summariser,
		indexQueue:        indexQueue,
		deleteQueue:       deleteQueue,
		assetAnalyseQueue: assetAnalyseQueue,
	}
}
