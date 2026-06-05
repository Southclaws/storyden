package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	oauthservice "github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type OAuth struct {
	oauth *oauthservice.Service
}

func NewOAuth(oauth *oauthservice.Service) OAuth {
	return OAuth{oauth: oauth}
}

type OAuthDiscoveryResponse struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	DeviceAuthorizationEndpoint      string   `json:"device_authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	JWKSURI                          string   `json:"jwks_uri"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	GrantTypesSupported              []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported    []string `json:"code_challenge_methods_supported"`
	ScopesSupported                  []string `json:"scopes_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

func (o OAuth) OAuthDiscovery(context.Context) OAuthDiscoveryResponse {
	discovery := o.oauth.Discovery()

	return OAuthDiscoveryResponse{
		Issuer:                           discovery.Issuer,
		AuthorizationEndpoint:            discovery.AuthorizationEndpoint,
		DeviceAuthorizationEndpoint:      discovery.DeviceAuthorizationEndpoint,
		TokenEndpoint:                    discovery.TokenEndpoint,
		UserinfoEndpoint:                 discovery.UserinfoEndpoint,
		JWKSURI:                          discovery.JWKSURI,
		ResponseTypesSupported:           discovery.ResponseTypesSupported,
		GrantTypesSupported:              discovery.GrantTypesSupported,
		CodeChallengeMethodsSupported:    discovery.CodeChallengeMethodsSupported,
		ScopesSupported:                  discovery.ScopesSupported,
		SubjectTypesSupported:            discovery.SubjectTypesSupported,
		IDTokenSigningAlgValuesSupported: discovery.IDTokenSigningAlgValuesSupported,
	}
}

func (o OAuth) OAuthJWKS(ctx context.Context, _ openapi.OAuthJWKSRequestObject) (openapi.OAuthJWKSResponseObject, error) {
	if !o.oauth.Enabled() {
		return nil, oauthDisabledError(ctx)
	}

	return openapi.OAuthJWKS200JSONResponse{
		OAuthJWKSOKJSONResponse: openapi.OAuthJWKSOKJSONResponse(openapi.OAuthJWKS{
			Keys: mapOAuthJWKs(o.oauth.JWKS()),
		}),
	}, nil
}

func oauthDisabledError(ctx context.Context) error {
	message := "OAuth and OpenID Connect are not enabled on this Storyden instance."
	suggested := "Ask the instance administrator to configure OAuth before using OAuth clients."

	ctx = fctx.WithMeta(ctx,
		"code", "oauth_disabled",
		"suggested", suggested,
	)

	return fault.New("oauth_disabled",
		fctx.With(ctx),
		fmsg.WithDesc("oauth_disabled", message),
		ftag.With(ftag.NotFound),
	)
}

func (o OAuth) OAuthDeviceAuthorisation(ctx context.Context, req openapi.OAuthDeviceAuthorisationRequestObject) (openapi.OAuthDeviceAuthorisationResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthDeviceAuthorisation400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{
				Error: "invalid_request",
			}),
		}, nil
	}

	result, oauthErr, err := o.oauth.StartDeviceAuthorization(ctx, req.Body.ClientId, opt.NewPtr(req.Body.Scope))
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthDeviceAuthorisation400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{
				Error: oauthErr.Code,
			}),
		}, nil
	}

	return openapi.OAuthDeviceAuthorisation200JSONResponse{
		OAuthDeviceAuthorisationOKJSONResponse: openapi.OAuthDeviceAuthorisationOKJSONResponse(openapi.OAuthDeviceAuthorisation{
			DeviceCode:              &result.DeviceCode,
			UserCode:                &result.UserCode,
			VerificationUri:         &result.VerificationURI,
			VerificationUriComplete: &result.VerificationURIComplete,
			ExpiresIn:               &result.ExpiresIn,
			Interval:                &result.Interval,
		}),
	}, nil
}

func (o OAuth) OAuthDeviceConsent(ctx context.Context, req openapi.OAuthDeviceConsentRequestObject) (openapi.OAuthDeviceConsentResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthDeviceConsent401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	userCode := ""
	if req.Params.UserCode != nil {
		userCode = string(*req.Params.UserCode)
	}

	consent, oauthErr, err := o.oauth.GetDeviceConsent(ctx, account.AccountID(acc), permissions, userCode)
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthDeviceConsent400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	return openapi.OAuthDeviceConsent200JSONResponse{
		OAuthDeviceConsentOKJSONResponse: openapi.OAuthDeviceConsentOKJSONResponse(openapi.OAuthDeviceConsent{
			UserCode:                consent.UserCode,
			ClientId:                consent.ClientID,
			ClientName:              consent.ClientName,
			ExpiresAt:               consent.ExpiresAt,
			RequestedScopes:         consent.RequestedScopes,
			GrantedScopes:           consent.GrantedScopes,
			InheritsUserPermissions: consent.InheritsUserPermissions,
		}),
	}, nil
}

func (o OAuth) OAuthDeviceConsentSubmit(ctx context.Context, req openapi.OAuthDeviceConsentSubmitRequestObject) (openapi.OAuthDeviceConsentSubmitResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthDeviceConsentSubmit400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: "invalid_request"}),
		}, nil
	}

	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthDeviceConsentSubmit401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	oauthErr := o.oauth.ApproveDeviceAuthorization(ctx, account.AccountID(acc), permissions, req.Body.UserCode, req.Body.Decision == openapi.OAuthDeviceDecisionApprove)
	if oauthErr != nil {
		return openapi.OAuthDeviceConsentSubmit400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	status := openapi.OAuthDeviceConsentResultStatusDenied
	if req.Body.Decision == openapi.OAuthDeviceDecisionApprove {
		status = openapi.OAuthDeviceConsentResultStatusApproved
	}

	return openapi.OAuthDeviceConsentSubmit200JSONResponse{
		OAuthDeviceConsentSubmitOKJSONResponse: openapi.OAuthDeviceConsentSubmitOKJSONResponse(openapi.OAuthDeviceConsentResult{
			Status: status,
		}),
	}, nil
}

func (o OAuth) OAuthAuthorise(ctx context.Context, req openapi.OAuthAuthoriseRequestObject) (openapi.OAuthAuthoriseResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthAuthorise302Response{
			Headers: openapi.OAuthAuthoriseFoundResponseHeaders{Location: "/login"},
		}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	result, oauthErr, err := o.oauth.Authorise(ctx, oauthservice.AuthoriseRequest{
		ResponseType:        string(req.Params.ResponseType),
		ClientID:            string(req.Params.ClientId),
		RedirectURI:         string(req.Params.RedirectUri),
		Scope:               opt.NewPtr(req.Params.Scope),
		State:               opt.NewPtr(req.Params.State),
		CodeChallenge:       string(req.Params.CodeChallenge),
		CodeChallengeMethod: string(req.Params.CodeChallengeMethod),
		AccountID:           account.AccountID(acc),
		AccountPermissions:  permissions,
	})
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthAuthorise400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{
				Error: oauthErr.Code,
			}),
		}, nil
	}

	return openapi.OAuthAuthorise302Response{
		Headers: openapi.OAuthAuthoriseFoundResponseHeaders{Location: result.Location},
	}, nil
}

func (o OAuth) OAuthAuthoriseConsent(ctx context.Context, req openapi.OAuthAuthoriseConsentRequestObject) (openapi.OAuthAuthoriseConsentResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthAuthoriseConsent401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	requestID := ""
	if req.Params.RequestId != nil {
		requestID = string(*req.Params.RequestId)
	}

	consent, oauthErr, err := o.oauth.GetAuthorisationConsent(ctx, account.AccountID(acc), permissions, requestID)
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthAuthoriseConsent400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	return openapi.OAuthAuthoriseConsent200JSONResponse{
		OAuthAuthoriseConsentOKJSONResponse: openapi.OAuthAuthoriseConsentOKJSONResponse(openapi.OAuthAuthoriseConsent{
			RequestId:               consent.RequestID,
			ClientId:                consent.ClientID,
			ClientName:              consent.ClientName,
			RedirectUri:             consent.RedirectURI,
			ExpiresAt:               consent.ExpiresAt,
			RequestedScopes:         consent.RequestedScopes,
			GrantedScopes:           consent.GrantedScopes,
			InheritsUserPermissions: consent.InheritsUserPermissions,
		}),
	}, nil
}

func (o OAuth) OAuthAuthoriseConsentSubmit(ctx context.Context, req openapi.OAuthAuthoriseConsentSubmitRequestObject) (openapi.OAuthAuthoriseConsentSubmitResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthAuthoriseConsentSubmit400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: "invalid_request"}),
		}, nil
	}

	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthAuthoriseConsentSubmit401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	result, oauthErr, err := o.oauth.SubmitAuthorisationConsent(ctx, account.AccountID(acc), permissions, req.Body.RequestId, req.Body.Decision == openapi.OAuthAuthoriseDecisionApprove)
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthAuthoriseConsentSubmit400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	status := openapi.OAuthAuthoriseConsentResultStatusDenied
	if result.Status == "approved" {
		status = openapi.OAuthAuthoriseConsentResultStatusApproved
	}

	return openapi.OAuthAuthoriseConsentSubmit200JSONResponse{
		OAuthAuthoriseConsentSubmitOKJSONResponse: openapi.OAuthAuthoriseConsentSubmitOKJSONResponse(openapi.OAuthAuthoriseConsentResult{
			Status:   status,
			Location: result.Location,
		}),
	}, nil
}

func (o OAuth) OAuthToken(ctx context.Context, req openapi.OAuthTokenRequestObject) (openapi.OAuthTokenResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthToken400JSONResponse{
			OAuthTokenErrorJSONResponse: openapi.OAuthTokenErrorJSONResponse(openapi.OAuthError{
				Error: "invalid_request",
			}),
		}, nil
	}

	token, oauthErr, err := o.oauth.ExchangeToken(ctx, oauthservice.TokenRequest{
		GrantType:    req.Body.GrantType,
		ClientID:     req.Body.ClientId,
		ClientSecret: opt.NewPtr(req.Body.ClientSecret),
		Scope:        opt.NewPtr(req.Body.Scope),
		DeviceCode:   opt.NewPtr(req.Body.DeviceCode),
		Code:         opt.NewPtr(req.Body.Code),
		RedirectURI:  opt.NewPtr(req.Body.RedirectUri),
		CodeVerifier: opt.NewPtr(req.Body.CodeVerifier),
		RefreshToken: opt.NewPtr(req.Body.RefreshToken),
	})
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthToken400JSONResponse{
			OAuthTokenErrorJSONResponse: openapi.OAuthTokenErrorJSONResponse(openapi.OAuthError{
				Error: oauthErr.Code,
			}),
		}, nil
	}

	return openapi.OAuthToken200JSONResponse{
		OAuthTokenOKJSONResponse: openapi.OAuthTokenOKJSONResponse(openapi.OAuthToken{
			AccessToken:  &token.AccessToken,
			TokenType:    &token.TokenType,
			ExpiresIn:    &token.ExpiresIn,
			Scope:        &token.Scope,
			IdToken:      token.IDToken.Ptr(),
			RefreshToken: token.RefreshToken.Ptr(),
		}),
	}, nil
}

func (o OAuth) OAuthUserInfo(ctx context.Context, _ openapi.OAuthUserInfoRequestObject) (openapi.OAuthUserInfoResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthUserInfo401Response{}, nil
	}
	scopes, err := session.GetOAuthScopes(ctx)
	if err != nil {
		return nil, err
	}

	userInfo, err := o.oauth.UserInfo(ctx, account.AccountID(acc), scopes)
	if err != nil {
		return nil, err
	}

	return openapi.OAuthUserInfo200JSONResponse{
		OAuthUserInfoOKJSONResponse: openapi.OAuthUserInfoOKJSONResponse(openapi.OAuthUserInfo{
			Sub:               &userInfo.Subject,
			Name:              userInfo.Name.Ptr(),
			Email:             userInfo.Email.Ptr(),
			EmailVerified:     userInfo.EmailVerified.Ptr(),
			PreferredUsername: userInfo.PreferredUsername.Ptr(),
		}),
	}, nil
}

func (o OAuth) OAuthRefreshTokenList(ctx context.Context, req openapi.OAuthRefreshTokenListRequestObject) (openapi.OAuthRefreshTokenListResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthRefreshTokenList401Response{}, nil
	}

	tokens, err := o.oauth.ListRefreshTokensByAccount(ctx, account.AccountID(acc))
	if err != nil {
		return nil, err
	}

	return openapi.OAuthRefreshTokenList200JSONResponse{
		OAuthRefreshTokenListOKJSONResponse: openapi.OAuthRefreshTokenListOKJSONResponse(openapi.OAuthRefreshTokenListResult{
			Tokens: serialiseOAuthRefreshTokenList(tokens),
		}),
	}, nil
}

func (o OAuth) OAuthRefreshTokenDelete(ctx context.Context, req openapi.OAuthRefreshTokenDeleteRequestObject) (openapi.OAuthRefreshTokenDeleteResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthRefreshTokenDelete401Response{}, nil
	}

	oauthErr := o.oauth.RevokeRefreshTokenByAccount(ctx, account.AccountID(acc), oauthresource.RefreshTokenID(deserialiseID(req.OauthRefreshTokenId)))
	if oauthErr != nil {
		return openapi.OAuthRefreshTokenDelete400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	return openapi.OAuthRefreshTokenDelete204Response{}, nil
}

func (o OAuth) OAuthClientList(ctx context.Context, req openapi.OAuthClientListRequestObject) (openapi.OAuthClientListResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthClientList401Response{}, nil
	}

	clients, err := o.oauth.ListClientsByAccount(ctx, account.AccountID(acc))
	if err != nil {
		return nil, err
	}

	return openapi.OAuthClientList200JSONResponse{
		OAuthClientListOKJSONResponse: openapi.OAuthClientListOKJSONResponse(openapi.OAuthClientListResult{
			Clients: serialiseOAuthClientList(clients),
		}),
	}, nil
}

func (o OAuth) OAuthClientCreate(ctx context.Context, req openapi.OAuthClientCreateRequestObject) (openapi.OAuthClientCreateResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthClientCreate400Response{}, nil
	}

	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthClientCreate401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	result, err := o.oauth.CreateClientForAccount(ctx, oauthservice.ClientSelfCreate{
		AccountID:          account.AccountID(acc),
		AccountPermissions: permissions,
		Name:               req.Body.Name,
		RedirectURIs:       opt.NewPtr(req.Body.RedirectUris).OrZero(),
		AllowedScopes:      req.Body.AllowedScopes,
	})
	if err != nil {
		return nil, err
	}

	return openapi.OAuthClientCreate200JSONResponse{
		OAuthClientIssuedOKJSONResponse: openapi.OAuthClientIssuedOKJSONResponse(serialiseOAuthClientIssued(result)),
	}, nil
}

func (o OAuth) OAuthClientGet(ctx context.Context, req openapi.OAuthClientGetRequestObject) (openapi.OAuthClientGetResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthClientGet401Response{}, nil
	}

	client, oauthErr, err := o.oauth.GetClientByAccount(ctx, account.AccountID(acc), oauthresource.ClientID(deserialiseID(req.OauthClientId)))
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthClientGet400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	return openapi.OAuthClientGet200JSONResponse{
		OAuthClientOKJSONResponse: openapi.OAuthClientOKJSONResponse(serialiseOAuthClient(client)),
	}, nil
}

func (o OAuth) OAuthClientUpdate(ctx context.Context, req openapi.OAuthClientUpdateRequestObject) (openapi.OAuthClientUpdateResponseObject, error) {
	if req.Body == nil {
		return openapi.OAuthClientUpdate400Response{}, nil
	}

	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthClientUpdate401Response{}, nil
	}
	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	client, oauthErr, err := o.oauth.UpdateClientByAccount(ctx, account.AccountID(acc), oauthresource.ClientID(deserialiseID(req.OauthClientId)), oauthservice.ClientSelfUpdate{
		AccountPermissions: permissions,
		Name:               opt.NewPtr(req.Body.Name),
		RedirectURIs:       opt.NewPtr(req.Body.RedirectUris),
		AllowedScopes:      opt.NewPtr(req.Body.AllowedScopes),
	})
	if err != nil {
		return nil, err
	}
	if oauthErr != nil {
		return openapi.OAuthClientUpdate400Response{}, nil
	}

	return openapi.OAuthClientUpdate200JSONResponse{
		OAuthClientOKJSONResponse: openapi.OAuthClientOKJSONResponse(serialiseOAuthClient(client)),
	}, nil
}

func (o OAuth) OAuthClientDelete(ctx context.Context, req openapi.OAuthClientDeleteRequestObject) (openapi.OAuthClientDeleteResponseObject, error) {
	acc, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.OAuthClientDelete401Response{}, nil
	}

	oauthErr := o.oauth.DeleteClientByAccount(ctx, account.AccountID(acc), oauthresource.ClientID(deserialiseID(req.OauthClientId)))
	if oauthErr != nil {
		return openapi.OAuthClientDelete400JSONResponse{
			OAuthErrorJSONResponse: openapi.OAuthErrorJSONResponse(openapi.OAuthError{Error: oauthErr.Code}),
		}, nil
	}

	return openapi.OAuthClientDelete204Response{}, nil
}

func mapOAuthJWKs(in []oauthservice.JWK) []openapi.OAuthJWK {
	out := make([]openapi.OAuthJWK, len(in))
	for i, key := range in {
		out[i] = openapi.OAuthJWK{
			Kty: key.Kty,
			Use: key.Use,
			Alg: key.Alg,
			Kid: key.Kid,
			N:   key.N,
			E:   key.E,
		}
	}

	return out
}

func serialiseOAuthClient(in *oauthresource.Client) openapi.OAuthClient {
	return openapi.OAuthClient{
		Id:            openapi.Identifier(in.ID.XID().String()),
		CreatedAt:     openapi.CreatedAt(in.CreatedAt),
		UpdatedAt:     openapi.UpdatedAt(in.UpdatedAt),
		AccountId:     opt.Map(in.AccountID, func(id account.AccountID) openapi.Identifier { return openapi.Identifier(id.String()) }).Ptr(),
		ClientId:      in.ClientID,
		Name:          in.Name,
		Type:          openapi.OAuthClientType(in.Type.String()),
		ScopePolicy:   openapi.OAuthClientScopePolicy(in.ScopePolicy.String()),
		RedirectUris:  in.RedirectURIs,
		AllowedScopes: in.AllowedScopes,
		AllowedGrants: in.AllowedGrants,
	}
}

func serialiseOAuthClientIssued(in *oauthservice.ClientSelfCreateResult) openapi.OAuthClientIssued {
	return openapi.OAuthClientIssued{
		Client:       serialiseOAuthClient(in.Client),
		ClientSecret: in.ClientSecret.Ptr(),
	}
}

func serialiseOAuthClientList(in []*oauthresource.Client) openapi.OAuthClientList {
	return dt.Map(in, serialiseOAuthClient)
}

func serialiseOAuthDeviceAuthorisation(in *oauthresource.DeviceAuthorisation) openapi.OAuthDeviceAuthorisationRecord {
	return openapi.OAuthDeviceAuthorisationRecord{
		Id:                  openapi.Identifier(in.ID.XID().String()),
		CreatedAt:           openapi.CreatedAt(in.CreatedAt),
		ClientId:            openapi.Identifier(in.ClientID.XID().String()),
		UserCode:            in.UserCodeDisplay,
		Scope:               in.Scope,
		ExpiresAt:           in.ExpiresAt,
		PollIntervalSeconds: in.PollIntervalSeconds,
		LastPolledAt:        in.LastPolledAt.Ptr(),
		ApprovedByAccountId: opt.Map(in.ApprovedByAccountID, func(id account.AccountID) openapi.Identifier { return openapi.Identifier(id.String()) }).Ptr(),
		ApprovedAt:          in.ApprovedAt.Ptr(),
		DeniedAt:            in.DeniedAt.Ptr(),
		ConsumedAt:          in.ConsumedAt.Ptr(),
	}
}

func serialiseOAuthDeviceAuthorisationList(in []*oauthresource.DeviceAuthorisation) openapi.OAuthDeviceAuthorisationList {
	return dt.Map(in, serialiseOAuthDeviceAuthorisation)
}

func serialiseOAuthRefreshToken(in *oauthresource.RefreshToken) openapi.OAuthRefreshToken {
	clientID := in.ClientIdentifier
	if clientID == "" {
		clientID = in.ClientID.XID().String()
	}

	clientName := in.ClientName
	if clientName == "" {
		clientName = clientID
	}

	return openapi.OAuthRefreshToken{
		Id:            openapi.Identifier(in.ID.XID().String()),
		CreatedAt:     openapi.CreatedAt(in.CreatedAt),
		OauthClientId: openapi.Identifier(in.ClientID.XID().String()),
		ClientId:      clientID,
		ClientName:    clientName,
		AccountId:     openapi.Identifier(in.AccountID.String()),
		Scope:         in.Scope,
		ExpiresAt:     in.ExpiresAt,
		RevokedAt:     in.RevokedAt.Ptr(),
		ReplacedByTokenId: opt.Map(in.ReplacedByTokenID, func(id oauthresource.RefreshTokenID) openapi.Identifier {
			return openapi.Identifier(id.XID().String())
		}).Ptr(),
		LastUsedAt: in.LastUsedAt.Ptr(),
	}
}

func serialiseOAuthRefreshTokenList(in []*oauthresource.RefreshToken) openapi.OAuthRefreshTokenList {
	return dt.Map(in, serialiseOAuthRefreshToken)
}
