package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/report/member_report"
	"github.com/Southclaws/storyden/app/services/report/report_manager"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Reports struct {
	memberReportMgr *member_report.Manager
	reportMgr       *report_manager.Manager
}

func NewReports(
	memberReportMgr *member_report.Manager,
	reportMgr *report_manager.Manager,
) Reports {
	return Reports{
		memberReportMgr: memberReportMgr,
		reportMgr:       reportMgr,
	}
}

func (h *Reports) ReportCreate(ctx context.Context, request openapi.ReportCreateRequestObject) (openapi.ReportCreateResponseObject, error) {
	targetID, err := xid.FromString(request.Body.TargetId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	targetKind, err := datagraph.NewKind(string(request.Body.TargetKind))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := h.memberReportMgr.Submit(
		ctx,
		targetID,
		targetKind,
		opt.NewPtr(request.Body.Comment),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ReportCreate200JSONResponse{
		ReportCreateOKJSONResponse: openapi.ReportCreateOKJSONResponse(serialiseReport(r)),
	}, nil
}

func (h *Reports) ReportList(ctx context.Context, request openapi.ReportListRequestObject) (openapi.ReportListResponseObject, error) {
	roles := session.GetRoles(ctx)
	permissions := roles.Permissions()

	page := deserialisePageParams(request.Params.Page, 10)

	canManageReports := permissions.HasAny(rbac.PermissionManageReports, rbac.PermissionAdministrator)

	var result pagination.Result[*report.Report]
	var err error

	if canManageReports {
		opts := report_manager.ListOpts{}

		if request.Params.Status != nil {
			statuses := []report.Status{}
			status, err := report.NewStatus(string(*request.Params.Status))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			statuses = append(statuses, status)
			opts.Status = opt.New(statuses)
		}

		if request.Params.Kind != nil {
			kinds := []datagraph.Kind{}
			kind, err := datagraph.NewKind(*request.Params.Kind)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			kinds = append(kinds, kind)
			opts.Kind = opt.New(kinds)
		}

		result, err = h.reportMgr.List(ctx, page, opts)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		result, err = h.memberReportMgr.List(ctx, page)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return openapi.ReportList200JSONResponse{
		ReportListOKJSONResponse: openapi.ReportListOKJSONResponse(serialiseReportList(result)),
	}, nil
}

func (h *Reports) ReportUpdate(ctx context.Context, request openapi.ReportUpdateRequestObject) (openapi.ReportUpdateResponseObject, error) {
	roles := session.GetRoles(ctx)
	permissions := roles.Permissions()

	reportID, err := xid.FromString(request.ReportId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	canManageReports := permissions.HasAny(rbac.PermissionManageReports, rbac.PermissionAdministrator)

	var r *report.Report

	if canManageReports {
		updateOpts := report_manager.UpdateOpts{}

		if request.Body.Status != nil {
			status, err := report.NewStatus(string(*request.Body.Status))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			updateOpts.Status = opt.New(status)
		}

		if request.Body.HandledBy != nil {
			handlerID, err := xid.FromString(*request.Body.HandledBy)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			updateOpts.HandledBy = opt.New(account.AccountID(handlerID))
		}

		r, err = h.reportMgr.Update(ctx, report.ID(reportID), updateOpts)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		if request.Body.Status != nil {
			status, err := report.NewStatus(string(*request.Body.Status))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			if status != report.StatusResolved {
				return nil, fault.Wrap(
					fault.New("members can only resolve their own reports"),
					fctx.With(ctx),
					ftag.With(ftag.PermissionDenied),
				)
			}
		}

		r, err = h.memberReportMgr.Resolve(ctx, report.ID(reportID))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return openapi.ReportUpdate200JSONResponse{
		ReportUpdateOKJSONResponse: openapi.ReportUpdateOKJSONResponse(serialiseReport(r)),
	}, nil
}

func serialiseReport(in *report.Report) openapi.Report {
	item := opt.Map(opt.NewSafe(in.TargetItem, in.TargetItem != nil), serialiseDatagraphItem)

	return openapi.Report{
		Id:         in.ID.String(),
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
		TargetId:   in.TargetItemID.String(),
		TargetKind: openapi.DatagraphItemKind(in.TargetItemKind.String()),
		Item:       item.Ptr(),
		ReportedBy: opt.Map(in.ReportedBy, serialiseProfileReferenceFromAccount).Ptr(),
		HandledBy:  opt.Map(in.HandledBy, serialiseProfileReferenceFromAccount).Ptr(),
		Comment:    in.Comment.Ptr(),
		Status:     openapi.ReportStatus(in.Status.String()),
	}
}

func serialiseReportList(in pagination.Result[*report.Report]) openapi.ReportListResult {
	items := dt.Map(in.Items, serialiseReport)

	return openapi.ReportListResult{
		Reports:     openapi.ReportList(items),
		CurrentPage: in.CurrentPage,
		NextPage:    in.NextPage.Ptr(),
		Results:     in.Results,
		TotalPages:  in.TotalPages,
	}
}
