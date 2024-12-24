package node_semdex

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
)

func (i *semdexer) autofill(ctx context.Context, id library.NodeID, autotitle bool, autotag bool) error {
	qk := library.NewID(xid.ID(id))

	p := node_mutate.Partial{}

	if autotitle {
		p.TitleFill = opt.New(datagraph.TitleFillCommand{FillRule: datagraph.TitleFillRuleReplace})
	}

	if autotag {
		p.TagFill = opt.New(tag.TagFillCommand{FillRule: tag.TagFillRuleReplace})
	}

	_, err := i.nodeUpdater.Update(ctx, qk, p)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
