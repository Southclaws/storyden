package analyse_job

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/asset/analyse"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type analyseConsumer struct {
	queue    pubsub.Topic[mq.AnalyseAsset]
	analyser *analyse.Analyser
}

func newAnalyseConsumer(
	queue pubsub.Topic[mq.AnalyseAsset],
	analyser *analyse.Analyser,
) *analyseConsumer {
	return &analyseConsumer{
		queue:    queue,
		analyser: analyser,
	}
}

func (i *analyseConsumer) analyseAsset(ctx context.Context, id asset.AssetID, fillrule opt.Optional[asset.ContentFillCommand]) error {
	return i.analyser.Analyse(ctx, id, fillrule)
}
