package oauth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	oauthGrantAuthorizationCode = "authorization_code"
	oauthGrantRefreshToken      = "refresh_token"
	oauthGrantClientCredentials = "client_credentials"
	oauthGrantDeviceCode        = "urn:ietf:params:oauth:grant-type:device_code"
)

var errNilResponse = errors.New("nil oauth token response")

type refreshResult struct {
	status  int
	errCode string
	err     error
}

type oauthTokenRequest struct {
	GrantType    string
	ClientId     string
	ClientSecret *string
	Scope        *string
	DeviceCode   *string
	Code         *string
	RedirectUri  *string
	CodeVerifier *string
	RefreshToken *string
}

func oauthToken(_ *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, req oauthTokenRequest) (*openapi.OAuthTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", req.GrantType)
	form.Set("client_id", req.ClientId)
	setFormValue(form, "client_secret", req.ClientSecret)
	setFormValue(form, "scope", req.Scope)
	setFormValue(form, "device_code", req.DeviceCode)
	setFormValue(form, "code", req.Code)
	setFormValue(form, "redirect_uri", req.RedirectUri)
	setFormValue(form, "code_verifier", req.CodeVerifier)
	setFormValue(form, "refresh_token", req.RefreshToken)

	return cl.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
}

func setFormValue(form url.Values, name string, value *string) {
	if value == nil {
		return
	}

	form.Set(name, *value)
}

func oauthConfig(t *testing.T) *config.Config {
	t.Helper()
	return oauthConfigWithAccessTTL(t, 15*time.Minute)
}

func oauthConfigWithAccessTTL(t *testing.T, accessTTL time.Duration) *config.Config {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	publicWebAddress, err := url.Parse("http://localhost:3000")
	require.NoError(t, err)
	publicAPIAddress, err := url.Parse("http://localhost:8000")
	require.NoError(t, err)

	return &config.Config{
		PublicWebAddress:      *publicWebAddress,
		PublicAPIAddress:      *publicAPIAddress,
		OAuthEnabled:          true,
		OAuthSigningKeyBase64: base64.StdEncoding.EncodeToString(pemBytes),
		OAuthSigningKeyID:     "test-key",
		OAuthAccessTokenTTL:   accessTTL,
		OAuthRefreshTokenTTL:  24 * time.Hour,
		OAuthDeviceCodeTTL:    10 * time.Minute,
		OAuthDevicePollEvery:  5 * time.Second,
	}
}

func createClient(
	t *testing.T,
	ctx context.Context,
	ow *oauth_writer.Writer,
	accountID account.AccountID,
	clientID string,
	clientType oauthresource.ClientType,
	scopePolicy oauthresource.ScopePolicy,
	secretHash opt.Optional[string],
	allowedScopes []string,
	allowedGrants []string,
) *oauthresource.Client {
	t.Helper()

	client, err := ow.CreateClient(ctx, oauth_writer.ClientCreate{
		AccountID:        opt.New(accountID),
		ClientID:         clientID,
		ClientSecretHash: secretHash,
		Name:             clientID,
		Type:             clientType,
		ScopePolicy:      opt.New(scopePolicy),
		RedirectURIs:     []string{"https://client.example/callback"},
		AllowedScopes:    allowedScopes,
		AllowedGrants:    allowedGrants,
	})
	require.NoError(t, err)

	return client
}

func grantOAuthClientUse(
	t *testing.T,
	ctx context.Context,
	roles *role_repo.Repository,
	assignments *role_assign.Assignment,
	accountID account.AccountID,
	extra ...rbac.Permission,
) role.RoleID {
	t.Helper()

	permissions := append(rbac.PermissionList{rbac.PermissionUseOauthClients}, extra...)
	created, err := roles.Create(ctx, "oauth-client-test-"+uuid.NewString(), "blue", permissions)
	require.NoError(t, err)

	err = assignments.UpdateRoles(ctx, account_ref.ID(accountID), role_assign.Add(created.ID))
	require.NoError(t, err)

	return created.ID
}

func revokeOAuthClientUse(t *testing.T, ctx context.Context, assignments *role_assign.Assignment, accountID account.AccountID, roleID role.RoleID) {
	t.Helper()

	err := assignments.UpdateRoles(ctx, account_ref.ID(accountID), role_assign.Remove(roleID))
	require.NoError(t, err)
}

func clientSecretHash(t *testing.T, secret string) string {
	t.Helper()

	hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
	require.NoError(t, err)

	return hash
}

func authorizeCode(t *testing.T, ctx context.Context, ts *httptest.Server, session openapi.RequestEditorFn, clientID, redirectURI, scope, verifier string) string {
	t.Helper()

	location := authorizeRedirect(t, ctx, ts, session, authorizeRequest{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		Scope:               scope,
		State:               "state-" + uuid.NewString(),
		CodeChallenge:       codeChallenge(verifier),
		CodeChallengeMethod: "S256",
	})

	u, err := url.Parse(location)
	require.NoError(t, err)
	require.Equal(t, redirectURI, u.Scheme+"://"+u.Host+u.Path)

	code := u.Query().Get("code")
	require.NotEmpty(t, code)

	return code
}

type authorizeRequest struct {
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

func authorizeRedirect(t *testing.T, ctx context.Context, ts *httptest.Server, session openapi.RequestEditorFn, req authorizeRequest) string {
	t.Helper()

	resp := authorizeHTTPResponse(t, ctx, ts, session, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	require.NotEmpty(t, location)

	return location
}

func authorizeHTTPResponse(t *testing.T, ctx context.Context, ts *httptest.Server, session openapi.RequestEditorFn, req authorizeRequest) *http.Response {
	t.Helper()

	responseType := req.ResponseType
	if responseType == "" {
		responseType = "code"
	}
	codeChallengeMethod := req.CodeChallengeMethod
	if codeChallengeMethod == "" {
		codeChallengeMethod = "S256"
	}

	q := url.Values{}
	q.Set("response_type", responseType)
	q.Set("client_id", req.ClientID)
	q.Set("redirect_uri", req.RedirectURI)
	if req.Scope != "" {
		q.Set("scope", req.Scope)
	}
	q.Set("state", req.State)
	q.Set("code_challenge", req.CodeChallenge)
	q.Set("code_challenge_method", codeChallengeMethod)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/oauth/authorize?"+q.Encode(), nil)
	require.NoError(t, err)
	require.NoError(t, session(ctx, httpReq))

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Do(httpReq)
	require.NoError(t, err)

	return resp
}

func refreshTwiceConcurrently(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, clientID, clientSecret, refreshToken string) []refreshResult {
	t.Helper()

	var wg sync.WaitGroup
	out := make([]refreshResult, 2)

	for i := range out {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := oauthToken(t, ctx, cl, oauthTokenRequest{
				GrantType:    oauthGrantRefreshToken,
				ClientId:     clientID,
				ClientSecret: &clientSecret,
				RefreshToken: &refreshToken,
			})
			if err != nil {
				out[i].err = err
				return
			}
			if resp == nil {
				out[i].err = errNilResponse
				return
			}

			out[i].status = resp.StatusCode()
			if resp.JSON400 != nil {
				out[i].errCode = resp.JSON400.Error
			}
		}(i)
	}

	wg.Wait()

	return out
}

func parseClaims(t *testing.T, raw string) jwt.MapClaims {
	t.Helper()

	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(raw, claims)
	require.NoError(t, err)

	return claims
}

func codeChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func bearer(token string) openapi.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

func standardScopes() []string {
	return []string{"openid", "profile", "email", "offline_access"}
}

func ptr[T any](v T) *T {
	return &v
}
