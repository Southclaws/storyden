package case_convert_test

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/Southclaws/storyden/internal/case_convert"
)

func TestToCamel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"whitespace only", " \t\r\n\v\f\u00a0\u2003", ""},
		{"delimiters only", "_ - . __--..", ""},
		{"punctuation only", "!@#$%^&*()[]{}:;,/?\\|+=~`\"'", ""},
		{"lowercase letter", "a", "A"},
		{"uppercase letter", "Z", "Z"},
		{"lowercase alphabet", "abcdefghijklmnopqrstuvwxyz", "Abcdefghijklmnopqrstuvwxyz"},
		{"uppercase alphabet", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "Abcdefghijklmnopqrstuvwxyz"},
		{"word", "test", "Test"},
		{"existing camel case is one word", "testCase", "Testcase"},
		{"existing pascal case is one word", "TestCase", "Testcase"},
		{"snake case", "test_case", "TestCase"},
		{"kebab case", "test-case", "TestCase"},
		{"dot separated", "test.case", "TestCase"},
		{"space separated", "test case", "TestCase"},
		{"many words", "many_many_words", "ManyManyWords"},
		{"repeated delimiters", "__test___case__", "TestCase"},
		{"mixed delimiters", ".-test _.- case-.", "TestCase"},
		{"padded whitespace", "\t\n test  case \r\n", "TestCase"},
		{"unicode padding", "\u00a0\u2003test_case\u202f\u3000", "TestCase"},
		{"single letter words", "a_b_c", "ABC"},
		{"uppercase words", "CONSTANT_CASE", "ConstantCase"},
		{"permission name", "READ_PUBLISHED_THREADS", "ReadPublishedThreads"},
		{"oauth permission", "USE_OAUTH_CLIENTS", "UseOauthClients"},
		{"two letter acronym", "ID", "Id"},
		{"acronym", "HTTP", "Http"},
		{"acronym within word", "HTTPServer", "Httpserver"},
		{"acronyms separated", "HTTP_API_ID", "HttpApiId"},
		{"digit", "0", "0"},
		{"digits only", "0123456789", "0123456789"},
		{"leading digits", "123HELLO", "123hello"},
		{"trailing digits", "hello123", "Hello123"},
		{"digits within word", "version2TOKEN", "Version2token"},
		{"digit followed by delimiter", "version2_token", "Version2Token"},
		{"separated digits", "v1_2.3-beta4", "V123Beta4"},
		{"punctuation between words", "foo!bar", "FooBar"},
		{"slash between words", "foo/bar", "FooBar"},
		{"backslash between words", "foo\\bar", "FooBar"},
		{"leading punctuation", "!hello", "Hello"},
		{"trailing punctuation", "hello!", "Hello"},
		{"internal controls", "one\ttwo\nthree\rfour\vfive\fsix", "OneTwoThreeFourFiveSix"},
		{"unicode spaces", "foo\u00a0bar\u2003baz", "FooBarBaz"},
		{"internal null", "foo\x00bar", "FooBar"},
		{"accented words", "ÉCOLE_déjà_vu", "ÉcoleDéjàVu"},
		{"greek", "ΚΑΛΗΜΈΡΑ_ΚΌΣΜΕ", "ΚαλημέραΚόσμε"},
		{"cyrillic", "ПРИВЕТ_МИР", "ПриветМир"},
		{"uncased letters", "東京_世界", "東京世界"},
		{"unicode digits", "版本٢_３", "版本٢３"},
		{"unicode compatibility letter", "KELVIN", "Kelvin"},
		{"lowercase mapping grows byte length", "HTTPȺ", "Httpⱥ"},
		{"lowercase mapping shrinks byte length", "FOOK", "Fook"},
		{"uppercase mapping grows byte length", "ȿ_name", "ⱾName"},
		{"simple case mapping does not expand sharp s", "ß_TEST", "ßTest"},
		{"combining marks separate words", "e\u0301clair", "EClair"},
		{"emoji only", "👩🏽‍💻", ""},
		{"emoji separates words", "hello👩🏽‍💻world", "HelloWorld"},
		{"byte order mark", "\ufeffhello", "Hello"},
		{"replacement character separates words", "a\ufffdb", "AB"},
		{"invalid UTF8 separates words", "\xffhello\xc3\x28_world\x80", "HelloWorld"},
		{"all invalid UTF8", "\xff\xfe\x80", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := case_convert.ToCamel(tc.input); got != tc.want {
				t.Errorf("ToCamel(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestToCamelLongInput(t *testing.T) {
	t.Parallel()

	input := strings.Repeat("HTTP2_api-東京.", 10_000)
	want := strings.Repeat("Http2Api東京", 10_000)

	if got := case_convert.ToCamel(input); got != want {
		t.Errorf("long input conversion differs: got %d bytes, want %d", len(got), len(want))
	}
}

func FuzzToCamel(f *testing.F) {
	for _, seed := range []string{
		"", "READ_PUBLISHED_THREADS", "USE_OAUTH_CLIENTS", "HTTPServer",
		"A_B", "a1b2c", "_!hello", "foo\tbar", "déjà_vu", "Kelvin",
		"👩🏽‍💻", "\xffhello\xc3\x28_world\x80", "\x00", "\u00a0hello\u3000",
		"ȿ_name", "ÉCOLE_東京_١٢٣", "e\u0301clair", "ı_İ", "HTTPȺ",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		got := case_convert.ToCamel(input)
		if !utf8.ValidString(got) {
			t.Fatalf("output is not valid UTF8: %q", got)
		}

		for _, character := range got {
			if !unicode.IsLetter(character) && !unicode.IsDigit(character) {
				t.Fatalf("output %q contains a non-letter/non-digit %q", got, character)
			}
		}

		lettersAndDigits := strings.Map(func(character rune) rune {
			if unicode.IsLetter(character) || unicode.IsDigit(character) {
				return character
			}

			return -1
		}, input)

		if utf8.RuneCountInString(got) != utf8.RuneCountInString(lettersAndDigits) {
			t.Fatalf("conversion changed the number of letters or digits: %q -> %q", input, got)
		}

		digits := func(value string) string {
			return strings.Map(func(character rune) rune {
				if unicode.IsDigit(character) {
					return character
				}

				return -1
			}, value)
		}

		if digits(got) != digits(input) {
			t.Fatalf("conversion changed the digits or their order: %q -> %q", input, got)
		}

		if utf8.RuneCountInString(input) == len(input) && !strings.EqualFold(got, lettersAndDigits) {
			t.Fatalf("conversion lost or reordered ASCII letters: %q -> %q", input, got)
		}
	})
}
