package oauth_http_client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"golang.org/x/oauth2"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			if err := ValidateRemoteOAuthURL(req.URL.String(), "redirect URL"); err != nil {
				return err
			}
			return nil
		},
	}
}

func ContextWithHTTPClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, client)
}

func ValidateRemoteOAuthURL(raw string, label string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return fmt.Errorf("%s is invalid: %w", label, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("%s must be an absolute URL", label)
	}
	if u.User != nil {
		return fmt.Errorf("%s must not include user info", label)
	}
	if u.Fragment != "" {
		return fmt.Errorf("%s must not include a fragment", label)
	}
	switch u.Scheme {
	case "https":
		return nil
	case "http":
		if isLoopbackHost(u.Hostname()) {
			return nil
		}
		return fmt.Errorf("%s must use https", label)
	default:
		return fmt.Errorf("%s must use http or https", label)
	}
}

func Endpoint(connection oauth_remote.Connection) oauth2.Endpoint {
	endpoint := oauth2.Endpoint{
		AuthURL:  connection.AuthorizationEndpoint,
		TokenURL: connection.TokenEndpoint,
	}
	switch connection.TokenEndpointAuthMethod {
	case "client_secret_basic":
		endpoint.AuthStyle = oauth2.AuthStyleInHeader
	case "client_secret_post":
		endpoint.AuthStyle = oauth2.AuthStyleInParams
	case "none":
		endpoint.AuthStyle = oauth2.AuthStyleInParams
	}
	return endpoint
}

func SupportsTokenEndpointAuth(methods []string, method string) bool {
	return len(methods) == 0 || Contains(methods, method)
}

func SupportedTokenEndpointAuthMethod(method string) bool {
	switch method {
	case "none", "client_secret_basic", "client_secret_post":
		return true
	default:
		return false
	}
}

func Contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func SplitScope(scope string) []string {
	return strings.Fields(scope)
}

func StringExtra(token *oauth2.Token, key string) string {
	v, _ := token.Extra(key).(string)
	return v
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
