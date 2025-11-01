package node_comment

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_comment"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/visibility"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
)

var errVisibility = fault.New("visibility")

func Build() fx.Option {
	return fx.Provide(New)
}

type Manager struct {
	repo          *node_comment.Repository
	nodeQuerier   *node_querier.Querier
	threadQuerier *thread_querier.Querier
	threadService thread_service.Service
}

func New(
	repo *node_comment.Repository,
	nodeQuerier *node_querier.Querier,
	threadQuerier *thread_querier.Querier,
	threadService thread_service.Service,
) *Manager {
	return &Manager{
		repo:          repo,
		nodeQuerier:   nodeQuerier,
		threadQuerier: threadQuerier,
		threadService: threadService,
	}
}

func (s *Manager) Create(
	ctx context.Context,
	nodeMark library.QueryKey,
	title string,
	authorID account.AccountID,
	meta map[string]any,
	partial thread_service.Partial,
) (*thread.Thread, error) {
	node, err := s.nodeQuerier.Get(ctx, nodeMark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get node"))
	}

	if node.Visibility != visibility.VisibilityPublished {
		return nil, fault.Wrap(errVisibility,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("node is not published", "Cannot comment on unpublished pages."),
		)
	}

	if !partial.Visibility.Ok() {
		// By default, comments are set to unlisted so they do not appear in
		// feeds and categories. The caller may opt to override this.
		partial.Visibility = opt.New(visibility.VisibilityUnlisted)
	}

	thr, err := s.threadService.Create(ctx, title, authorID, meta, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	err = s.repo.Create(ctx, library.NodeID(node.GetID()), thr.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to link thread to node"))
	}

	return thr, nil
}

func (s *Manager) List(
	ctx context.Context,
	qk library.QueryKey,
	page int,
	size int,
	accountID opt.Optional[account.AccountID],
) (*thread_querier.Result, error) {
	pp := pagination.NewPageParams(uint(page), uint(size))

	paginationResult, err := s.repo.GetThreadIDs(ctx, qk, pp)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(paginationResult.Items) == 0 {
		return &thread_querier.Result{
			PageSize:    size,
			Results:     0,
			TotalPages:  0,
			CurrentPage: page,
			NextPage:    opt.NewEmpty[int](),
			Threads:     []*thread.Thread{},
		}, nil
	}

	result, err := s.threadQuerier.List(ctx, 0, len(paginationResult.Items), accountID, thread_querier.WithIDs(paginationResult.Items...))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}
