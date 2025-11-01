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
	// NOTE: See if there's a better way to design the analyser in such a way
	// that can avoid reading the file into memory before checking the fill rule
	// because currently, fill rules are the only use-case for analysing the
	// file but there may be other use-cases in future so it's probably useless.
	rule, ok := fillrule.Get()
	if !ok {
		return nil // no fill rule, nothing to do
	}

	targetNode, ok := rule.TargetNodeID.Get()
	if !ok {
		return fault.New("target node ID not set", fctx.With(ctx))
	}

	node, err := a.nodereader.Probe(ctx, library.NodeID(targetNode))
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

	_, err = a.nodewriter.Update(ctx, library.NewQueryKey(node.Mark), node_mutate.Partial{
		Content: opt.New(rich),
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
