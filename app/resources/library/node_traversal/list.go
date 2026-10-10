package node_traversal

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
)

type ListParams struct {
	CreatedSince  opt.Optional[time.Time]
	CreatedBefore opt.Optional[time.Time]
	NewestFirst   opt.Optional[bool]
}

type ListResult struct {
	pagination.Result[*library.Node]
	TotalResults int
}

func (d *database) List(ctx context.Context, page pagination.Parameters, params ListParams, opts ...Filter) (*ListResult, error) {
	opts = append(opts, func(f *filters) { f.includeHiddenChildren = true })
	nodes, err := d.Subtree(ctx, opt.NewEmpty[library.NodeID](), true, opts...)
	if err != nil {
		return nil, err
	}

	nodes = slices.DeleteFunc(nodes, func(n *library.Node) bool {
		if since, ok := params.CreatedSince.Get(); ok && n.CreatedAt.Before(since) {
			return true
		}
		if before, ok := params.CreatedBefore.Get(); ok && !n.CreatedAt.Before(before) {
			return true
		}
		return false
	})

	if newest, ok := params.NewestFirst.Get(); ok {
		slices.SortFunc(nodes, func(a, b *library.Node) int {
			order := a.CreatedAt.Compare(b.CreatedAt)
			if order == 0 {
				order = strings.Compare(a.GetID().String(), b.GetID().String())
			}
			if newest {
				return -order
			}
			return order
		})
	}

	total := len(nodes)
	start := min(page.Offset(), total)
	end := start + min(page.Limit(), total-start)
	return &ListResult{
		Result:       pagination.NewPageResult(page, total, nodes[start:end]),
		TotalResults: total,
	}, nil
}
