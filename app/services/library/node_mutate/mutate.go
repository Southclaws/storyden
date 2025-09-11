package node_mutate

import (
	"net/url"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
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
	URL          deletable.Value[url.URL]
	Description  opt.Optional[string]
	PrimaryImage deletable.Value[asset.AssetID]
	Content      opt.Optional[datagraph.Content]
	Parent       opt.Optional[library.QueryKey]
	HideChildren opt.Optional[bool]
	Properties   opt.Optional[library.PropertyMutationList]
	Tags         opt.Optional[tag_ref.Names]
	Visibility   opt.Optional[visibility.Visibility]
	Metadata     opt.Optional[map[string]any]
	AssetsAdd    opt.Optional[[]asset.AssetID]
	AssetsRemove opt.Optional[[]asset.AssetID]
	AssetSources opt.Optional[[]string]
}

type Manager struct {
	accountQuery *account_querier.Querier
	nodeQuerier  *node_querier.Querier
	nodeWriter   *node_writer.Writer
	schemaWriter *node_properties.SchemaWriter
	propWriter   *node_properties.Writer
	tagWriter    *tag_writer.Writer
	titler       generative.Titler
	tagger       *autotagger.Tagger
	nc           *node_children.Writer
	fetcher      *fetcher.Fetcher
	summariser   generative.Summariser
	bus          *pubsub.Bus
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	schemaWriter *node_properties.SchemaWriter,
	propWriter *node_properties.Writer,
	tagWriter *tag_writer.Writer,
	titler generative.Titler,
	tagger *autotagger.Tagger,
	nc *node_children.Writer,
	fetcher *fetcher.Fetcher,
	summariser generative.Summariser,
	bus *pubsub.Bus,
) *Manager {
	return &Manager{
		accountQuery: accountQuery,
		nodeQuerier:  nodeQuerier,
		nodeWriter:   nodeWriter,
		schemaWriter: schemaWriter,
		propWriter:   propWriter,
		tagWriter:    tagWriter,
		titler:       titler,
		tagger:       tagger,
		nc:           nc,
		fetcher:      fetcher,
		summariser:   summariser,
		bus:          bus,
	}
}
