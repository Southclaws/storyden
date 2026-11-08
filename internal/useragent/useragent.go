// Package useragent produces descriptive browser and OS labels, not trusted device identities.
// Tokenization and detection precedence are adapted from mileusna/useragent v1.3.5.
// See LICENSE and README.md for attribution and intentional differences.
package useragent

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxHeaderBytes = 4096

type UserAgent struct {
	Name string
	OS   string
}

type browserRule struct {
	token           string
	name            string
	requiresVersion bool
}

var browserRules = [...]browserRule{
	{"Googlebot", "Googlebot", false},
	{"GoogleProber", "GoogleProber", false},
	{"GoogleProducer", "GoogleProducer", false},
	{"Bytespider", "Bytespider", false},
	{"Applebot", "Applebot", false},
	{"Opera Mini", "Opera Mini", true},
	{"OPR", "Opera", true},
	{"OPT", "Opera Touch", true},
	{"OPiOS", "Opera", true},
	{"CriOS", "Chrome", true},
	{"FxiOS", "Firefox", true},
	{"Firefox", "Firefox", true},
	{"Vivaldi", "Vivaldi", true},
	{"MSIE", "Internet Explorer", false},
	{"EdgiOS", "Edge", true},
	{"Edge", "Edge", true},
	{"Edg", "Edge", true},
	{"EdgA", "Edge", true},
	{"bingbot", "Bingbot", true},
	{"YandexBot", "YandexBot", true},
	{"YandexAdNet", "YandexAdNet", true},
	{"SamsungBrowser", "Samsung Browser", true},
	{"HeadlessChrome", "Headless Chrome", true},
	{"AdsBot-Google-Mobile", "Google Ads Bot", false},
	{"Mediapartners-Google", "Google Ads Bot", false},
	{"AdsBot-Google", "Google Ads Bot", false},
	{"Yahoo Ad monitoring", "Yahoo Ad monitoring", false},
	{"MiuiBrowser", "Miui Browser", true},
	{"FBAN", "Facebook App", false},
	{"FB_IAB", "Facebook App", false},
	{"Instagram", "Instagram App", false},
	{"BytedanceWebview", "TikTok App", false},
	{"HuaweiBrowser", "Huawei Browser", true},
	{"BlackBerry", "BlackBerry", false},
	{"NetFront", "NetFront", false},
	{"QtWebEngine", "QtWebEngine", true},
	{"ArkWeb", "ArkWeb", true},
	{"Mobile DuckDuckGo", "Mobile DuckDuckGo", true},
	{"Chrome", "Chrome", false},
	{"Brave Chrome", "Chrome", false},
	{"Safari", "Safari", false},
}

var operatingSystems = [...]struct {
	token string
	name  string
}{
	{"Android", "Android"},
	{"iPhone", "iOS"},
	{"iPad", "iOS"},
	{"iPod", "iOS"},
	{"Windows NT", "Windows"},
	{"Windows Phone OS", "Windows Phone"},
	{"Macintosh", "macOS"},
	{"Linux", "Linux"},
	{"FreeBSD", "FreeBSD"},
	{"CrOS", "ChromeOS"},
	{"BlackBerry", "BlackBerry"},
	{"OpenHarmony", "Harmony"},
}

func Parse(header string) UserAgent {
	ua := UserAgent{Name: "Unknown", OS: "Unknown"}
	if !validHeader(header) {
		return ua
	}

	tokens := parseTokens(strings.ReplaceAll(header, "\t", " "))
	ua.OS = detectOS(tokens)
	ua.Name = detectBrowser(tokens, ua.OS)

	switch ua.Name {
	case "Applebot":
		ua.OS = "Unknown"
	case "Samsung Browser":
		ua.OS = "Android"
	}

	return ua
}

func validHeader(header string) bool {
	if len(header) > maxHeaderBytes || !utf8.ValidString(header) {
		return false
	}

	for _, c := range header {
		if c == '\t' {
			continue
		}

		isControl := unicode.IsControl(c) || unicode.Is(unicode.Cf, c)
		isLineSeparator := c == '\u2028' || c == '\u2029'
		if isControl || isLineSeparator {
			return false
		}
	}

	return true
}

func detectOS(tokens properties) string {
	for _, rule := range operatingSystems {
		if tokens.has(rule.token) {
			return rule.name
		}
	}

	return "Unknown"
}

func detectBrowser(tokens properties, os string) string {
	for _, rule := range browserRules {
		value, exists := tokens[rule.token]
		if !exists {
			continue
		}

		if rule.requiresVersion && value == "" {
			continue
		}

		return rule.name
	}

	if os == "Android" && tokens["Version"] != "" {
		return "Android browser"
	}

	return "Unknown"
}

type properties map[string]string

func (p properties) has(key string) bool {
	_, exists := p[key]
	return exists
}

func parseTokens(header string) properties {
	tokens := make(properties)
	var key bytes.Buffer
	var value bytes.Buffer
	readingValue := false
	isURL := false

	addToken := func() {
		name := strings.TrimSpace(key.String())
		version := strings.TrimSpace(value.String())
		if name != "" && !isURL && !ignore(name) {
			if version == "" {
				name, version = splitVersion(name)
			}

			if name == "XiaoMi" && strings.HasPrefix(version, "MiuiBrowser/") {
				name = "MiuiBrowser"
				version = strings.TrimPrefix(version, "MiuiBrowser/")
			}

			if _, exists := tokens[name]; !exists {
				tokens[name] = version
			}
		}

		key.Reset()
		value.Reset()
		readingValue = false
		isURL = false
	}

	for i := 0; i < len(header); i++ {
		c := header[i]
		switch {
		case c == '(' || c == ')' || c == '[' || c == ']' || c == ';':
			addToken()

		case c == ':':
			if bytes.HasSuffix(key.Bytes(), []byte("http")) || bytes.HasSuffix(key.Bytes(), []byte("https")) {
				key.WriteByte(c)
			} else if i+1 < len(header) && header[i+1] != ' ' {
				key.WriteByte(' ')
			}

		case readingValue && c == ' ':
			addToken()

		case readingValue:
			value.WriteByte(c)

		case c == '/' && !isURL:
			nextIsSlash := i+1 < len(header) && header[i+1] == '/'
			hasURLScheme := bytes.HasSuffix(key.Bytes(), []byte("http:")) || bytes.HasSuffix(key.Bytes(), []byte("https:"))
			if nextIsSlash && hasURLScheme {
				key.WriteByte(c)
				isURL = true
			} else if ignore(key.String()) {
				key.Reset()
			} else {
				readingValue = true
			}

		default:
			key.WriteByte(c)
		}
	}

	addToken()
	return tokens
}

func splitVersion(token string) (string, string) {
	i := strings.LastIndexByte(token, ' ')
	if i == -1 {
		return token, ""
	}

	switch token[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android", "OpenHarmony", "Instagram":
		return token[:i], token[i+1:]
	case "CrOS x86_64", "CrOS aarch64", "CrOS armv7l":
		return "CrOS", token[i+1:]
	default:
		return token, ""
	}
}

func ignore(token string) bool {
	switch token {
	case "KHTML, like Gecko", "U", "compatible", "Mozilla", "WOW64", "en", "en-us", "en-gb", "ru-ru", "Browser":
		return true
	default:
		return false
	}
}
