package bindings

import (
	"context"
	"net/url"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_badge"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/account/account_auth"
	"github.com/Southclaws/storyden/app/services/account/account_email"
	"github.com/Southclaws/storyden/app/services/account/account_manage"
	"github.com/Southclaws/storyden/app/services/account/account_update"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Accounts struct {
	profile_cache *profile_cache.Cache
	avatarService avatar.Service
	authManager   *authentication.Manager
	accountQuery  *account_querier.Querier
	profileQuery  *profile_querier.Querier
	accountUpdate *account_update.Updater
	accountAuth   *account_auth.Manager
	accountEmail  *account_email.Manager
	accountManage *account_manage.Manager
	roleAssign    *role_assign.Assignment
	roleBadge     *role_badge.Writer
	webAddress    url.URL
}

func NewAccounts(
	cfg config.Config,
	profile_cache *profile_cache.Cache,
	avatarService avatar.Service,
	authManager *authentication.Manager,
	accountQuery *account_querier.Querier,
	profileQuery *profile_querier.Querier,
	accountUpdate *account_update.Updater,
	accountAuth *account_auth.Manager,
	accountEmail *account_email.Manager,
	accountManage *account_manage.Manager,
	roleAssign *role_assign.Assignment,
	roleBadge *role_badge.Writer,
) Accounts {
	return Accounts{
		profile_cache: profile_cache,
		avatarService: avatarService,
		authManager:   authManager,
		accountQuery:  accountQuery,
		profileQuery:  profileQuery,
		accountUpdate: accountUpdate,
		accountAuth:   accountAuth,
		accountEmail:  accountEmail,
		accountManage: accountManage,
		roleAssign:    roleAssign,
		roleBadge:     roleBadge,
		webAddress:    cfg.PublicWebAddress,
	}
}

var (
	ErrSelfAdminRoleChange = fault.New("cannot change own admin role", ftag.With(ftag.InvalidArgument), fmsg.WithDesc("admin role", "You cannot change your own admin role."))
	ErrEveryoneRole        = fault.New("cannot change default role", ftag.With(ftag.InvalidArgument), fmsg.WithDesc("default role", "You cannot change the default role."))
)

const accountGetCacheControl = "private, no-cache"

func (i *Accounts) AccountGet(ctx context.Context, request openapi.AccountGetRequestObject) (openapi.AccountGetResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	etag, notModified := i.profile_cache.Check(ctx, reqinfo.GetCacheQuery(ctx), xid.ID(accountID))
	if notModified {
		return openapi.AccountGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: accountGetCacheControl,
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		}, nil
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if etag == nil {
		i.profile_cache.Store(ctx, xid.ID(accountID), acc.UpdatedAt)
		etag = cachecontrol.NewETag(acc.UpdatedAt)
	}

	return openapi.AccountGet200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
			Headers: openapi.AccountGetOKResponseHeaders{
				CacheControl: accountGetCacheControl,
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		},
	}, nil
}

func (i *Accounts) AccountView(ctx context.Context, request openapi.AccountViewRequestObject) (openapi.AccountViewResponseObject, error) {
	targetID, err := xid.FromString(request.AccountId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument), fmsg.WithDesc("invalid account ID", "The account ID provided is invalid."))
	}

	acc, err := i.accountManage.GetByID(ctx, account.AccountID(targetID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountView200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
			Headers: openapi.AccountGetOKResponseHeaders{
				CacheControl: "private, no-cache",
				LastModified: acc.UpdatedAt.Format(time.RFC1123),
			},
		},
	}, nil
}

func (i *Accounts) AccountUpdate(ctx context.Context, request openapi.AccountUpdateRequestObject) (openapi.AccountUpdateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	links, err := opt.MapErr(opt.NewPtr(request.Body.Links), deserialiseExternalLinkList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountUpdate.Update(ctx, accountID, account_update.Partial{
		Handle:    opt.NewPtrMap(request.Body.Handle, func(i openapi.AccountHandle) string { return string(i) }),
		Name:      opt.NewPtr(request.Body.Name),
		Bio:       opt.NewPtr(request.Body.Bio),
		Links:     links,
		Meta:      opt.NewPtr((*map[string]any)(request.Body.Meta)),
		Interests: opt.NewPtrMap(request.Body.Interests, tagsIDs),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountUpdate200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func deserialiseExternalLinkList(i openapi.ProfileExternalLinkList) ([]account.ExternalLink, error) {
	return dt.MapErr(i, deserialiseExternalLink)
}

func deserialiseExternalLink(l openapi.ProfileExternalLink) (account.ExternalLink, error) {
	u, err := url.Parse(string(l.Url))
	if err != nil {
		return account.ExternalLink{}, err
	}

	return account.ExternalLink{
		Text: string(l.Text),
		URL:  *u,
	}, nil
}

func (i *Accounts) AccountAuthProviderList(ctx context.Context, request openapi.AccountAuthProviderListRequestObject) (openapi.AccountAuthProviderListResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := i.authManager.GetProviderList(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authmethods, err := i.accountAuth.GetAuthMethods(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	available, err := dt.MapErr(providers, serialiseAuthProvider(buildRedirectURL(i.webAddress)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	active, err := dt.MapErr(authmethods, serialiseAuthMethod(i.webAddress))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAuthProviderList200JSONResponse{
		AccountAuthProviderListOKJSONResponse: openapi.AccountAuthProviderListOKJSONResponse{
			Available: available,
			Active:    active,
		},
	}, nil
}

func (i *Accounts) AccountAuthMethodDelete(ctx context.Context, request openapi.AccountAuthMethodDeleteRequestObject) (openapi.AccountAuthMethodDeleteResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	id, err := xid.FromString(request.AuthMethodId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	err = i.accountAuth.DeleteAuthMethod(ctx, accountID, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := i.authManager.GetProviderList(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authmethods, err := i.accountAuth.GetAuthMethods(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	available, err := dt.MapErr(providers, serialiseAuthProvider(buildRedirectURL(i.webAddress)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	active, err := dt.MapErr(authmethods, serialiseAuthMethod(i.webAddress))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAuthMethodDelete200JSONResponse{
		AccountAuthProviderListOKJSONResponse: openapi.AccountAuthProviderListOKJSONResponse{
			Available: available,
			Active:    active,
		},
	}, nil
}

func (i *Accounts) AccountGetAvatar(ctx context.Context, request openapi.AccountGetAvatarRequestObject) (openapi.AccountGetAvatarResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, size, err := i.avatarService.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountGetAvatar200ImagepngResponse{
		AccountGetAvatarImagepngResponse: openapi.AccountGetAvatarImagepngResponse{
			Body:          r,
			ContentLength: size,
			Headers: openapi.AccountGetAvatarResponseHeaders{
				CacheControl: "public, max-age=3600",
			},
		},
	}, nil
}

func (i *Accounts) AccountSetAvatar(ctx context.Context, request openapi.AccountSetAvatarRequestObject) (openapi.AccountSetAvatarResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := i.avatarService.Set(ctx, accountID, request.Body, int64(request.Params.ContentLength)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountSetAvatar200Response{}, nil
}

func (h *Accounts) AccountRemoveRole(ctx context.Context, request openapi.AccountRemoveRoleRequestObject) (openapi.AccountRemoveRoleResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, found, err := h.accountQuery.LookupByHandle(ctx, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		return nil, fault.New("account not found", fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	roleID := role.RoleID(openapi.ParseID(request.RoleId))

	if accountID == acc.ID {
		// Updating self roles
		if roleID == role.DefaultRoleAdminID {
			return nil, fault.Wrap(ErrSelfAdminRoleChange, fctx.With(ctx))
		}
	}
	if roleID == role.DefaultRoleMemberID {
		return nil, fault.Wrap(ErrEveryoneRole, fctx.With(ctx))
	}

	acc, err = h.roleAssign.UpdateRoles(ctx, acc.ID, role_assign.Remove(roleID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountRemoveRole200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (h *Accounts) AccountAddRole(ctx context.Context, request openapi.AccountAddRoleRequestObject) (openapi.AccountAddRoleResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, found, err := h.accountQuery.LookupByHandle(ctx, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !found {
		return nil, fault.New("account not found", fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	roleID := role.RoleID(openapi.ParseID(request.RoleId))

	if accountID == acc.ID {
		// Updating self roles
		if roleID == role.DefaultRoleAdminID {
			return nil, fault.Wrap(ErrSelfAdminRoleChange, fctx.With(ctx))
		}
	}
	if roleID == role.DefaultRoleMemberID {
		return nil, fault.Wrap(ErrEveryoneRole, fctx.With(ctx))
	}

	acc, err = h.roleAssign.UpdateRoles(ctx, acc.ID, role_assign.Add(roleID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAddRole200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (h *Accounts) AccountRoleSetBadge(ctx context.Context, request openapi.AccountRoleSetBadgeRequestObject) (openapi.AccountRoleSetBadgeResponseObject, error) {
	roleID := role.RoleID(openapi.ParseID(request.RoleId))

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := h.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	target, found, err := h.accountQuery.LookupByHandle(ctx, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !found {
		return nil, fault.New("account not found", fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	updatingSelf := accountID == target.ID

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if updatingSelf {
			return nil
		}
		return errNotAuthorised
	}, rbac.PermissionManageRoles); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err = h.roleBadge.Update(ctx, target.ID, roleID, true)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountRoleSetBadge200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (h *Accounts) AccountRoleRemoveBadge(ctx context.Context, request openapi.AccountRoleRemoveBadgeRequestObject) (openapi.AccountRoleRemoveBadgeResponseObject, error) {
	roleID := role.RoleID(openapi.ParseID(request.RoleId))

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := h.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	target, found, err := h.accountQuery.LookupByHandle(ctx, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !found {
		return nil, fault.New("account not found", fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	updatingSelf := accountID == target.ID

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if updatingSelf {
			return nil
		}
		return errNotAuthorised
	}, rbac.PermissionManageRoles); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err = h.roleBadge.Update(ctx, target.ID, roleID, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountRoleRemoveBadge200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func serialiseAuthMethod(webAddress url.URL) func(in *account_auth.AuthMethod) (openapi.AccountAuthMethod, error) {
	return func(in *account_auth.AuthMethod) (openapi.AccountAuthMethod, error) {
		p, err := serialiseAuthProvider(buildRedirectURL(webAddress))(in.Provider)
		if err != nil {
			return openapi.AccountAuthMethod{}, fault.Wrap(err)
		}

		return openapi.AccountAuthMethod{
			Id:         in.Instance.ID.String(),
			CreatedAt:  in.Instance.Created,
			Name:       in.Instance.Name.Or("Unknown"),
			Identifier: in.Instance.Identifier,
			Provider:   p,
		}, nil
	}
}
