package report_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/internal/ent"
	entreport "github.com/Southclaws/storyden/internal/ent/report"
)

type Querier struct {
	db             *ent.Client
	postSearcher   post_search.Repository
	profileQuerier *profile_querier.Querier
	nodeQuerier    *node_querier.Querier
}

func New(
	db *ent.Client,
	postSearcher post_search.Repository,
	profileQuerier *profile_querier.Querier,
	nodeQuerier *node_querier.Querier,
) *Querier {
	return &Querier{
		db:             db,
		postSearcher:   postSearcher,
		profileQuerier: profileQuerier,
		nodeQuerier:    nodeQuerier,
	}
}

type Query func(*ent.ReportQuery)

func WithStatus(statuses ...report.Status) Query {
	return func(q *ent.ReportQuery) {
		if len(statuses) == 0 {
			return
		}
		statusStrings := dt.Map(statuses, func(s report.Status) string {
			return s.String()
		})
		q.Where(entreport.StatusIn(statusStrings...))
	}
}

func WithKind(kinds ...datagraph.Kind) Query {
	return func(q *ent.ReportQuery) {
		if len(kinds) == 0 {
			return
		}
		kindStrings := dt.Map(kinds, func(k datagraph.Kind) string {
			return k.String()
		})
		q.Where(entreport.TargetKindIn(kindStrings...))
	}
}

func WithReporter(accountID account.AccountID) Query {
	return func(q *ent.ReportQuery) {
		q.Where(entreport.ReportedByID(xid.ID(accountID)))
	}
}

func (q *Querier) List(
	ctx context.Context,
	page pagination.Parameters,
	opts ...Query,
) (pagination.Result[*report.Report], error) {
	query := q.db.Report.Query()

	for _, fn := range opts {
		fn(query)
	}

	query.
		WithReportedBy().
		WithHandledBy().
		Order(
			ent.Desc(entreport.FieldUpdatedAt),
			ent.Desc(entreport.FieldID),
		)

	total, err := query.Count(ctx)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	query.
		Limit(page.Limit()).
		Offset(page.Offset())

	result, err := query.All(ctx)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	refs, err := dt.MapErr(result, report.Map)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	reports, err := q.hydrateRefs(ctx, refs)
	if err != nil {
		return pagination.Result[*report.Report]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return pagination.NewPageResult(page, total, reports), nil
}

func (q *Querier) hydrateRefs(ctx context.Context, refs report.ReportRefs) (report.Reports, error) {
	grouped := lo.GroupBy(refs, func(r *report.ReportRef) datagraph.Kind {
		return r.TargetRef.Kind
	})

	// post related reports

	pids := dt.Map(grouped[datagraph.KindPost], func(r *report.ReportRef) post.ID {
		return post.ID(r.TargetRef.ID)
	})
	tids := dt.Map(grouped[datagraph.KindThread], func(r *report.ReportRef) post.ID {
		return post.ID(r.TargetRef.ID)
	})
	rids := dt.Map(grouped[datagraph.KindReply], func(r *report.ReportRef) post.ID {
		return post.ID(r.TargetRef.ID)
	})
	allIDs := append(append(pids, tids...), rids...)
	posts, err := q.postSearcher.GetMany(ctx, allIDs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	pg := lo.KeyBy(posts, func(p *post.Post) post.ID { return p.ID })

	// profile related reports

	profileIDs := dt.Map(grouped[datagraph.KindProfile], func(r *report.ReportRef) account.AccountID {
		return account.AccountID(r.TargetRef.ID)
	})
	profiles, err := q.profileQuerier.GetMany(ctx, profileIDs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	profilesMap := lo.KeyBy(profiles, func(p *profile.Public) account.AccountID { return p.ID })

	// node related reports

	nodeIDs := dt.Map(grouped[datagraph.KindNode], func(r *report.ReportRef) library.NodeID {
		return library.NodeID(r.TargetRef.ID)
	})
	nodes, err := q.nodeQuerier.ProbeMany(ctx, nodeIDs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	nodesMap := lo.KeyBy(nodes, func(n *library.Node) xid.ID { return n.GetID() })

	reports := dt.Map(refs, func(r *report.ReportRef) *report.Report {
		var item datagraph.Item

		switch r.TargetRef.Kind {
		case datagraph.KindPost, datagraph.KindThread, datagraph.KindReply:
			p := pg[post.ID(r.TargetRef.ID)]
			if p != nil {
				item = p
			}

		case datagraph.KindProfile:
			p := profilesMap[account.AccountID(r.TargetRef.ID)]
			if p != nil {
				item = p
			}

		case datagraph.KindNode:
			n := nodesMap[r.TargetRef.ID]
			if n != nil {
				item = n
			}
		}

		return &report.Report{
			ID:             r.ID,
			TargetItemKind: r.TargetRef.Kind,
			TargetItemID:   xid.ID(r.TargetRef.ID),
			TargetItem:     item,
			ReportedBy:     r.ReportedBy,
			HandledBy:      r.HandledBy,
			Comment:        r.Comment,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		}
	})

	return reports, nil
}

func (q *Querier) Get(ctx context.Context, id report.ID) (*report.Report, error) {
	r, err := q.db.Report.Query().
		Where(entreport.ID(xid.ID(id))).
		WithReportedBy().
		WithHandledBy().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ref, err := report.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reports, err := q.hydrateRefs(ctx, report.ReportRefs{ref})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(reports) == 0 {
		return nil, fault.Wrap(fault.New("report not found"), fctx.With(ctx))
	}

	return reports[0], nil
}
