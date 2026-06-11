package oauth

import (
	"context"
	"net/url"
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

type AuthoriseRequest struct {
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               opt.Optional[string]
	State               opt.Optional[string]
	Nonce               opt.Optional[string]
	CodeChallenge       string
	CodeChallengeMethod string
	AccountID           account.AccountID
	AccountPermissions  rbac.Permissions
}

type AuthoriseResult struct {
	Location string
}

type AuthorisationConsent struct {
	RequestID               string
	ClientID                string
	ClientName              string
	RedirectURI             string
	ExpiresAt               time.Time
	RequestedScopes         []string
	GrantedScopes           []string
	InheritsUserPermissions bool
}

type AuthorisationConsentResult struct {
	Status   string
	Location string
}

func (s *Service) Authorise(ctx context.Context, input AuthoriseRequest) (*AuthoriseResult, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "OAuth is not enabled on this instance"), nil
	}

	if input.ResponseType != "code" || input.CodeChallengeMethod != CodeChallengeMethodS256 || !validCodeVerifier(input.CodeChallenge) {
		return nil, oauthError("invalid_request", "Invalid response_type, code_challenge_method, or code_challenge"), nil
	}
	if !canAuthoriseOAuthClients(input.AccountPermissions) {
		return nil, oauthError("access_denied", "Account is not permitted to authorise OAuth clients"), nil
	}

	cl, oauthErr, err := s.resolveClient(ctx, input.ClientID)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}
	if oauthErr != nil {
		return nil, oauthErr, err
	}

	if !contains(cl.AllowedGrants, GrantTypeAuthorizationCode) {
		return nil, oauthError("unauthorized_client", "Client is not authorized for authorization_code grant"), nil
	}
	if !contains(cl.RedirectURIs, input.RedirectURI) {
		return nil, oauthError("invalid_request", "Redirect URI is not registered for this client"), nil
	}

	scope := strings.TrimSpace(input.Scope.OrZero())
	if _, err := grantScope(scope, cl, input.AccountPermissions); err != nil {
		return nil, oauthError("invalid_scope", "Requested scope is not permitted for this account"), nil
	}

	requestID, err := randomToken(32)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.tokens.CreateAuthorisationRequest(ctx, oauth_writer.AuthorisationRequestCreate{
		ClientID:            cl.ID,
		AccountID:           input.AccountID,
		RequestIDHash:       hashString(requestID),
		RedirectURI:         input.RedirectURI,
		Scope:               scope,
		State:               input.State,
		Nonce:               input.Nonce,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: CodeChallengeMethodS256,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		return nil, nil, err
	}

	return &AuthoriseResult{Location: s.authorizationCodeConsentURL(requestID)}, nil, nil
}

func (s *Service) GetAuthorisationConsent(ctx context.Context, accountID account.AccountID, accountPermissions rbac.Permissions, requestID string) (*AuthorisationConsent, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "OAuth is not enabled on this instance"), nil
	}

	rec, oauthErr, err := s.getPendingAuthorisationRequest(ctx, accountID, requestID)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}
	if !canAuthoriseOAuthClients(accountPermissions) {
		return nil, oauthError("access_denied", "Account is not permitted to authorise OAuth clients"), nil
	}

	cl, err := s.clients.GetClient(ctx, rec.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client", "Client not found"), nil
	}

	grantedScope, err := grantScope(rec.Scope, cl, accountPermissions)
	if err != nil {
		return nil, oauthError("invalid_scope", "Requested scope is not permitted for this account"), nil
	}

	return &AuthorisationConsent{
		RequestID:               requestID,
		ClientID:                cl.ClientID,
		ClientName:              cl.Name,
		RedirectURI:             rec.RedirectURI,
		ExpiresAt:               rec.ExpiresAt,
		RequestedScopes:         splitScope(rec.Scope),
		GrantedScopes:           splitScope(grantedScope),
		InheritsUserPermissions: shouldInheritUserPermissions(cl),
	}, nil, nil
}

func (s *Service) SubmitAuthorisationConsent(ctx context.Context, accountID account.AccountID, accountPermissions rbac.Permissions, requestID string, approved bool) (*AuthorisationConsentResult, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "OAuth is not enabled on this instance"), nil
	}

	rec, oauthErr, err := s.getPendingAuthorisationRequest(ctx, accountID, requestID)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	if !approved {
		ok, err := s.tokens.DenyAuthorisationRequest(ctx, rec.ID, time.Now())
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			return nil, oauthError("invalid_request", "Failed to record denial"), nil
		}

		return &AuthorisationConsentResult{
			Status:   "denied",
			Location: authorizationErrorRedirect(rec.RedirectURI, "access_denied", rec.State),
		}, nil, nil
	}
	if !canAuthoriseOAuthClients(accountPermissions) {
		return nil, oauthError("access_denied", "Account is not permitted to authorise OAuth clients"), nil
	}

	cl, err := s.clients.GetClient(ctx, rec.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client", "Client not found"), nil
	}

	grantedScope, err := grantScope(rec.Scope, cl, accountPermissions)
	if err != nil {
		return nil, oauthError("invalid_scope", "Requested scope is not permitted for this account"), nil
	}

	code, err := randomToken(32)
	if err != nil {
		return nil, nil, err
	}
	ok, err := s.tokens.ApproveAuthorisationRequestAndCreateCode(ctx, rec.ID, oauth_writer.AuthorisationCodeCreate{
		ClientID:            rec.ClientID,
		AccountID:           rec.AccountID,
		CodeHash:            hashString(code),
		RedirectURI:         rec.RedirectURI,
		Scope:               grantedScope,
		Nonce:               rec.Nonce,
		CodeChallenge:       rec.CodeChallenge,
		CodeChallengeMethod: CodeChallengeMethodS256,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
	}, time.Now())
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, oauthError("invalid_request", "Failed to approve authorisation request"), nil
	}

	return &AuthorisationConsentResult{
		Status:   "approved",
		Location: authorizationCodeRedirect(rec.RedirectURI, code, rec.State),
	}, nil, nil
}

func (s *Service) getPendingAuthorisationRequest(ctx context.Context, accountID account.AccountID, requestID string) (*oauthresource.AuthorisationRequest, *Error, error) {
	if strings.TrimSpace(requestID) == "" {
		return nil, oauthError("invalid_request", "Missing request_id"), nil
	}

	rec, err := s.clients.GetAuthorisationRequestByRequestIDHash(ctx, hashString(requestID))
	if err != nil {
		return nil, oauthError("invalid_request", "Request ID not found"), nil
	}

	if rec.AccountID != accountID || rec.ExpiresAt.Before(time.Now()) || rec.ApprovedAt.Ok() || rec.DeniedAt.Ok() {
		return nil, oauthError("invalid_request", "Request is invalid, expired, or already processed"), nil
	}

	return rec, nil, nil
}

func authorizationCodeRedirect(redirectURI string, code string, state opt.Optional[string]) string {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("code", code)
	state.Call(func(value string) {
		q.Set("state", value)
	})
	u.RawQuery = q.Encode()

	return u.String()
}

func authorizationErrorRedirect(redirectURI string, errorCode string, state opt.Optional[string]) string {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("error", errorCode)
	state.Call(func(value string) {
		q.Set("state", value)
	})
	u.RawQuery = q.Encode()

	return u.String()
}
