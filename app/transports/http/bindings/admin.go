package bindings

import (
	"context"
	"strconv"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/audit/audit_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/resources/timerange"
	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/services/admin/settings_manager"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/moderation/action_dispatcher"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Admin struct {
	accountQuery     *account_querier.Querier
	profileQuery     *profile_querier.Querier
	auditQuerier     *audit_querier.Querier
	as               account_suspension.Service
	settingsManager  *settings_manager.Manager
	akr              *access_key.Repository
	actionDispatcher *action_dispatcher.Service
}

func NewAdmin(
	accountQuery *account_querier.Querier,
	profileQuery *profile_querier.Querier,
	auditQuerier *audit_querier.Querier,
	as account_suspension.Service,
	settingsManager *settings_manager.Manager,
	akr *access_key.Repository,
	actionDispatcher *action_dispatcher.Service,
) Admin {
	return Admin{
		accountQuery:     accountQuery,
		profileQuery:     profileQuery,
		auditQuerier:     auditQuerier,
		as:               as,
		settingsManager:  settingsManager,
		akr:              akr,
		actionDispatcher: actionDispatcher,
	}
}

func (a *Admin) AdminSettingsGet(ctx context.Context, request openapi.AdminSettingsGetRequestObject) (openapi.AdminSettingsGetResponseObject, error) {
	settings, err := a.settingsManager.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminSettingsGet200JSONResponse{
		AdminSettingsGetOKJSONResponse: openapi.AdminSettingsGetOKJSONResponse(serialiseSettings(settings)),
	}, nil
}

func (a *Admin) AdminSettingsUpdate(ctx context.Context, request openapi.AdminSettingsUpdateRequestObject) (openapi.AdminSettingsUpdateResponseObject, error) {
	content, err := opt.MapErr(opt.NewPtr(request.Body.Content), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authMode, err := opt.MapErr(opt.NewPtr(request.Body.AuthenticationMode), deserialiseAuthMode)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var services opt.Optional[settings.ServiceSettings]
	if request.Body.Services != nil && request.Body.Services.Moderation != nil {
		moderation := request.Body.Services.Moderation
		services = opt.New(settings.ServiceSettings{
			Moderation: opt.New(settings.ModerationServiceSettings{
				ThreadBodyLengthMax: opt.NewPtr(moderation.ThreadBodyLengthMax),
				ReplyBodyLengthMax:  opt.NewPtr(moderation.ReplyBodyLengthMax),
				WordBlockList:       opt.NewPtr(moderation.WordBlockList),
				WordReportList:      opt.NewPtr(moderation.WordReportList),
			}),
		})
	}

	settings, err := a.settingsManager.Set(ctx, settings.Settings{
		Title:              opt.NewPtr(request.Body.Title),
		Description:        opt.NewPtr(request.Body.Description),
		Content:            content,
		AccentColour:       opt.NewPtr(request.Body.AccentColour),
		AuthenticationMode: authMode,
		Services:           services,
		Metadata:           opt.NewPtr((*map[string]any)(request.Body.Metadata)),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminSettingsUpdate200JSONResponse{
		AdminSettingsUpdateOKJSONResponse: openapi.AdminSettingsUpdateOKJSONResponse(serialiseSettings(settings)),
	}, nil
}

func (a *Admin) AuditEventList(ctx context.Context, request openapi.AuditEventListRequestObject) (openapi.AuditEventListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	page := opt.NewPtrMap(request.Params.Page, func(pq openapi.PaginationQuery) int {
		v, err := strconv.ParseInt(pq, 10, 32)
		if err != nil {
			return 1
		}
		return max(1, int(v))
	}).Or(1)

	params := pagination.NewPageParams(uint(page), 50)

	var eventTypes opt.Optional[[]audit.EventType]
	if request.Params.Types != nil && len(*request.Params.Types) > 0 {
		types, err := dt.MapErr(*request.Params.Types, func(t openapi.AuditEventType) (audit.EventType, error) {
			return audit.NewEventType(string(t))
		})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		eventTypes = opt.New(types)
	}

	var filterTimeRange opt.Optional[audit_querier.TimeRange]
	if request.Params.Range != nil {
		tr, err := timerange.Parse(*request.Params.Range)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		if tr.Start.Ok() || tr.End.Ok() {
			var start, end time.Time
			tr.Start.Call(func(t time.Time) { start = t })
			tr.End.Call(func(t time.Time) { end = t })

			filterTimeRange = opt.New(audit_querier.TimeRange{
				Start: start,
				End:   end,
			})
		}
	}

	result, err := a.auditQuerier.List(ctx, params, audit_querier.Filter{
		Types:     eventTypes,
		TimeRange: filterTimeRange,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	events := dt.Map(result.Items, serialiseAuditEvent)
	eventList := openapi.AuditEventList(events)

	return openapi.AuditEventList200JSONResponse{
		AuditEventListOKJSONResponse: openapi.AuditEventListOKJSONResponse{
			CurrentPage: result.CurrentPage,
			Events:      &eventList,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
		},
	}, nil
}

func (a *Admin) AuditEventGet(ctx context.Context, request openapi.AuditEventGetRequestObject) (openapi.AuditEventGetResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	id := deserialiseID(request.AuditEventId)

	auditLog, err := a.auditQuerier.Get(ctx, audit.AuditLogID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	event := serialiseAuditEvent(auditLog)

	return openapi.AuditEventGet200JSONResponse{
		AuditEventGetOKJSONResponse: openapi.AuditEventGetOKJSONResponse(event),
	}, nil
}

func (a *Admin) ModerationActionCreate(ctx context.Context, request openapi.ModerationActionCreateRequestObject) (openapi.ModerationActionCreateResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	enactedBy, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	discriminator, err := request.Body.Discriminator()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	switch discriminator {
	case "purge_account":
		purgeBody, err := request.Body.AsModerationActionCreatePurgeAccount()
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		accountID := account.AccountID(deserialiseID(purgeBody.AccountId))

		contentTypes, err := dt.MapErr(purgeBody.Include, func(ct openapi.ModerationActionPurgeAccountContentType) (action_dispatcher.ContentType, error) {
			return action_dispatcher.NewContentType(string(ct))
		})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		auditLog, err := a.actionDispatcher.PurgeAccountContent(
			ctx,
			accountID,
			opt.New(enactedBy),
			contentTypes,
		)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		props := serialiseAuditEventProps(auditLog)

		return openapi.ModerationActionCreate201JSONResponse{
			AuditEventCreatedOKJSONResponse: openapi.AuditEventCreatedOKJSONResponse(props),
		}, nil

	default:
		return nil, fault.Wrap(
			fault.New("unknown moderation action: "+discriminator),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}
}

func (i *Admin) AdminAccountBanCreate(ctx context.Context, request openapi.AdminAccountBanCreateRequestObject) (openapi.AdminAccountBanCreateResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	acc, err = i.as.Suspend(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccountBanCreate200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
		},
	}, nil
}

func (i *Admin) AdminAccountBanRemove(ctx context.Context, request openapi.AdminAccountBanRemoveRequestObject) (openapi.AdminAccountBanRemoveResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	acc, err = i.as.Reinstate(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccountBanRemove200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
		},
	}, nil
}

func (i *Admin) AdminAccessKeyList(ctx context.Context, request openapi.AdminAccessKeyListRequestObject) (openapi.AdminAccessKeyListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := i.akr.ListAllAsAdmin(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccessKeyList200JSONResponse{
		AdminAccessKeyListOKJSONResponse: openapi.AdminAccessKeyListOKJSONResponse{
			Keys: serialiseOwnedAccessKeyList(list),
		},
	}, nil
}

func (i *Admin) AdminAccessKeyDelete(ctx context.Context, request openapi.AdminAccessKeyDeleteRequestObject) (openapi.AdminAccessKeyDeleteResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err := i.akr.RevokeAsAdmin(ctx, deserialiseID(request.AccessKeyId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NoContentResponse{}, nil
}

func serialiseSettings(in *settings.Settings) openapi.AdminSettingsProps {
	return openapi.AdminSettingsProps{
		AccentColour:       in.AccentColour.OrZero(),
		Description:        in.Description.OrZero(),
		Content:            in.Content.OrZero().HTML(),
		Title:              in.Title.OrZero(),
		AuthenticationMode: openapi.AuthMode(in.AuthenticationMode.Or(authentication.ModeHandle).String()),
		Services:           opt.Map(in.Services, serialiseServiceSettings).Ptr(),
		Metadata:           (*openapi.Metadata)(in.Metadata.Ptr()),
	}
}

func serialiseServiceSettings(in settings.ServiceSettings) openapi.AdminSettingsServiceProps {
	return openapi.AdminSettingsServiceProps{
		Moderation: opt.Map(in.Moderation, serialiseModerationSettings).Ptr(),
	}
}

func serialiseModerationSettings(in settings.ModerationServiceSettings) openapi.ModerationServiceSettings {
	return openapi.ModerationServiceSettings{
		ThreadBodyLengthMax: in.ThreadBodyLengthMax.Ptr(),
		ReplyBodyLengthMax:  in.ReplyBodyLengthMax.Ptr(),
		WordBlockList:       in.WordBlockList.Ptr(),
		WordReportList:      in.WordReportList.Ptr(),
	}
}

func serialiseOwnedAccessKey(in *authentication.Authentication) openapi.OwnedAccessKey {
	return openapi.OwnedAccessKey{
		Id:        in.ID.String(),
		CreatedAt: in.Account.CreatedAt,
		ExpiresAt: in.Expires.Ptr(),
		Enabled:   !in.Disabled,
		Name:      in.Name.Or("Unnamed"),
		CreatedBy: serialiseProfileReferenceFromAccount(in.Account),
	}
}

func serialiseOwnedAccessKeyList(in []*authentication.Authentication) openapi.OwnedAccessKeyList {
	return dt.Map(in, serialiseOwnedAccessKey)
}

func serialiseAuditEvent(in *audit.AuditLog) openapi.AuditEvent {
	out := openapi.AuditEvent{
		Id:        openapi.Identifier(in.ID.String()),
		Type:      openapi.AuditEventType(in.Type.String()),
		Timestamp: in.CreatedAt,
	}

	in.EnactedBy.Call(func(acc account.Account) {
		ref := serialiseProfileReferenceFromAccount(acc)
		out.EnactedBy = &ref
	})

	var err error

	switch in.Type {
	case audit.EventTypeThreadDeleted:
		err = out.FromAuditEventThreadDeleted(openapi.AuditEventThreadDeleted{
			Type:     openapi.ThreadDeleted,
			ThreadId: openapi.Identifier(in.Target.OrZero().ID.String()),
		})

	case audit.EventTypeThreadReplyDeleted:
		err = out.FromAuditEventThreadReplyDeleted(openapi.AuditEventThreadReplyDeleted{
			Type:    openapi.ThreadReplyDeleted,
			ReplyId: openapi.Identifier(in.Target.OrZero().ID.String()),
		})

	case audit.EventTypeAccountSuspended:
		accountID := in.Metadata["account_id"].(string)
		err = out.FromAuditEventAccountSuspended(openapi.AuditEventAccountSuspended{
			Type:      openapi.AccountSuspended,
			AccountId: openapi.Identifier(accountID),
		})

	case audit.EventTypeAccountUnsuspended:
		accountID := in.Metadata["account_id"].(string)
		err = out.FromAuditEventAccountUnsuspended(openapi.AuditEventAccountUnsuspended{
			Type:      openapi.AccountUnsuspended,
			AccountId: openapi.Identifier(accountID),
		})

	case audit.EventTypeAccountContentPurged:
		accountID := in.Metadata["account_id"].(string)
		var included []openapi.ModerationActionPurgeAccountContentType
		if inc, ok := in.Metadata["included"].([]interface{}); ok {
			for _, item := range inc {
				if str, ok := item.(string); ok {
					included = append(included, openapi.ModerationActionPurgeAccountContentType(str))
				}
			}
		}

		err = out.FromAuditEventAccountContentPurged(openapi.AuditEventAccountContentPurged{
			Type:      openapi.AccountContentPurged,
			AccountId: openapi.Identifier(accountID),
			Included:  &included,
		})
	}

	if err != nil {
		panic(err)
	}

	return out
}

func serialiseAuditEventProps(in *audit.AuditLog) openapi.AuditEventProps {
	var enactedBy *openapi.ProfileReference
	in.EnactedBy.Call(func(acc account.Account) {
		ref := serialiseProfileReferenceFromAccount(acc)
		enactedBy = &ref
	})

	return openapi.AuditEventProps{
		Id:        openapi.Identifier(in.ID.String()),
		Type:      openapi.AuditEventType(in.Type.String()),
		Timestamp: in.CreatedAt,
		EnactedBy: enactedBy,
	}
}
