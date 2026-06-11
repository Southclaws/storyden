package oauth

import (
	"context"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type DeviceAuthorisation struct {
	DeviceCode              string
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresIn               int
	Interval                int
}

type DeviceConsent struct {
	UserCode                string
	ClientID                string
	ClientName              string
	ExpiresAt               time.Time
	RequestedScopes         []string
	GrantedScopes           []string
	InheritsUserPermissions bool
}

func (s *Service) StartDeviceAuthorization(ctx context.Context, clientID string, requestedScope opt.Optional[string]) (*DeviceAuthorisation, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "OAuth is not enabled on this instance"), nil
	}

	cl, err := s.getClientForDeviceAuthorization(ctx, clientID)
	if err != nil {
		return nil, oauthError("invalid_client", "Client not found"), nil
	}
	if cl.Type == oauthresource.ClientTypeConfidential {
		return nil, oauthError("unauthorized_client", "Confidential clients cannot use device authorization"), nil
	}
	if !contains(cl.AllowedGrants, GrantTypeDeviceCode) {
		return nil, oauthError("unauthorized_client", "Client is not authorized for device_code grant"), nil
	}

	scope := strings.TrimSpace(requestedScope.OrZero())
	if cl.ClientID == StorydenCLIClientID && !isStorydenDeviceScope(scope) {
		return nil, oauthError("invalid_scope", "Storyden CLI requires openid, profile, and offline_access scopes"), nil
	}
	if err := validateScopeNames(scope); err != nil {
		return nil, oauthError("invalid_scope", "Requested scope contains invalid scope names"), nil
	}
	if err := authorizeScopeNames(scope, cl.AllowedScopes); err != nil {
		return nil, oauthError("invalid_scope", "Requested scope is not permitted for this client"), nil
	}

	deviceCode, err := randomToken(32)
	if err != nil {
		return nil, nil, err
	}
	rawToken, err := randomToken(8)
	if err != nil {
		return nil, nil, err
	}
	raw := strings.ToUpper(rawToken)
	userCode := raw[:4] + "-" + raw[4:8]

	if _, err := s.tokens.DeleteExpiredDeviceAuthorisations(ctx, time.Now()); err != nil {
		return nil, nil, err
	}

	_, err = s.tokens.CreateDeviceAuthorisation(ctx, oauth_writer.DeviceAuthorisationCreate{
		ClientID:            cl.ID,
		DeviceCodeHash:      hashString(deviceCode),
		UserCodeHash:        hashString(normalizeCode(userCode)),
		UserCodeDisplay:     userCode,
		Scope:               scope,
		ExpiresAt:           time.Now().Add(s.cfg.OAuthDeviceCodeTTL),
		PollIntervalSeconds: int(s.cfg.OAuthDevicePollEvery.Seconds()),
	})
	if err != nil {
		return nil, nil, err
	}

	verificationURI := s.deviceAuthorizationConsentURL("")
	verificationURIComplete := s.deviceAuthorizationConsentURL(userCode)

	return &DeviceAuthorisation{
		DeviceCode:              deviceCode,
		UserCode:                userCode,
		VerificationURI:         verificationURI,
		VerificationURIComplete: verificationURIComplete,
		ExpiresIn:               int(s.cfg.OAuthDeviceCodeTTL.Seconds()),
		Interval:                int(s.cfg.OAuthDevicePollEvery.Seconds()),
	}, nil, nil
}

func (s *Service) getClientForDeviceAuthorization(ctx context.Context, clientID string) (*oauthresource.Client, error) {
	cl, err := s.clients.GetClientByClientID(ctx, clientID)
	if err == nil {
		if clientID == StorydenCLIClientID && cl.Name != "Storyden" {
			cl, err = s.tokens.UpdateClient(ctx, cl.ID, oauth_writer.ClientUpdate{
				Name: opt.New("Storyden"),
			})
			if err != nil {
				return nil, err
			}
		}
		return cl, nil
	}

	if clientID == StorydenCLIClientID {
		_, err = s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
			AccountID:        opt.NewEmpty[account.AccountID](),
			ClientID:         StorydenCLIClientID,
			ClientSecretHash: opt.NewEmpty[string](),
			Name:             "Storyden",
			Type:             oauthresource.ClientTypePublic,
			ScopePolicy:      opt.New(oauthresource.ScopePolicyInheritUserPermissions),
			RedirectURIs:     []string{},
			AllowedScopes:    supportedScopes(),
			AllowedGrants:    []string{GrantTypeDeviceCode, GrantTypeRefreshToken},
		})
		if err != nil {
			cl, readErr := s.clients.GetClientByClientID(ctx, clientID)
			if readErr == nil {
				return cl, nil
			}
			return nil, err
		}
		return s.clients.GetClientByClientID(ctx, clientID)
	}

	if isCIMDClientID(clientID) {
		cl, oauthErr, err := s.resolveCIMDClient(ctx, clientID)
		if oauthErr != nil || err != nil {
			if err != nil {
				return nil, err
			}
			return nil, fault.New("invalid CIMD client for device authorization")
		}
		return cl, nil
	}

	return nil, err
}

func (s *Service) GetDeviceConsent(ctx context.Context, accountID account.AccountID, accountPermissions rbac.Permissions, userCode string) (*DeviceConsent, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "OAuth is not enabled on this instance"), nil
	}

	rec, oauthErr, err := s.getPendingDeviceAuthorisation(ctx, userCode)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}
	if !canAuthoriseOAuthClients(accountPermissions) {
		return nil, oauthError("access_denied", "Account is not permitted to authorise OAuth clients"), nil
	}

	claimed, err := s.tokens.ClaimDeviceAuthorisation(ctx, rec.ID, accountID)
	if err != nil {
		return nil, nil, err
	}
	if !claimed {
		return nil, oauthError("access_denied", "Device authorisation has already been claimed by another account"), nil
	}

	cl, err := s.clients.GetClient(ctx, rec.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client", "Client not found"), nil
	}

	grantedScope, err := grantScope(rec.Scope, cl, accountPermissions)
	if err != nil {
		return nil, oauthError("invalid_scope", "Requested scope is not permitted for this account"), nil
	}

	return &DeviceConsent{
		UserCode:                rec.UserCodeDisplay,
		ClientID:                cl.ClientID,
		ClientName:              cl.Name,
		ExpiresAt:               rec.ExpiresAt,
		RequestedScopes:         strings.Fields(rec.Scope),
		GrantedScopes:           strings.Fields(grantedScope),
		InheritsUserPermissions: shouldInheritUserPermissions(cl),
	}, nil, nil
}

func (s *Service) ApproveDeviceAuthorization(ctx context.Context, accountID account.AccountID, accountPermissions rbac.Permissions, userCode string, approved bool) *Error {
	if !s.Enabled() {
		return oauthError("temporarily_unavailable", "OAuth is not enabled on this instance")
	}

	rec, oauthErr, err := s.getPendingDeviceAuthorisation(ctx, userCode)
	if oauthErr != nil || err != nil {
		return oauthErr
	}

	if claimant, ok := rec.ClaimedByAccountID.Get(); !ok || claimant != accountID {
		return oauthError("access_denied", "Only the claiming account can approve or deny this request")
	}

	if approved {
		if !canAuthoriseOAuthClients(accountPermissions) {
			return oauthError("access_denied", "Account is not permitted to authorise OAuth clients")
		}

		cl, err := s.clients.GetClient(ctx, rec.ClientID)
		if err != nil {
			return oauthError("invalid_client", "Client not found")
		}

		grantedScope, err := grantScope(rec.Scope, cl, accountPermissions)
		if err != nil {
			return oauthError("invalid_scope", "Requested scope is not permitted for this account")
		}

		ok, err := s.tokens.ApproveDeviceAuthorisation(ctx, rec.ID, accountID, grantedScope, time.Now())
		if err != nil || !ok {
			return oauthError("invalid_request", "Failed to approve device authorisation")
		}
		return nil
	}

	ok, err := s.tokens.DenyDeviceAuthorisation(ctx, rec.ID, time.Now())
	if err != nil || !ok {
		return oauthError("invalid_request", "Failed to deny device authorisation")
	}
	return nil
}

func (s *Service) getPendingDeviceAuthorisation(ctx context.Context, userCode string) (*oauthresource.DeviceAuthorisation, *Error, error) {
	if strings.TrimSpace(userCode) == "" {
		return nil, oauthError("invalid_request", "Missing user_code"), nil
	}

	rec, err := s.clients.GetDeviceAuthorisationByUserCodeHash(ctx, hashString(normalizeCode(userCode)))
	if err != nil {
		return nil, oauthError("invalid_request", "User code not found"), nil
	}

	if rec.ExpiresAt.Before(time.Now()) || rec.ConsumedAt.Ok() || rec.ApprovedAt.Ok() || rec.DeniedAt.Ok() {
		return nil, oauthError("invalid_request", "User code has expired or already been used"), nil
	}

	return rec, nil, nil
}

func (s *Service) exchangeDeviceCode(ctx context.Context, input TokenRequest) (*Token, *Error, error) {
	deviceCode, ok := input.DeviceCode.Get()
	if !ok {
		return nil, oauthError("invalid_request", "Missing device_code"), nil
	}

	cl, oauthErr, err := s.resolveClient(ctx, input.ClientID)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}
	if oauthErr != nil {
		return nil, oauthErr, err
	}

	if cl.Type == oauthresource.ClientTypeConfidential || !contains(cl.AllowedGrants, GrantTypeDeviceCode) {
		return nil, oauthError("unauthorized_client", "Client is not authorized for device_code grant"), nil
	}

	rec, err := s.clients.GetDeviceAuthorisationByDeviceCodeHash(ctx, hashString(deviceCode))
	if err != nil || rec.ClientID != cl.ID {
		return nil, oauthError("invalid_grant", "Device code is invalid or does not match client"), nil
	}

	now := time.Now()
	if rec.ExpiresAt.Before(now) {
		return nil, oauthError("expired_token", "Device code has expired"), nil
	}
	if rec.DeniedAt.Ok() {
		return nil, oauthError("access_denied", "Device authorisation was denied"), nil
	}
	if rec.ConsumedAt.Ok() {
		return nil, oauthError("invalid_grant", "Device code has already been used"), nil
	}
	if lastPolledAt, ok := rec.LastPolledAt.Get(); ok && now.Sub(lastPolledAt) < time.Duration(rec.PollIntervalSeconds)*time.Second {
		if err := s.tokens.RecordDeviceAuthorisationPoll(ctx, rec.ID, now, rec.PollIntervalSeconds+5); err != nil {
			return nil, nil, err
		}
		return nil, oauthError("slow_down", "Polling too frequently; increase the interval"), nil
	}

	if err := s.tokens.RecordDeviceAuthorisationPoll(ctx, rec.ID, now, rec.PollIntervalSeconds); err != nil {
		return nil, nil, err
	}

	accountID, ok := rec.ApprovedByAccountID.Get()
	if !ok {
		return nil, oauthError("authorization_pending", "User has not yet approved the device authorisation"), nil
	}

	consumed, err := s.tokens.ConsumeDeviceAuthorisation(ctx, rec.ID, now)
	if err != nil {
		return nil, nil, err
	}
	if !consumed {
		return nil, oauthError("invalid_grant", "Device authorisation has already been consumed"), nil
	}

	token, err := s.issueTokens(ctx, cl, accountID, rec.Scope)
	if err != nil {
		return nil, nil, err
	}

	return token, nil, nil
}

func isStorydenDeviceScope(scope string) bool {
	scopes := splitScope(scope)
	if len(scopes) != 3 {
		return false
	}

	return contains(scopes, "openid") &&
		contains(scopes, "profile") &&
		contains(scopes, "offline_access")
}
