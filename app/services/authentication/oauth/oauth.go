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

	cleanupInterval = time.Hour
)

type Error struct {
	Code        string
	Description string
}

type Service struct {
	cfg     config.Config
	clients *oauth_querier.Querier
	tokens  *oauth_writer.Writer
	account *account_querier.Querier
	signer  *rsa.PrivateKey
	kid     string
	issuer  string
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
			cfg:     cfg,
			clients: clients,
			tokens:  tokens,
			account: account,
			issuer:  issuer,
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
		cfg:     cfg,
		clients: clients,
		tokens:  tokens,
		account: account,
		signer:  pk,
		kid:     kid,
		issuer:  issuer,
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

	DeviceAuthorisations, err := s.tokens.DeleteExpiredDeviceAuthorisations(ctx, now)
	if err != nil {
		logger.Error("failed to clean expired oauth device authorizations", slog.Any("error", err))
		return
	}

	authorizationRequests, err := s.tokens.DeleteExpiredAuthorisationRequests(ctx, now)
	if err != nil {
		logger.Error("failed to clean expired oauth authorization requests", slog.Any("error", err))
		return
	}

	if DeviceAuthorisations > 0 || authorizationRequests > 0 {
		logger.Debug(
			"cleaned expired oauth records",
			slog.Int("device_authorizations", DeviceAuthorisations),
			slog.Int("authorization_requests", authorizationRequests),
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
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(v), "-", ""))
}
