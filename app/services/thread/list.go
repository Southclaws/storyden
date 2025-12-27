package thread

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Params struct {
	Query         opt.Optional[string]
	CreatedBefore opt.Optional[time.Time]
	UpdatedBefore opt.Optional[time.Time]
	AccountID     opt.Optional[account.AccountID]
	Visibility    opt.Optional[[]visibility.Visibility]
	Tags          opt.Optional[[]xid.ID]
	Categories    opt.Optional[thread_querier.CategoryFilter]
}

func (s *service) List(ctx context.Context,
	page int,
	size int,
	opts Params,
) (*thread_querier.Result, error) {
	accountID := session.GetOptAccountID(ctx)

	q := []thread_querier.Query{
		thread_querier.HasNotBeenDeleted(),
	}

	opts.Query.Call(func(value string) { q = append(q, thread_querier.HasKeyword(value)) })
	opts.CreatedBefore.Call(func(value time.Time) { q = append(q, thread_querier.HasCreatedDateBefore(value)) })
	opts.UpdatedBefore.Call(func(value time.Time) { q = append(q, thread_querier.HasUpdatedDateBefore(value)) })
	opts.AccountID.Call(func(a account.AccountID) { q = append(q, thread_querier.HasAuthor(a)) })
	opts.Tags.Call(func(a []xid.ID) { q = append(q, thread_querier.HasTags(a)) })
	opts.Categories.Call(func(cf thread_querier.CategoryFilter) { q = append(q, thread_querier.HasCategories(cf)) })

	vq := func() thread_querier.Query {
		v, ok := opts.Visibility.Get()
		if !ok {
			// NOTE: In the default path, we're querying for the main feed, and
			// there are no additional filtering options specified so the engine
			// builds a query specifically for moderators to see in-review posts
			// as well as authors to see their own posts in-review while other
			// members will just see published posts.
			isModerator := false
			if accountID.Ok() {
				roles := session.GetRoles(ctx)
				isModerator = roles.Permissions().HasAny(rbac.PermissionManagePosts, rbac.PermissionAdministrator)
			}
			return thread_querier.HasPublishedOrOwnInReview(accountID, isModerator)
		}

		onlyRequestingPublished := len(v) == 1 && v[0] == visibility.VisibilityPublished
		if onlyRequestingPublished {
			return thread_querier.HasStatus(visibility.VisibilityPublished)
		}

		filterByAccount, ok := opts.AccountID.Get()
		if !ok {
			// Not filtering by specific account - check if user has permission to see all review threads
			if accountID.Ok() {
				roles := session.GetRoles(ctx)
				if roles.Permissions().HasAny(rbac.PermissionManagePosts, rbac.PermissionAdministrator) {
					return thread_querier.HasStatus(v...)
				}
			}

			return thread_querier.HasStatus(visibility.VisibilityPublished)
		}

		requestedByAccount, ok := accountID.Get()
		if !ok {
			return thread_querier.HasStatus(visibility.VisibilityPublished)
		}

		requestingOwnThreads := filterByAccount == requestedByAccount

		if !requestingOwnThreads {
			// Viewing someone else's threads - check if user has permission
			roles := session.GetRoles(ctx)
			if roles.Permissions().HasAny(rbac.PermissionManagePosts, rbac.PermissionAdministrator) {
				return thread_querier.HasStatus(v...)
			}

			return thread_querier.HasStatus(visibility.VisibilityPublished)
		}

		// Viewing own threads - allow all visibilities
		return thread_querier.HasStatus(v...)
	}()
	q = append(q, vq)

	thr, err := s.threadQuerier.List(ctx, page, size, accountID, q...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to list threads"))
	}

	return thr, nil
}
