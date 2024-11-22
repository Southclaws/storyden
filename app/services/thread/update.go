package thread

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

func (s *service) Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error) {
	if content, ok := partial.Content.Get(); ok {
		if err := s.cpm.CheckContent(ctx, content); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thr, err := s.thread_repo.Get(ctx, threadID, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := authoriseThreadUpdate(ctx, acc, thr); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := partial.Opts()

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

		opts = append(opts, thread.WithTagsAdd(addIDs...))
		opts = append(opts, thread.WithTagsRemove(removeIDs...))
	}

	if u, ok := partial.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u, fetcher.Options{})
		if err == nil {
			opts = append(opts, thread.WithLink(xid.ID(ln.ID)))
		}
	}

	thr, err = s.thread_repo.Update(ctx, threadID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if thr.Visibility == visibility.VisibilityPublished {
		if err := s.indexQueue.Publish(ctx, mq.IndexThread{
			ID: thr.ID,
		}); err != nil {
			s.l.Error("failed to publish index post message", zap.Error(err))
		}
	} else {
		if err := s.deleteQueue.Publish(ctx, mq.DeleteThread{
			ID: thr.ID,
		}); err != nil {
			s.l.Error("failed to publish index post message", zap.Error(err))
		}
	}

	return thr, nil
}

func authoriseThreadUpdate(ctx context.Context, acc *account.Account, thr *thread.Thread) error {
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
