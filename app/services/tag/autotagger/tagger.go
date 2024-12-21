package autotagger

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag/tag_querier"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/generative"
)

type Tagger struct {
	querier *tag_querier.Querier
	tagger  generative.Tagger
}

func New(
	querier *tag_querier.Querier,
	tagger generative.Tagger,
) *Tagger {
	return &Tagger{
		querier: querier,
		tagger:  tagger,
	}
}

func (t *Tagger) Gather(ctx context.Context, content datagraph.Content) (tag_ref.Names, error) {
	available, err := t.querier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := t.tagger.SuggestTags(ctx, content, available.Names())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}
