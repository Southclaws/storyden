package oauth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

var (
	cimdFetchTimeout           = 10 * time.Second
	cimdMaxResponseBytes int64 = 5 * 1024
	cimdDefaultCacheTTL        = 5 * time.Minute
	cimdMaxCacheTTL            = time.Hour
	cimdMaxCacheEntries        = 4096
)

var cimdDisallowedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
}

var cimdLookupNetIP = net.DefaultResolver.LookupNetIP

// cimdPrivilegedScopes are permission scopes that grant administrative or
// moderation capabilities. They are never granted to CIMD clients unless the
// server explicitly opts in via OAUTH_CIMD_ALLOW_PRIVILEGED_SCOPES.
var cimdPrivilegedScopes = map[string]struct{}{
	rbac.PermissionAdministrator.String():         {},
	rbac.PermissionManageSettings.String():        {},
	rbac.PermissionManageAccounts.String():        {},
	rbac.PermissionManageRoles.String():           {},
	rbac.PermissionManageReports.String():         {},
	rbac.PermissionManageWarnings.String():        {},
	rbac.PermissionManageSuspensions.String():     {},
	rbac.PermissionViewAccounts.String():          {},
	rbac.PermissionViewModerationNotes.String():   {},
	rbac.PermissionManageModerationNotes.String(): {},
	rbac.PermissionUseOauthClients.String():       {},
	rbac.PermissionUsePersonalAccessKeys.String(): {},
}

// cimdDefaultAllowedPermissionScopes is the conservative read-only allowlist
// used when the deployment does not configure OAUTH_CIMD_ALLOWED_SCOPES.
var cimdDefaultAllowedPermissionScopes = []string{
	rbac.PermissionReadPublishedThreads.String(),
	rbac.PermissionReadPublishedLibrary.String(),
	rbac.PermissionListProfiles.String(),
	rbac.PermissionReadProfile.String(),
	rbac.PermissionListCollections.String(),
	rbac.PermissionReadCollection.String(),
}

// clientMetadataDocument mirrors the JSON document hosted at a CIMD client_id URL.
type clientMetadataDocument struct {
	ClientID                string   `json:"client_id"`
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	Scope                   string   `json:"scope"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	ClientSecret            string   `json:"client_secret"`
	ClientSecretExpiresAt   int64    `json:"client_secret_expires_at"`
	LogoURI                 string   `json:"logo_uri"`
	ClientURI               string   `json:"client_uri"`
	TOSURI                  string   `json:"tos_uri"`
	PolicyURI               string   `json:"policy_uri"`
	JWKSURI                 string   `json:"jwks_uri"`
}

type cimdCache struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

func newCIMDCache() *cimdCache {
	return &cimdCache{entries: map[string]time.Time{}}
}

func (c *cimdCache) fresh(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	exp, ok := c.entries[key]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(c.entries, key)
		return false
	}
	return true
}

func (c *cimdCache) store(key string, ttl time.Duration) {
	if ttl <= 0 {
		c.invalidate(key)
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = time.Now().Add(ttl)
	if len(c.entries) > cimdMaxCacheEntries {
		for k := range c.entries {
			delete(c.entries, k)
			break
		}
	}
}

func (c *cimdCache) invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (s *Service) cimdEnabled() bool {
	return s.Enabled() && s.cfg.OAuthClientIDMetadataDocumentEnabled
}

func (s *Service) cimdAllowInsecure() bool {
	return s.cfg.OAuthCIMDAllowInsecureFetch
}

// isCIMDClientID reports whether a client_id should be treated as a Client ID
// Metadata Document URL rather than a locally registered client identifier.
func isCIMDClientID(clientID string) bool {
	u, err := url.Parse(clientID)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(u.Scheme)
	return (scheme == "https" || scheme == "http") && u.Host != ""
}

func (s *Service) resolveClient(ctx context.Context, clientID string) (*oauth.Client, *Error, error) {
	// CIMD client IDs (https/https URLs) are always resolved through the CIMD
	// path when the feature is enabled. This ensures the in-memory cimdCache
	// freshness check + potential re-fetch + upsert happens, so that updates
	// to the hosted metadata document (new redirect_uris, grants, scopes, name)
	// are observed instead of being permanently stuck at the first-seen DB
	// snapshot. The DB row is treated as a durable cache of the last snapshot.
	// When the feature is disabled we fall back to plain DB lookup so that any
	// previously created CIMD snapshots continue to work as ordinary clients.
	if isCIMDClientID(clientID) && s.cimdEnabled() {
		return s.resolveCIMDClient(ctx, clientID)
	}

	cl, err := s.clients.GetClientByClientID(ctx, clientID)
	if err == nil {
		return cl, nil, nil
	}

	if ftag.Get(err) != ftag.NotFound {
		return cl, nil, err
	}

	return nil, oauthError("invalid_client", "Client not found"), nil
}

func (s *Service) resolveCIMDClient(ctx context.Context, rawURL string) (*oauth.Client, *Error, error) {
	if !s.cimdEnabled() {
		return nil, oauthError("unauthorized_client", "CIMD is not enabled"), nil
	}

	u, err := s.parseCIMDClientID(rawURL)
	if err != nil {
		return nil, oauthError("invalid_client", fmt.Sprintf("Invalid CIMD client ID: '%s' %v", rawURL, err)), nil
	}

	if s.cimdCache.fresh(rawURL) {
		if cl, err := s.clients.GetClientByClientID(ctx, rawURL); err == nil {
			return cl, nil, nil
		}
		s.cimdCache.invalidate(rawURL)
	}

	doc, ttl, oauthErr := s.fetchClientMetadata(ctx, u)
	if oauthErr != nil {
		s.cimdCache.invalidate(rawURL)
		return nil, oauthErr, nil
	}

	if err := validateClientMetadata(doc, u.String()); err != nil {
		s.cimdCache.invalidate(rawURL)
		return nil, oauthError("invalid_client", "Invalid CIMD client metadata"), nil
	}

	allowedScopes := s.cimdAllowedScopes(doc.Scope)
	allowedGrants := cimdAllowedGrants(doc, allowedScopes)

	cl, err := s.upsertCIMDClient(ctx, rawURL, doc, allowedScopes, allowedGrants)
	if err != nil {
		s.cimdCache.invalidate(rawURL)
		return nil, nil, err
	}

	s.cimdCache.store(rawURL, ttl)

	return cl, nil, nil
}

func (s *Service) parseCIMDClientID(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "https" {
		if !(s.cimdAllowInsecure() && scheme == "http") {
			return nil, fault.New("cimd client_id must use https")
		}
	}
	if u.Host == "" {
		return nil, fault.New("cimd client_id host is required")
	}
	if u.User != nil {
		return nil, fault.New("cimd client_id must not include user info")
	}
	if u.Fragment != "" {
		return nil, fault.New("cimd client_id must not include a fragment")
	}
	if u.Path == "" || u.Path == "/" {
		return nil, fault.New("cimd client_id must contain a path")
	}
	for _, seg := range strings.Split(u.EscapedPath(), "/") {
		if seg == "." || seg == ".." {
			return nil, fault.New("cimd client_id must not contain dot segments")
		}
	}
	// Query strings are allowed (some providers such as ChatGPT encode hints
	// such as ?token_endpoint_auth_method=none directly in the client_id URL).
	// The full URL (including query) is used both for fetching the metadata
	// document and as the canonical client identifier (the document's
	// "client_id" value must match it per the CIMD draft).
	if !s.cimdAllowInsecure() && cimdIsDisallowedHost(u.Hostname()) {
		return nil, fault.New("cimd client_id host is not allowed")
	}

	return u, nil
}

func (s *Service) fetchClientMetadata(ctx context.Context, u *url.URL) (*clientMetadataDocument, time.Duration, *Error) {
	allowInsecure := s.cimdAllowInsecure()

	if !allowInsecure {
		if err := cimdPrevalidateHost(ctx, u.Hostname()); err != nil {
			return nil, 0, oauthError("invalid_client", "CIMD client host is not allowed")
		}
	}

	client := &http.Client{
		Timeout: cimdFetchTimeout,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(dialCtx context.Context, network, address string) (net.Conn, error) {
				host, _, err := net.SplitHostPort(address)
				if err != nil {
					host = address
				}
				if !allowInsecure {
					if err := cimdValidateResolvedHost(dialCtx, host); err != nil {
						return nil, err
					}
				}
				return (&net.Dialer{Timeout: cimdFetchTimeout}).DialContext(dialCtx, network, address)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: allowInsecure},
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return fault.New("cimd metadata fetch must not follow redirects")
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, 0, oauthError("invalid_client", "Failed to create request for CIMD client metadata")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, oauthError("invalid_client", "Failed to fetch CIMD client metadata")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, oauthError("invalid_client", "Failed to fetch CIMD client metadata")
	}
	if !cimdIsJSONContentType(resp.Header.Get("Content-Type")) {
		return nil, 0, oauthError("invalid_client", "CIMD client metadata is not JSON")
	}

	body, err := cimdReadAllBounded(resp.Body, cimdMaxResponseBytes)
	if err != nil {
		return nil, 0, oauthError("invalid_client", "Failed to read CIMD client metadata")
	}

	var doc clientMetadataDocument
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, 0, oauthError("invalid_client", "Failed to parse CIMD client metadata")
	}

	return &doc, cimdCacheTTL(resp.Header), nil
}

func validateClientMetadata(doc *clientMetadataDocument, docURL string) error {
	if strings.TrimSpace(doc.ClientID) != docURL {
		return fault.New("cimd document client_id does not match its URL")
	}
	if len(doc.RedirectURIs) == 0 {
		return fault.New("cimd document must declare at least one redirect_uri")
	}
	for _, ru := range doc.RedirectURIs {
		if err := validateDCRRedirectURI(ru); err != nil {
			return fault.New("cimd document redirect_uri is invalid")
		}
	}
	if doc.TokenEndpointAuthMethod != "" && doc.TokenEndpointAuthMethod != "none" {
		return fault.New("cimd document token_endpoint_auth_method must be none")
	}
	if doc.ClientSecret != "" || doc.ClientSecretExpiresAt != 0 {
		return fault.New("cimd document must not contain client_secret")
	}
	return nil
}

func (s *Service) cimdAllowedScopes(docScope string) []string {
	allow := s.cimdScopeAllowSet()

	requested := splitScope(docScope)
	if len(requested) == 0 {
		// No declared scope: default to the standard OIDC scopes the server allows.
		requested = []string{"openid", "profile", "email", "offline_access"}
	}

	out := []string{}
	seen := map[string]struct{}{}
	for _, sc := range requested {
		if _, ok := allow[sc]; !ok {
			continue
		}
		if _, dup := seen[sc]; dup {
			continue
		}
		seen[sc] = struct{}{}
		out = append(out, sc)
	}

	return out
}

func (s *Service) cimdScopeAllowSet() map[string]struct{} {
	allow := map[string]struct{}{
		"openid":         {},
		"profile":        {},
		"email":          {},
		"offline_access": {},
	}

	configured := s.cfg.OAuthCIMDAllowedScopes
	if len(configured) == 0 {
		configured = cimdDefaultAllowedPermissionScopes
	}

	for _, sc := range configured {
		sc = strings.TrimSpace(sc)
		if sc == "" {
			continue
		}
		if cimdIsPrivilegedScope(sc) && !s.cfg.OAuthCIMDAllowPrivilegedScopes {
			continue
		}
		allow[sc] = struct{}{}
	}

	return allow
}

func cimdIsPrivilegedScope(scope string) bool {
	_, ok := cimdPrivilegedScopes[scope]
	return ok
}

func cimdAllowedGrants(doc *clientMetadataDocument, allowedScopes []string) []string {
	grants := []string{GrantTypeAuthorizationCode}

	if contains(allowedScopes, "offline_access") && cimdDocAllowsRefresh(doc) {
		grants = append(grants, GrantTypeRefreshToken)
	}

	return grants
}

func cimdDocAllowsRefresh(doc *clientMetadataDocument) bool {
	if len(doc.GrantTypes) == 0 {
		return true
	}
	return contains(doc.GrantTypes, GrantTypeRefreshToken)
}

func (s *Service) upsertCIMDClient(ctx context.Context, rawURL string, doc *clientMetadataDocument, scopes, grants []string) (*oauth.Client, error) {
	name := strings.TrimSpace(doc.ClientName)
	if name == "" {
		name = rawURL
	}

	if existing, err := s.clients.GetClientByClientID(ctx, rawURL); err == nil {
		return s.tokens.UpdateClient(ctx, existing.ID, oauth_writer.ClientUpdate{
			Name:          opt.New(name),
			ScopePolicy:   opt.New(oauth.ScopePolicyExplicit),
			RedirectURIs:  opt.New(doc.RedirectURIs),
			AllowedScopes: opt.New(scopes),
			AllowedGrants: opt.New(grants),
		})
	}

	cl, err := s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
		ClientID:      rawURL,
		Name:          name,
		Type:          oauth.ClientTypePublic,
		ScopePolicy:   opt.New(oauth.ScopePolicyExplicit),
		RedirectURIs:  doc.RedirectURIs,
		AllowedScopes: scopes,
		AllowedGrants: grants,
	})
	if err != nil {
		// A concurrent authorize may have created the row first; fall back to it.
		if existing, getErr := s.clients.GetClientByClientID(ctx, rawURL); getErr == nil {
			return existing, nil
		}
		return nil, err
	}

	return cl, nil
}

func cimdReadAllBounded(r io.Reader, maxBytes int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, fault.New("cimd metadata document exceeds maximum size")
	}
	return b, nil
}

func cimdIsJSONContentType(value string) bool {
	if value == "" {
		return false
	}
	mediaType := strings.ToLower(strings.TrimSpace(strings.SplitN(value, ";", 2)[0]))
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}

func cimdCacheTTL(h http.Header) time.Duration {
	cc := strings.ToLower(h.Get("Cache-Control"))
	if strings.Contains(cc, "no-store") || strings.Contains(cc, "no-cache") {
		return 0
	}

	for _, directive := range strings.Split(cc, ",") {
		directive = strings.TrimSpace(directive)
		if !strings.HasPrefix(directive, "max-age=") {
			continue
		}
		seconds, err := strconv.Atoi(strings.TrimPrefix(directive, "max-age="))
		if err != nil || seconds <= 0 {
			return 0
		}
		ttl := time.Duration(seconds) * time.Second
		if ttl > cimdMaxCacheTTL {
			ttl = cimdMaxCacheTTL
		}
		return ttl
	}

	return cimdDefaultCacheTTL
}

func cimdIsDisallowedHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return true
	}
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	return cimdIsDisallowedAddr(addr)
}

func cimdIsDisallowedAddr(addr netip.Addr) bool {
	addr = addr.Unmap()
	if addr.IsLoopback() || addr.IsLinkLocalMulticast() || addr.IsLinkLocalUnicast() || addr.IsMulticast() || addr.IsUnspecified() {
		return true
	}
	for _, prefix := range cimdDisallowedPrefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

func cimdValidateResolvedHost(ctx context.Context, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return fault.New("cimd metadata host is required")
	}
	if cimdIsDisallowedHost(host) {
		return fault.New("cimd metadata host is not allowed")
	}

	if _, err := netip.ParseAddr(host); err == nil {
		return nil
	}

	addrs, err := cimdLookupNetIP(ctx, "ip", host)
	if err != nil {
		return fault.Wrap(err)
	}
	if len(addrs) == 0 {
		return fault.New("cimd metadata host did not resolve to any addresses")
	}
	for _, addr := range addrs {
		if cimdIsDisallowedAddr(addr) {
			return fault.New("cimd metadata host resolves to a disallowed address")
		}
	}

	return nil
}

func cimdPrevalidateHost(ctx context.Context, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return fault.New("cimd metadata host is required")
	}
	if cimdIsDisallowedHost(host) {
		return fault.New("cimd metadata host is not allowed")
	}
	if _, err := netip.ParseAddr(host); err == nil {
		return nil
	}
	addrs, err := cimdLookupNetIP(ctx, "ip", host)
	if err != nil {
		return fault.Wrap(err)
	}
	if len(addrs) == 0 {
		return fault.New("cimd metadata host did not resolve to any addresses")
	}
	for _, addr := range addrs {
		if cimdIsDisallowedAddr(addr) {
			return fault.New("cimd metadata host resolves to a disallowed address")
		}
	}
	return nil
}
