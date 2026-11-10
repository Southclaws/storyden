package link_url

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{" HTTPS://Example.COM ", "https://example.com/"},
		{"http://EXAMPLE.com", "http://example.com/"},
		{"https://www.example.com", "https://www.example.com/"},
		{"https://example.com/Watch?v=2#chapter", "https://example.com/Watch?v=2#chapter"},
		{"https://example.com/a%2Fb?q=a%2Fb", "https://example.com/a%2Fb?q=a%2Fb"},
		{"https://example.com./", "https://example.com/"},
		{"https://bücher.de/", "https://xn--bcher-kva.de/"},
		{"https://例え.テスト/", "https://xn--r8jz45g.xn--zckzah/"},
		{"https://a1-b2.example.com/", "https://a1-b2.example.com/"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			u, err := Parse(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, u.String())
		})
	}
}

func TestParseRejectsInvalidAndUnsafeURLs(t *testing.T) {
	for _, raw := range []string{
		"", "example.com", "/relative", "//example.com", "ftp://example.com/",
		"https://example.com/%zz", "https://example.com:abc/", "javascript:alert(1)", "https:", "https://", "https:///path",
		"https://user:password@example.com/", "https://example.com:0/", "https://example.com:65536/",
		"https://example.com:443/", "http://example.com:80/", "https://example.com:8443/", "https://example.com:/",
		"https://8.8.8.8/", "https://8.8.8.8./", "https://[2606:4700:4700::1111]/",
		"https://2606:4700:4700::1111/", "https://2130706433/", "https://127.1/",
		"https://0x7f000001/", "https://0177.0.0.1/", "https://8.8.8.0x8/", "https://８.８.８.８/",
		"https://intranet/", "https://example.123/", "https://.example.com/", "https://example..com/",
		"https://example.com../", "https://-example.com/", "https://example-.com/", "https://under_score.example.com/",
		"https://" + strings.Repeat("a", 64) + ".com/",
		"https://" + strings.Repeat(strings.Repeat("a", 63)+".", 4) + "com/",
		"http://127.0.0.1/", "http://localhost./", "http://foo.localhost/",
		"http://169.254.169.254/", "http://10.0.0.1/", "http://[::1]/",
		"http://[::ffff:127.0.0.1]/", "http://[fe80::1%25en0]/",
	} {
		t.Run(raw, func(t *testing.T) {
			_, err := Parse(raw)
			require.Error(t, err)
			require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
			require.NotEmpty(t, fmsg.GetIssue(err), "URL errors must include a message for the caller")
		})
	}
}

func TestURLNetURLReturnsIndependentCopy(t *testing.T) {
	u, err := Parse("https://example.com/articles?q=one#section")
	require.NoError(t, err)

	raw := u.NetURL()
	raw.Scheme = "file"
	raw.Host = "127.0.0.1:8080"
	raw.Path = "/other"

	require.Equal(t, "https://example.com/articles?q=one#section", u.String())
	require.Equal(t, "example.com", u.Hostname())
}

func TestURLJSONValidatesBoundary(t *testing.T) {
	var u URL
	require.NoError(t, json.Unmarshal([]byte(`"HTTPS://Example.COM"`), &u))
	require.Equal(t, "https://example.com/", u.String())
	encoded, err := json.Marshal(u)
	require.NoError(t, err)
	require.JSONEq(t, `"https://example.com/"`, string(encoded))

	for _, raw := range []string{`"http://127.0.0.1/"`, `"https://example.com:443/"`, `"ftp://example.com/"`, `""`, `null`, `{}`, `42`} {
		t.Run(raw, func(t *testing.T) {
			var invalidURL URL
			require.Error(t, json.Unmarshal([]byte(raw), &invalidURL))
			require.True(t, invalidURL.IsZero())
		})
	}
}

func TestURLZeroValueCannotBeSerialized(t *testing.T) {
	var u URL
	require.True(t, u.IsZero())
	_, err := json.Marshal(u)
	require.ErrorIs(t, err, ErrEmpty)
}
