package session_cookie

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/internal/config"
)

func TestSessionCookieTransport(t *testing.T) {
	for _, tt := range []struct {
		name       string
		api        string
		web        string
		secure     bool
		webCookies int
	}{
		{"https", "https://example.com", "https://example.com", true, 1},
		{"http", "http://example.com", "http://example.com", false, 1},
		{"http_api", "http://example.com:8000", "https://example.com", false, 1},
		{"http_web", "https://example.com", "http://example.com:3000", true, 0},
		{"localhost", "http://localhost:8000", "http://localhost:3000", false, 1},
		{"lan_ip", "http://192.168.1.10:8000", "http://192.168.1.10:3000", false, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			api, web := u(tt.api), u(tt.web)
			jar, err := New(config.Config{PublicAPIAddress: api, PublicWebAddress: web}, nil)
			r.NoError(err)
			sessionToken, err := token.Generate()
			r.NoError(err)

			response := httptest.NewRecorder()
			http.SetCookie(response, jar.Create(sessionToken))
			cookies := response.Result().Cookies()
			r.Len(cookies, 1)
			r.Equal(tt.secure, cookies[0].Secure)
			r.True(cookies[0].HttpOnly)

			browser, err := cookiejar.New(nil)
			r.NoError(err)
			browser.SetCookies(&api, cookies)
			r.Len(browser.Cookies(&api), 1)
			r.Len(browser.Cookies(&web), tt.webCookies)
			if tt.secure {
				insecure := api
				insecure.Scheme = "http"
				r.Empty(browser.Cookies(&insecure))
			}

			logout := jar.Destroy()
			r.Equal(tt.secure, logout.Secure)
			browser.SetCookies(&api, []*http.Cookie{logout})
			r.Empty(browser.Cookies(&api))
			r.Empty(browser.Cookies(&web))
		})
	}
}
