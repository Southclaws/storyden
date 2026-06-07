package node_versioning

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_writer"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Service struct {
	nodeMutator    *node_mutate.Manager
	nodeReader     *node_read.HydratedQuerier
	versionQuerier *node_version_querier.Querier
	versionWriter  *node_version_writer.Writer
	bus            *pubsub.Bus
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	nodeMutator *node_mutate.Manager,
	nodeReader *node_read.HydratedQuerier,
	nodeQuerier *node_querier.Querier,
	versionQuerier *node_version_querier.Querier,
	versionWriter *node_version_writer.Writer,
	bus *pubsub.Bus,
	notifier *notify.Notifier,
) *Service {
	s := &Service{
		nodeMutator:    nodeMutator,
		nodeReader:     nodeReader,
		versionQuerier: versionQuerier,
		versionWriter:  versionWriter,
		bus:            bus,
	}

	s.subscribeNotifications(ctx, lc, nodeQuerier, notifier)

	return s
}
