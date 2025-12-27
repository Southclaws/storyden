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
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/post/thread_writer"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

func (s *service) Create(ctx context.Context,
	title string,
	authorID account.AccountID,
	meta map[string]any,
	partial Partial,
) (*thread.Thread, error) {
	opts := partial.Opts()
	opts = append(opts,
		thread_writer.WithMeta(meta),
	)

	// Small hack: default to zero-value of content, which is actually not zero
	// it's <body></body>. Why? who knows... oh, me, yes I should know. I don't.
	if !partial.Content.Ok() {
		c, _ := datagraph.NewRichText("")
		opts = append(opts, thread_writer.WithContent(c))
	}

	if u, ok := partial.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u, fetcher.Options{})
		if err == nil {
			opts = append(opts, thread_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	if tags, ok := partial.Tags.Get(); ok {
		newTags, err := s.tagWriter.Add(ctx, tags...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		tagIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })

		opts = append(opts, thread_writer.WithTagsAdd(tagIDs...))
	}

	thr, err := s.threadWriter.Create(ctx,
		title,
		authorID,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	if content, ok := partial.Content.Get(); ok && thr.Visibility == visibility.VisibilityPublished {
		result, err := s.cpm.CheckContent(ctx, xid.ID(thr.ID), datagraph.KindThread, title, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if result.RequiresReview {
			thr, err = s.threadWriter.Update(ctx, thr.ID, thread_writer.WithVisibility(visibility.VisibilityReview))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	if err := s.cache.Invalidate(ctx, xid.ID(thr.ID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if thr.Visibility == visibility.VisibilityPublished {
		s.bus.Publish(ctx, &message.EventThreadPublished{
			ID: thr.ID,
		})
	}

	// TODO: Do this using event consumer.
	s.mentioner.Send(ctx, authorID, *datagraph.NewRef(thr), thr.Content.References()...)

	return thr, nil
}
