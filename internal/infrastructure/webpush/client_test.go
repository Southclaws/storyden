package webpush

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func pushConfig(t *testing.T) config.Config {
	t.Helper()

	r := require.New(t)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	r.NoError(err)

	private, err := key.Bytes()
	r.NoError(err)

	public, err := key.PublicKey.Bytes()
	r.NoError(err)

	return config.Config{
		WebPushEnabled:          true,
		WebPushMaxSubscriptions: 5,
		WebPushVAPIDPrivateKey:  base64.RawURLEncoding.EncodeToString(private),
		WebPushVAPIDPublicKey:   base64.RawURLEncoding.EncodeToString(public),
		WebPushVAPIDSubject:     "mailto:admin@example.com",
	}
}

func TestConfiguration(t *testing.T) {
	valid := pushConfig(t)
	for _, tc := range []struct {
		name   string
		change func(*config.Config)
		valid  bool
	}{
		{
			name: "valid",
			change: func(*config.Config) {
			},
			valid: true,
		},
		{
			name: "https_contact",
			change: func(c *config.Config) {
				c.WebPushVAPIDSubject = "https://example.com/contact"
			},
			valid: true,
		},
		{
			name: "missing_private",
			change: func(c *config.Config) {
				c.WebPushVAPIDPrivateKey = ""
			},
			valid: false,
		},
		{
			name: "invalid_private",
			change: func(c *config.Config) {
				c.WebPushVAPIDPrivateKey = "bad!"
			},
			valid: false,
		},
		{
			name: "zero_private",
			change: func(c *config.Config) {
				c.WebPushVAPIDPrivateKey = strings.Repeat("A", 43)
			},
			valid: false,
		},
		{
			name: "missing_public",
			change: func(c *config.Config) {
				c.WebPushVAPIDPublicKey = ""
			},
			valid: false,
		},
		{
			name: "invalid_public",
			change: func(c *config.Config) {
				c.WebPushVAPIDPublicKey = "bad!"
			},
			valid: false,
		},
		{
			name: "mismatched_public",
			change: func(c *config.Config) {
				c.WebPushVAPIDPublicKey = pushConfig(t).WebPushVAPIDPublicKey
			},
			valid: false,
		},
		{
			name: "missing_subject",
			change: func(c *config.Config) {
				c.WebPushVAPIDSubject = ""
			},
			valid: false,
		},
		{
			name: "invalid_contact",
			change: func(c *config.Config) {
				c.WebPushVAPIDSubject = "mailto:"
			},
			valid: false,
		},
		{
			name: "insecure_contact",
			change: func(c *config.Config) {
				c.WebPushVAPIDSubject = "http://example.com"
			},
			valid: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			cfg := valid
			tc.change(&cfg)
			_, err := New(cfg)
			if tc.valid {
				r.NoError(err)
			} else {
				r.Error(err)
			}

			cfg.WebPushEnabled = false
			disabled, err := New(cfg)
			r.NoError(err, "disabled push must ignore VAPID configuration")
			a.False(disabled.Enabled())
			a.ErrorIs(disabled.Send(context.Background(), Subscription{}, nil), ErrDisabled)
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestEncryptedRequest(t *testing.T) {
	cfg := pushConfig(t)

	for _, tc := range []struct {
		name     string
		endpoint string
		audience string
	}{
		{
			name:     "hostname",
			endpoint: "https://push.example/messages/secret?token=secret",
			audience: "https://push.example",
		},
		{
			name:     "default_port",
			endpoint: "https://PUSH.example:443/messages/secret?token=secret",
			audience: "https://push.example",
		},
		{
			name:     "non_default_port",
			endpoint: "https://push.example:8443/messages/secret?token=secret",
			audience: "https://push.example:8443",
		},
		{
			name:     "ipv6",
			endpoint: "https://[2001:4860:4860::8888]/messages/secret?token=secret",
			audience: "https://[2001:4860:4860::8888]",
		},
		{
			name:     "ipv6_default_port",
			endpoint: "https://[2001:4860:4860::8888]:443/messages/secret?token=secret",
			audience: "https://[2001:4860:4860::8888]",
		},
		{
			name:     "ipv6_non_default_port",
			endpoint: "https://[2001:4860:4860::8888]:8443/messages/secret?token=secret",
			audience: "https://[2001:4860:4860::8888]:8443",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			client, err := New(cfg)
			r.NoError(err)

			input, receiver := NewTestSubscription(t, tc.endpoint)
			subscription, err := ParseSubscription(input)
			r.NoError(err)

			client.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				a.Equal(http.MethodPost, req.Method)
				a.Equal(tc.endpoint, req.URL.String())
				a.Equal(int64(123), req.ContentLength)
				a.Equal("aes128gcm", req.Header.Get("Content-Encoding"))
				a.Equal("application/octet-stream", req.Header.Get("Content-Type"))
				a.Equal("86400", req.Header.Get("TTL"))

				body, err := io.ReadAll(req.Body)
				r.NoError(err)

				plaintext, err := DecryptTestPayload(t, receiver, subscription.Auth, body)
				r.NoError(err)
				a.Equal("0123456789abcdefghij\x02", string(plaintext))

				auth, ok := strings.CutPrefix(req.Header.Get("Authorization"), "vapid t=")
				r.True(ok)
				token, public, ok := strings.Cut(auth, ", k=")
				r.True(ok)
				a.Equal(cfg.WebPushVAPIDPublicKey, public)

				parsed, err := jwt.Parse(token, func(*jwt.Token) (any, error) {
					return &client.privateKey.PublicKey, nil
				},
					jwt.WithValidMethods([]string{"ES256"}),
					jwt.WithAudience(tc.audience),
					jwt.WithSubject("mailto:admin@example.com"),
				)
				r.NoError(err)

				expires, err := parsed.Claims.GetExpirationTime()
				r.NoError(err)
				a.WithinDuration(time.Now().Add(12*time.Hour), expires.Time, 5*time.Second)

				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			})

			r.NoError(client.Send(context.Background(), *subscription, []byte("0123456789abcdefghij")))
		})
	}
}

func TestPushResponses(t *testing.T) {
	r := require.New(t)

	cfg := pushConfig(t)
	client, err := New(cfg)
	r.NoError(err)

	for _, status := range []int{201, 202, 404, 410, 401, 403, 429, 500, 307} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Location": {"https://other.example/"}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			})
			input, _ := NewTestSubscription(t, "https://push.example/message")
			subscription, err := ParseSubscription(input)
			r.NoError(err)

			err = client.Send(context.Background(), *subscription, []byte("0123456789abcdefghij"))
			switch status {
			case 201, 202:
				r.NoError(err)
			case 404, 410:
				r.ErrorIs(err, ErrExpired)
			default:
				r.Error(err)
				a.False(errors.Is(err, ErrExpired))
			}
		})
	}
}

func TestSubscriptionValidation(t *testing.T) {
	public := pushConfig(t).WebPushVAPIDPublicKey

	for _, tc := range []struct {
		endpoint string
		public   string
		auth     string
		valid    bool
	}{
		{
			endpoint: "https://push.example/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    true,
		},
		{
			endpoint: "http://push.example/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://127.0.0.1/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://[::1]/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://10.0.0.1/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://user:password@push.example/token",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://push.example/token#fragment",
			public:   public,
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://push.example/token",
			public:   "bad!",
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://push.example/token",
			public:   "AAAA",
			auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			valid:    false,
		},
		{
			endpoint: "https://push.example/token",
			public:   public,
			auth:     "AAAA",
			valid:    false,
		},
	} {
		_, err := ParseSubscription(SubscriptionInput{
			Endpoint: tc.endpoint,
			P256DH:   tc.public,
			Auth:     tc.auth,
		})
		assert.Equal(t, tc.valid, err == nil, "%+v: %v", tc, err)
	}
}
