package thread

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/post/thread_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
)

func (s *service) Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error) {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thr, err := s.threadQuerier.Get(ctx, threadID, pagination.Parameters{}, opt.NewEmpty[account.AccountID]())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := authoriseThreadUpdate(ctx, acc, thr); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := authoriseMutation(ctx, partial); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	oldVisibility := thr.Visibility
	opts := partial.Opts()

	newContent, contentChanged := partial.Content.Get()
	newTitle, titleChanged := partial.Title.Get()
	if contentChanged || titleChanged {
		result, err := s.cpm.CheckContent(ctx, xid.ID(threadID), datagraph.KindThread, newTitle, newContent)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if result.Action == checker.ActionReport {
			opts = append(opts, thread_writer.WithVisibility(visibility.VisibilityReview))
		}
	}

	if tags, ok := partial.Tags.Get(); ok {
		currentTagNames := thr.Tags.Names()

		toCreate, toRemove := lo.Difference(tags, currentTagNames)

		newTags, err := s.tagWriter.Add(ctx, toCreate...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		addIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })
		removeIDs := dt.Reduce(thr.Tags, func(acc []tag_ref.ID, prev *tag_ref.Tag) []tag_ref.ID {
			if lo.Contains(toRemove, prev.Name) {
				acc = append(acc, prev.ID)
			}
			return acc
		}, []tag_ref.ID{})

		opts = append(opts, thread_writer.WithTagsAdd(addIDs...))
		opts = append(opts, thread_writer.WithTagsRemove(removeIDs...))
	}

	if u, ok := partial.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u, fetcher.Options{})
		if err == nil {
			opts = append(opts, thread_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	if err := s.cache.Invalidate(ctx, xid.ID(threadID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thr, err = s.threadWriter.Update(ctx, threadID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Always emit a general update event
	s.bus.Publish(ctx, &message.EventThreadUpdated{
		ID: thr.ID,
	})

	// Emit visibility-specific events when visibility changes
	if oldVisibility != thr.Visibility {
		if thr.Visibility == visibility.VisibilityPublished {
			s.bus.Publish(ctx, &message.EventThreadPublished{
				ID: thr.ID,
			})
		} else {
			s.bus.Publish(ctx, &message.EventThreadUnpublished{
				ID: thr.ID,
			})
		}
	}

	return thr, nil
}

func authoriseThreadUpdate(ctx context.Context, acc *account.AccountWithEdges, thr *thread.Thread) error {
	return acc.Roles.Permissions().Authorise(ctx, func() error {
		if thr.Author.ID != acc.ID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not author", "You are not the author of the thread and do not have the Manage Posts permission."),
			)
		}
		return nil
	}, rbac.PermissionManagePosts)
}
