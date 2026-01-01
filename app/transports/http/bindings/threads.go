package bindings

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Threads struct {
	thread_cache    *thread_cache.Cache
	thread_svc      thread_service.Service
	thread_mark_svc thread_mark.Service
	accountQuery    *account_querier.Querier
	profileQuery    *profile_querier.Querier
}

func NewThreads(
	thread_cache *thread_cache.Cache,
	thread_svc thread_service.Service,
	thread_mark_svc thread_mark.Service,
	accountQuery *account_querier.Querier,
	profileQuery *profile_querier.Querier,
) Threads {
	return Threads{thread_cache, thread_svc, thread_mark_svc, accountQuery, profileQuery}
}

func (i *Threads) ThreadCreate(ctx context.Context, request openapi.ThreadCreateRequestObject) (openapi.ThreadCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), deserialiseThreadStatus)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	tags := opt.Map(opt.NewPtr(request.Body.Tags), func(tags []string) tag_ref.Names {
		return dt.Map(tags, deserialiseTagName)
	})

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	category := opt.NewPtrMap(request.Body.Category, func(cat openapi.Identifier) xid.ID {
		return openapi.ParseID(cat)
	})

	url, err := opt.MapErr(opt.NewPtr(request.Body.Url), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	pinned := opt.NewPtrMap(request.Body.Pinned, func(p openapi.PinnedRank) int {
		return int(p)
	})

	thread, err := i.thread_svc.Create(ctx,
		request.Body.Title,
		accountID,
		meta,
		thread_service.Partial{
			Content:    richContent,
			Category:   category,
			Tags:       tags,
			Visibility: status,
			URL:        url,
			Pinned:     pinned,
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadCreate200JSONResponse{
		ThreadCreateOKJSONResponse: openapi.ThreadCreateOKJSONResponse(serialiseThread(thread)),
	}, nil
}

func (i *Threads) ThreadUpdate(ctx context.Context, request openapi.ThreadUpdateRequestObject) (openapi.ThreadUpdateResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := opt.Map(opt.NewPtr(request.Body.Tags), func(tags []string) tag_ref.Names {
		return dt.Map(tags, deserialiseTagName)
	})

	Visibility, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), deserialiseThreadStatus)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	pinned := opt.NewPtrMap(request.Body.Pinned, func(p openapi.PinnedRank) int {
		return int(p)
	})

	thread, err := i.thread_svc.Update(ctx, postID, thread_service.Partial{
		Title:      opt.NewPtr(request.Body.Title),
		Content:    richContent,
		Tags:       tags,
		Category:   opt.NewPtrMap(request.Body.Category, deserialiseID),
		Visibility: Visibility,
		Pinned:     pinned,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadUpdate200JSONResponse{
		ThreadUpdateOKJSONResponse: openapi.ThreadUpdateOKJSONResponse(serialiseThread(thread)),
	}, nil
}

func (i *Threads) ThreadDelete(ctx context.Context, request openapi.ThreadDeleteRequestObject) (openapi.ThreadDeleteResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = i.thread_svc.Delete(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadDelete200Response{}, nil
}

func (i *Threads) ThreadList(ctx context.Context, request openapi.ThreadListRequestObject) (openapi.ThreadListResponseObject, error) {
	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}

		return max(1, int(v))
	}).Or(1)

	query := opt.NewPtr(request.Params.Q)

	author, err := openapi.OptionalID(ctx, i.profileQuery, request.Params.Author)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	visibilities, err := opt.MapErr(opt.NewPtr(request.Params.Visibility), deserialiseVisibilityList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := opt.NewPtrMap(request.Params.Tags, func(t []openapi.Identifier) []xid.ID {
		return dt.Map(t, func(i openapi.Identifier) xid.ID {
			return openapi.ParseID(i)
		})
	})

	cats := deserialiseCategorySlugQueryParam(request.Params.Categories)
	ignorePinned := opt.NewPtr(request.Params.IgnorePinned)

	page = max(0, page-1)
	result, err := i.thread_svc.List(ctx, page, pageSize, thread_service.Params{
		Query:        query,
		AccountID:    author,
		Visibility:   visibilities,
		Tags:         tags,
		Categories:   cats,
		IgnorePinned: ignorePinned,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	page = result.CurrentPage + 1
	nextPage := opt.Map(result.NextPage, func(i int) int { return i + 1 })

	return openapi.ThreadList200JSONResponse{
		ThreadListOKJSONResponse: openapi.ThreadListOKJSONResponse{
			Body: openapi.ThreadListResult{
				CurrentPage: page,
				NextPage:    nextPage.Ptr(),
				PageSize:    result.PageSize,
				Results:     result.Results,
				Threads:     dt.Map(result.Threads, serialiseThreadReference),
				TotalPages:  result.TotalPages,
			},
			Headers: openapi.ThreadListOKResponseHeaders{
				CacheControl: "no-store",
			},
		},
	}, nil
}

func (i *Threads) ThreadGet(ctx context.Context, request openapi.ThreadGetRequestObject) (openapi.ThreadGetResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	etag, notModified := i.thread_cache.Check(ctx, reqinfo.GetCacheQuery(ctx), xid.ID(postID))
	if notModified {
		return openapi.ThreadGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		}, nil
	}

	pp := deserialisePageParams(request.Params.Page, reply.RepliesPerPage)

	thread, err := i.thread_svc.Get(ctx, postID, pp)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if etag == nil {
		i.thread_cache.Store(ctx, xid.ID(thread.ID), thread.UpdatedAt)
		etag = cachecontrol.NewETag(thread.UpdatedAt)
	}

	return openapi.ThreadGet200JSONResponse{
		ThreadGetJSONResponse: openapi.ThreadGetJSONResponse{
			Body: serialiseThread(thread),
			Headers: openapi.ThreadGetResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		},
	}, nil
}

func deserialiseThreadStatus(in openapi.Visibility) (visibility.Visibility, error) {
	s, err := visibility.NewVisibility(string(in))
	if err != nil {
		return visibility.Visibility{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return s, nil
}

func deserialiseCategorySlugQueryParam(in *openapi.CategorySlugListQuery) opt.Optional[thread_querier.CategoryFilter] {
	// Do not filter by any categorise, return all threads.
	if in == nil {
		return opt.NewEmpty[thread_querier.CategoryFilter]()
	}

	// Fetch uncategorised threads only.
	_, isExplicitlyNull := lo.Find(*in, func(s string) bool { return s == "null" })
	if isExplicitlyNull {
		return opt.New(thread_querier.CategoryFilter{
			Uncategorised: true,
		})
	}

	// Filter by these categories.
	return opt.New(thread_querier.CategoryFilter{
		Slugs:         *in,
		Uncategorised: false,
	})
}
