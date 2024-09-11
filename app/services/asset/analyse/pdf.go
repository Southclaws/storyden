package analyse

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"golang.org/x/net/html"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
)

func (a *Analyser) analysePDF(ctx context.Context, buf []byte, fillrule opt.Optional[asset.ContentFillCommand]) error {
	rule, ok := fillrule.Get()
	if !ok {
		return nil // no fill rule, nothing to do
	}

	node, err := a.nodereader.GetByID(ctx, library.NodeID(rule.TargetNodeID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	result, err := a.pdfextractor.Extract(ctx, buf)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	pr, pw := io.Pipe()
	err = html.Render(pw, result.HTML)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	rich, err := datagraph.NewRichTextFromReader(pr)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	switch rule.FillRule {
	case asset.ContentFillRulePrepend:
		// TODO: Content prepend API that handles HTML properly
		// node.Content.Prepend(rich)

	case asset.ContentFillRuleAppend:
		// TODO: Content append API that handles HTML properly
		// node.Content.Append(rich)

	case asset.ContentFillRuleReplace:
		// rich = rich
	}

	_, err = a.nodewriter.Update(ctx, library.NodeSlug(node.Slug), node_mutate.Partial{
		Content: opt.New(rich),
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
