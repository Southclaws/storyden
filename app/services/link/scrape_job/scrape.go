package scrape_job

import (
	"context"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_writer"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

type scrapeConsumer struct {
	fetcher    *fetcher.Fetcher
	posts      *post_writer.PostWriter
	nodeWriter *node_writer.Writer
}

func newScrapeConsumer(
	fetcher *fetcher.Fetcher,
	posts *post_writer.PostWriter,
	nodeWriter *node_writer.Writer,
) *scrapeConsumer {
	return &scrapeConsumer{
		fetcher:    fetcher,
		posts:      posts,
		nodeWriter: nodeWriter,
	}
}

func (s *scrapeConsumer) scrapeLink(ctx context.Context, u url.URL, item opt.Optional[datagraph.Ref]) error {
	ln, err := s.fetcher.ScrapeAndStore(ctx, u)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if i, ok := item.Get(); ok {
		switch i.Kind {
		case datagraph.KindPost:
			_, err := s.posts.Update(ctx, post.ID(i.ID), post_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

		case datagraph.KindNode:
			qk := library.QueryKey{mark.NewQueryKeyID(i.ID)}
			_, err := s.nodeWriter.Update(ctx, qk, node_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	return nil
}
