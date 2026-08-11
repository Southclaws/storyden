package oauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"log/slog"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_querier"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeDeviceCode        = "urn:ietf:params:oauth:grant-type:device_code"

	CodeChallengeMethodS256 = "S256"

	StorydenCLIClientID = "storyden-cli"

	cleanupInterval          = time.Hour
	dcrClientRetentionPeriod = 7 * 24 * time.Hour

	// rotation keeps every superseded row so revokeRefreshTokenFamily can walk
	// the chain, they are dropped once they have been unusable for this long
	refreshTokenRetentionPeriod = 7 * 24 * time.Hour
	maxUserCodeInputLength      = 32
)

type Error struct {
	Code        string
	Description string
}

type Service struct {
	cfg       config.Config
	clients   *oauth_querier.Querier
	tokens    *oauth_writer.Writer
	account   *account_querier.Querier
	signer    *rsa.PrivateKey
	kid       string
	issuer    string
	cimdCache *cimdCache
}

func (s *Service) Enabled() bool {
	return s.cfg.OAuthEnabled
}

func canAuthoriseOAuthClients(permissions rbac.Permissions) bool {
	return permissions.HasAny(rbac.PermissionUseOauthClients, rbac.PermissionAdministrator)
}

func New(
	lc fx.Lifecycle,
	logger *slog.Logger,
	cfg config.Config,
	clients *oauth_querier.Querier,
	tokens *oauth_writer.Writer,
	account *account_querier.Querier,
) (*Service, error) {
	issuer := strings.TrimSuffix(cfg.PublicAPIAddress.String(), "/")

	if !cfg.OAuthEnabled {
		service := &Service{
			cfg:       cfg,
			clients:   clients,
			tokens:    tokens,
			account:   account,
			issuer:    issuer,
			cimdCache: newCIMDCache(),
		}
		service.registerCleanupJob(lc, logger)

		return service, nil
	}

	b, err := base64.StdEncoding.DecodeString(cfg.OAuthSigningKeyBase64)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fault.New("invalid oauth private key pem")
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		rsaPK, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, err
		}
		parsed = rsaPK
	}

	pk, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fault.New("only RSA keys supported")
	}

	kid := cfg.OAuthSigningKeyID
	if kid == "" {
		h := sha256.Sum256(x509.MarshalPKCS1PublicKey(&pk.PublicKey))
		kid = hex.EncodeToString(h[:8])
	}

	service := &Service{
		cfg:       cfg,
		clients:   clients,
		tokens:    tokens,
		account:   account,
		signer:    pk,
		kid:       kid,
		issuer:    issuer,
		cimdCache: newCIMDCache(),
	}
	service.registerCleanupJob(lc, logger)

	return service, nil
}

func (s *Service) registerCleanupJob(lc fx.Lifecycle, logger *slog.Logger) {
	var cancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ctx, stop := context.WithCancel(context.Background())
			cancel = stop

			go s.cleanupExpiredRecordsLoop(ctx, logger)

			return nil
		},
		OnStop: func(context.Context) error {
			if cancel != nil {
				cancel()
			}

			return nil
		},
	})
}

func (s *Service) cleanupExpiredRecordsLoop(ctx context.Context, logger *slog.Logger) {
	s.cleanupExpiredRecords(ctx, logger)

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupExpiredRecords(ctx, logger)
		}
	}
}

func (s *Service) cleanupExpiredRecords(ctx context.Context, logger *slog.Logger) {
	now := time.Now()

	deviceAuthorisations, err := s.tokens.DeleteExpiredDeviceAuthorisations(ctx, now)
	if err != nil {
		logger.Error("failed to clean expired oauth device authorizations", slog.Any("error", err))
		return
	}

	authorizationRequests, err := s.tokens.DeleteExpiredAuthorisationRequests(ctx, now)
	if err != nil {
		logger.Error("failed to clean expired oauth authorization requests", slog.Any("error", err))
		return
	}

	unusedDCRClients, err := s.tokens.DeleteUnusedDCRClients(ctx, now.Add(-dcrClientRetentionPeriod))
	if err != nil {
		logger.Error("failed to clean unused dcr clients", slog.Any("error", err))
		return
	}

	refreshTokens, err := s.tokens.DeleteExpiredRefreshTokens(ctx, now.Add(-refreshTokenRetentionPeriod))
	if err != nil {
		logger.Error("failed to clean expired oauth refresh tokens", slog.Any("error", err))
		return
	}

	if deviceAuthorisations > 0 || authorizationRequests > 0 || unusedDCRClients > 0 || refreshTokens > 0 {
		logger.Debug(
			"cleaned expired oauth records",
			slog.Int("device_authorizations", deviceAuthorisations),
			slog.Int("authorization_requests", authorizationRequests),
			slog.Int("unused_dcr_clients", unusedDCRClients),
			slog.Int("refresh_tokens", refreshTokens),
		)
	}
}

func oauthError(code string, description string) *Error {
	return &Error{Code: code, Description: description}
}

func (s *Service) deviceAuthorizationConsentURL(userCode string) string {
	base := s.cfg.OAuthDeviceAuthorisationConsentURL
	if base.String() == "" {
		base = s.cfg.PublicWebAddress
		base.Path = strings.TrimRight(base.Path, "/") + "/oauth/consent"
	}

	u := base
	q := u.Query()
	if userCode != "" {
		q.Set("user_code", userCode)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (s *Service) authorizationCodeConsentURL(requestID string) string {
	base := s.cfg.OAuthAuthorisationCodeConsentURL
	if base.String() == "" {
		base = s.cfg.PublicWebAddress
		base.Path = strings.TrimRight(base.Path, "/") + "/oauth/authorize/consent"
	}

	u := base
	q := u.Query()
	q.Set("request_id", requestID)
	u.RawQuery = q.Encode()

	return u.String()
}

func (s *Service) LoginURL() string {
	base := s.cfg.OAuthAuthorisationLoginURL
	if !base.IsAbs() || base.Host == "" {
		base = s.cfg.PublicWebAddress
		base.Path = strings.TrimRight(base.Path, "/") + "/login"
	}

	return base.String()
}

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return b64url(b), nil
}

// userCodeAlphabet is Crockford's base32 alphabet. It excludes the visually
// ambiguous letters I, L, O, and U; normalizeCode accepts I/L and O as aliases
// for the digits 1 and 0 when a user transcribes a code.
const userCodeAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

var userCodeAliases = strings.NewReplacer(
	"O", "0",
	"I", "1",
	"L", "1",
)

// generateUserCode produces an 8-character user code formatted as XXXX-XXXX
// from userCodeAlphabet, providing exactly 40 bits of entropy. 256 is an exact
// multiple of len(userCodeAlphabet) (32), so masking a random byte to 5 bits
// introduces no modulo bias.
func generateUserCode() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	out := make([]byte, 8)
	for i, v := range b {
		out[i] = userCodeAlphabet[v&0x1F]
	}

	return string(out[:4]) + "-" + string(out[4:]), nil
}

func hashString(v string) string {
	s := sha256.Sum256([]byte(v))
	return hex.EncodeToString(s[:])
}

func splitScope(v string) []string {
	return strings.Fields(strings.TrimSpace(v))
}

func contains(in []string, v string) bool {
	for _, it := range in {
		if it == v {
			return true
		}
	}
	return false
}

func normalizeCode(v string) string {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(v), "-", ""))

	return userCodeAliases.Replace(normalized)
}

func parseUserCode(v string) (string, bool) {
	// Bound work before trimming, case folding, or removing separators. The
	// canonical representation is only eight bytes; the extra room permits
	// harmless surrounding whitespace and separators without accepting an
	// arbitrarily large value at this low-entropy lookup boundary.
	if len(v) > maxUserCodeInputLength {
		return "", false
	}

	normalized := normalizeCode(v)
	if len(normalized) != 8 {
		return "", false
	}

	for _, c := range normalized {
		if !strings.ContainsRune(userCodeAlphabet, c) {
			return "", false
		}
	}

	return normalized, true
}
