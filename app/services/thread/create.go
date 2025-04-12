package thread

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

func (s *service) Create(ctx context.Context,
	title string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	status visibility.Visibility,
	meta map[string]any,
	partial Partial,
) (*thread.Thread, error) {
	if content, ok := partial.Content.Get(); ok {
		if err := s.cpm.CheckContent(ctx, content); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	opts := partial.Opts()
	opts = append(opts,
		thread.WithVisibility(status),
		thread.WithMeta(meta),
	)

	if u, ok := partial.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u, fetcher.Options{})
		if err == nil {
			opts = append(opts, thread.WithLink(xid.ID(ln.ID)))
		}
	}

	if tags, ok := partial.Tags.Get(); ok {
		newTags, err := s.tagWriter.Add(ctx, tags...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		tagIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })

		opts = append(opts, thread.WithTagsAdd(tagIDs...))
	}

	thr, err := s.thread_repo.Create(ctx,
		title,
		authorID,
		categoryID,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	if partial.Visibility.OrZero() == visibility.VisibilityPublished {
		s.indexQueue.PublishAndForget(ctx, mq.IndexThread{
			ID: thr.ID,
		})
	}

	s.fetcher.HydrateContentURLs(ctx, thr)

	s.mentioner.Send(ctx, *datagraph.NewRef(thr), thr.Content.References()...)

	return thr, nil
}
