package mark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cyrillic",
			input:    "Документация",
			expected: "документация",
		},
		{
			name:     "japanese with spaces",
			input:    "日本語 テスト",
			expected: "日本語-テスト",
		},
		{
			name:     "greek",
			input:    "Παράδειγμα",
			expected: "παράδειγμα",
		},
		{
			name:     "hindi",
			input:    "परीक्षण दस्तावेज़",
			expected: "परीक्षण-दस्तावेज़",
		},
		{
			name:     "korean with full-width space",
			input:    "문서　테스트",
			expected: "문서-테스트",
		},
		{
			name:     "hebrew with hyphen",
			input:    "תיעוד-מערכת",
			expected: "תיעוד-מערכת",
		},
		{
			name:     "persian with hyphen",
			input:    "مثالِ-سادِه",
			expected: "مثالِ-سادِه",
		},
		{
			name:     "basic english",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "uppercase to lowercase",
			input:    "HELLO WORLD",
			expected: "hello-world",
		},
		{
			name:     "mixed case",
			input:    "HeLLo WoRLd",
			expected: "hello-world",
		},
		{
			name:     "leading spaces",
			input:    "   hello world",
			expected: "hello-world",
		},
		{
			name:     "trailing spaces",
			input:    "hello world   ",
			expected: "hello-world",
		},
		{
			name:     "leading and trailing spaces",
			input:    "   hello world   ",
			expected: "hello-world",
		},
		{
			name:     "multiple spaces",
			input:    "hello    world",
			expected: "hello-world",
		},
		{
			name:     "multiple hyphens",
			input:    "hello----world",
			expected: "hello-world",
		},
		{
			name:     "leading hyphens",
			input:    "---hello-world",
			expected: "hello-world",
		},
		{
			name:     "trailing hyphens",
			input:    "hello-world---",
			expected: "hello-world",
		},
		{
			name:     "leading underscores",
			input:    "___hello_world",
			expected: "hello_world",
		},
		{
			name:     "trailing underscores",
			input:    "hello_world___",
			expected: "hello_world",
		},
		{
			name:     "special characters",
			input:    "hello@world!test#123",
			expected: "hello-world-test-123",
		},
		{
			name:     "punctuation",
			input:    "hello, world. test?",
			expected: "hello-world-test",
		},
		{
			name:     "brackets and parens",
			input:    "hello (world) [test]",
			expected: "hello-world-test",
		},
		{
			name:     "emojis",
			input:    "hello 👋 world 🌍",
			expected: "hello-world",
		},
		{
			name:     "mixed emojis and text",
			input:    "🎉 Party Time 🎊",
			expected: "party-time",
		},
		{
			name:     "numbers",
			input:    "test 123 456",
			expected: "test-123-456",
		},
		{
			name:     "numbers with letters",
			input:    "test123abc456",
			expected: "test123abc456",
		},
		{
			name:     "accented characters",
			input:    "café résumé",
			expected: "café-résumé",
		},
		{
			name:     "german umlauts",
			input:    "Über Größe",
			expected: "über-größe",
		},
		{
			name:     "full-width characters",
			input:    "ｈｅｌｌｏ　ｗｏｒｌｄ",
			expected: "hello-world",
		},
		{
			name:     "mixed full-width and half-width",
			input:    "hello ｗｏｒｌｄ",
			expected: "hello-world",
		},
		{
			name:     "tabs",
			input:    "hello\tworld",
			expected: "hello-world",
		},
		{
			name:     "newlines",
			input:    "hello\nworld",
			expected: "hello-world",
		},
		{
			name:     "carriage returns",
			input:    "hello\rworld",
			expected: "hello-world",
		},
		{
			name:     "mixed whitespace",
			input:    "hello \t\n\r world",
			expected: "hello-world",
		},
		{
			name:     "zero width space",
			input:    "hello\u200Bworld",
			expected: "hello-world",
		},
		{
			name:     "zero width non-joiner",
			input:    "hello\u200Cworld",
			expected: "hello-world",
		},
		{
			name:     "zero width joiner",
			input:    "hello\u200Dworld",
			expected: "hello-world",
		},
		{
			name:     "soft hyphen",
			input:    "hello\u00ADworld",
			expected: "hello-world",
		},
		{
			name:     "non-breaking space",
			input:    "hello\u00A0world",
			expected: "hello-world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "only hyphens",
			input:    "-----",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "only emojis",
			input:    "👋🌍🎉",
			expected: "",
		},
		{
			name:     "url with protocol",
			input:    "https://example.com",
			expected: "https-example-com",
		},
		{
			name:     "email address",
			input:    "user@example.com",
			expected: "user-example-com",
		},
		{
			name:     "path-like string",
			input:    "path/to/file.txt",
			expected: "path-to-file-txt",
		},
		{
			name:     "mixed scripts",
			input:    "English 日本語 Русский",
			expected: "english-日本語-русский",
		},
		{
			name:     "right-to-left scripts",
			input:    "العربية עברית",
			expected: "العربية-עברית",
		},
		{
			name:     "chinese characters",
			input:    "中文测试",
			expected: "中文测试",
		},
		{
			name:     "chinese with spaces",
			input:    "中文 测试",
			expected: "中文-测试",
		},
		{
			name:     "thai",
			input:    "ทดสอบ ภาษาไทย",
			expected: "ทดสอบ-ภาษาไทย",
		},
		{
			name:     "quotes and apostrophes",
			input:    "it's a \"test\"",
			expected: "it-s-a-test",
		},
		{
			name:     "slashes and backslashes",
			input:    "test/slash\\backslash",
			expected: "test-slash-backslash",
		},
		{
			name:     "currency symbols",
			input:    "$100 €50 ¥1000",
			expected: "100-50-1000",
		},
		{
			name:     "math symbols",
			input:    "x + y = z",
			expected: "x-y-z",
		},
		{
			name:     "already valid slug",
			input:    "hello-world",
			expected: "hello-world",
		},
		{
			name:     "underscores preserved",
			input:    "hello_world_test",
			expected: "hello_world_test",
		},
		{
			name:     "mixed hyphens and underscores",
			input:    "hello-world_test",
			expected: "hello-world_test",
		},
		{
			name:     "control characters",
			input:    "hello\x00\x01\x02world",
			expected: "hello-world",
		},
		{
			name:     "ligatures",
			input:    "ﬁle ﬂag",
			expected: "file-flag",
		},
		{
			name:     "superscripts and subscripts",
			input:    "x² + y₃",
			expected: "x2-y3",
		},
		{
			name:     "fractions",
			input:    "½ + ¼",
			expected: "1-2-1-4",
		},
		{
			name:     "combining diacritics",
			input:    "e\u0301" + "cole",
			expected: "école",
		},
		{
			name:     "arabic numerals in arabic",
			input:    "مثال ١٢٣",
			expected: "مثال-١٢٣",
		},
		{
			name:     "devanagari numerals",
			input:    "परीक्षण १२३",
			expected: "परीक्षण-१२३",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
