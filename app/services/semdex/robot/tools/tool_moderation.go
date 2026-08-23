package tools

import (
	"context"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/report/member_report"
	"github.com/Southclaws/storyden/app/services/report/report_manager"
	"github.com/Southclaws/storyden/lib/mcp"
)

type moderationTools struct {
	memberReportMgr *member_report.Manager
	reportMgr       *report_manager.Manager
	suspensions     account_suspension.Service
}

func newModerationTools(
	registry *Registry,
	memberReportMgr *member_report.Manager,
	reportMgr *report_manager.Manager,
	suspensions account_suspension.Service,
) *moderationTools {
	t := &moderationTools{
		memberReportMgr: memberReportMgr,
		reportMgr:       reportMgr,
		suspensions:     suspensions,
	}

	registry.Register(t.newReportCreateTool())
	registry.Register(t.newReportListTool())
	registry.Register(t.newReportGetTool())
	registry.Register(t.newReportUpdateTool())
	registry.Register(t.newMemberSuspendTool())
	registry.Register(t.newMemberReinstateTool())

	return t
}

func (mt *moderationTools) newReportCreateTool() *Tool {
	toolDef := mcp.GetReportCreateTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolReportCreateInput) (*mcp.ToolReportCreateOutput, error) {
				return mt.ExecuteReportCreate(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteReportCreate),
	}
}

func (mt *moderationTools) ExecuteReportCreate(ctx context.Context, args mcp.ToolReportCreateInput) (*mcp.ToolReportCreateOutput, error) {
	targetID, err := xid.FromString(args.TargetId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	targetKind, err := datagraph.NewKind(string(args.TargetKind))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	rep, err := mt.memberReportMgr.Submit(ctx, targetID, targetKind, opt.New(args.Comment))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &mcp.ToolReportCreateOutput{Report: mapModerationReport(rep)}, nil
}

func (mt *moderationTools) newReportListTool() *Tool {
	toolDef := mcp.GetReportListTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolReportListInput) (*mcp.ToolReportListOutput, error) {
				return mt.ExecuteReportList(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteReportList),
	}
}

func (mt *moderationTools) ExecuteReportList(ctx context.Context, args mcp.ToolReportListInput) (*mcp.ToolReportListOutput, error) {
	page, pageSize := 1, 20
	if args.Page != nil {
		page = *args.Page
	}
	if args.PageSize != nil {
		pageSize = *args.PageSize
	}

	listOpts := report_manager.ListOpts{}
	if args.Status != nil {
		status, err := report.NewStatus(string(*args.Status))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		listOpts.Status = opt.New([]report.Status{status})
	}
	if args.TargetKind != nil {
		kind, err := datagraph.NewKind(string(*args.TargetKind))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		listOpts.Kind = opt.New([]datagraph.Kind{kind})
	}

	result, err := mt.reportMgr.List(ctx, pagination.NewPageParams(uint(page), uint(pageSize)), listOpts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reports := make([]mcp.ModerationToolReportYaml, 0, len(result.Items))
	for _, rep := range result.Items {
		reports = append(reports, mapModerationReport(rep))
	}

	return &mcp.ToolReportListOutput{
		Reports:  reports,
		Results:  result.Results,
		Page:     result.CurrentPage,
		NextPage: result.NextPage.Ptr(),
	}, nil
}

func (mt *moderationTools) newReportGetTool() *Tool {
	toolDef := mcp.GetReportGetTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolReportGetInput) (*mcp.ToolReportGetOutput, error) {
				return mt.ExecuteReportGet(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteReportGet),
	}
}

func (mt *moderationTools) ExecuteReportGet(ctx context.Context, args mcp.ToolReportGetInput) (*mcp.ToolReportGetOutput, error) {
	id, err := xid.FromString(args.ReportId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	rep, err := mt.reportMgr.Get(ctx, report.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &mcp.ToolReportGetOutput{Report: mapModerationReport(rep)}, nil
}

func (mt *moderationTools) newReportUpdateTool() *Tool {
	toolDef := mcp.GetReportUpdateTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolReportUpdateInput) (*mcp.ToolReportUpdateOutput, error) {
				return mt.ExecuteReportUpdate(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteReportUpdate),
	}
}

func (mt *moderationTools) ExecuteReportUpdate(ctx context.Context, args mcp.ToolReportUpdateInput) (*mcp.ToolReportUpdateOutput, error) {
	id, err := xid.FromString(args.ReportId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	status, err := report.NewStatus(string(args.Status))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	rep, err := mt.reportMgr.Update(ctx, report.ID(id), report_manager.UpdateOpts{Status: opt.New(status)})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &mcp.ToolReportUpdateOutput{Report: mapModerationReport(rep)}, nil
}

func (mt *moderationTools) newMemberSuspendTool() *Tool {
	toolDef := mcp.GetMemberSuspendTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolMemberSuspendInput) (*mcp.ToolMemberSuspendOutput, error) {
				return mt.ExecuteMemberSuspend(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteMemberSuspend),
	}
}

func (mt *moderationTools) ExecuteMemberSuspend(ctx context.Context, args mcp.ToolMemberSuspendInput) (*mcp.ToolMemberSuspendOutput, error) {
	id, err := xid.FromString(args.AccountId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	actingAccountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if account.AccountID(id) == actingAccountID {
		return nil, fault.Wrap(fault.New("a Robot cannot suspend its own account"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	acc, err := mt.suspensions.Suspend(ctx, account.AccountID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &mcp.ToolMemberSuspendOutput{AccountId: acc.ID.String(), Handle: acc.Handle, Suspended: true}, nil
}

func (mt *moderationTools) newMemberReinstateTool() *Tool {
	toolDef := mcp.GetMemberReinstateTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema}, func(ctx adkagent.Context, args mcp.ToolMemberReinstateInput) (*mcp.ToolMemberReinstateOutput, error) {
				return mt.ExecuteMemberReinstate(ctx, args)
			})
		},
		Handler: makeHandler(mt.ExecuteMemberReinstate),
	}
}

func (mt *moderationTools) ExecuteMemberReinstate(ctx context.Context, args mcp.ToolMemberReinstateInput) (*mcp.ToolMemberReinstateOutput, error) {
	id, err := xid.FromString(args.AccountId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	acc, err := mt.suspensions.Reinstate(ctx, account.AccountID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &mcp.ToolMemberReinstateOutput{AccountId: acc.ID.String(), Handle: acc.Handle, Suspended: false}, nil
}

func mapModerationReport(rep *report.Report) mcp.ModerationToolReportYaml {
	out := mcp.ModerationToolReportYaml{
		Id:         rep.ID.String(),
		TargetId:   rep.TargetItemID.String(),
		TargetKind: mcp.ModerationToolReportYamlTargetKind(rep.TargetItemKind.String()),
		Status:     mcp.ModerationToolReportYamlStatus(rep.Status.String()),
		Comment:    rep.Comment.Ptr(),
		CreatedAt:  rep.CreatedAt,
		UpdatedAt:  rep.UpdatedAt,
	}
	if reportedBy, ok := rep.ReportedBy.Get(); ok {
		id := reportedBy.ID.String()
		out.ReportedById = &id
	}
	if handledBy, ok := rep.HandledBy.Get(); ok {
		id := handledBy.ID.String()
		out.HandledById = &id
	}
	return out
}
