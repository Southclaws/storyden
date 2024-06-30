package link_getter

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph/link_graph"
	"github.com/Southclaws/storyden/app/services/semdex"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Getter struct {
	l   *zap.Logger
	lg  link_graph.Repository
	rec semdex.Recommender
}

func New(
	l *zap.Logger,
	lg link_graph.Repository,
	rec semdex.Recommender,
) *Getter {
	return &Getter{
		l:   l.With(zap.String("service", "link_retreival")),
		lg:  lg,
		rec: rec,
	}
}

func (s *Getter) Get(ctx context.Context, slug string) (*link_graph.WithRefs, error) {
	ln, err := s.lg.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ln, nil
}
