package node_children

import (
	"log/slog"

	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	logger *slog.Logger
	db     *ent.Client
	nr     *node_querier.Querier
}

func New(logger *slog.Logger, db *ent.Client, nr *node_querier.Querier) *Writer {
	return &Writer{logger, db, nr}
}
