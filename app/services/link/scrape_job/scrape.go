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
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

type scrapeConsumer struct {
	fetcher *fetcher.Fetcher
	threads thread.Repository
	replies reply.Repository
	nodes   library.Repository
}

func newScrapeConsumer(
	fetcher *fetcher.Fetcher,
	threads thread.Repository,
	replies reply.Repository,
	nodes library.Repository,
) *scrapeConsumer {
	return &scrapeConsumer{
		fetcher: fetcher,
		threads: threads,
		replies: replies,
		nodes:   nodes,
	}
}

func (s *scrapeConsumer) scrapeLink(ctx context.Context, u url.URL, item opt.Optional[datagraph.Item]) error {
	ln, err := s.fetcher.ScrapeAndStore(ctx, u)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if i, ok := item.Get(); ok {
		switch t := i.(type) {
		case *thread.Thread:
			_, err := s.threads.Update(ctx, t.ID, thread.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

		case *reply.Reply:
			_, err := s.replies.Update(ctx, t.ID, reply.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

		case *library.Node:
			_, err := s.nodes.Update(ctx, t.ID, library.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	return nil
}
