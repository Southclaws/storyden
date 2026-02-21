import { test } from "uvu";
import * as assert from "uvu/assert";

import { slugify } from "./slugify";

test("cyrillic", () => {
  assert.is(slugify("Документация"), "документация");
});

test("japanese with spaces", () => {
  assert.is(slugify("日本語 テスト"), "日本語-テスト");
});

test("greek", () => {
  assert.is(slugify("Παράδειγμα"), "παράδειγμα");
});

test("hindi", () => {
  assert.is(slugify("परीक्षण दस्तावेज़"), "परीक्षण-दस्तावेज़");
});

test("korean with full-width space", () => {
  assert.is(slugify("문서　테스트"), "문서-테스트");
});

test("hebrew with hyphen", () => {
  assert.is(slugify("תיעוד-מערכת"), "תיעוד-מערכת");
});

test("persian with hyphen", () => {
  assert.is(slugify("مثالِ-سادِه"), "مثالِ-سادِه");
});

test("basic english", () => {
  assert.is(slugify("Hello World"), "hello-world");
});

test("uppercase to lowercase", () => {
  assert.is(slugify("HELLO WORLD"), "hello-world");
});

test("mixed case", () => {
  assert.is(slugify("HeLLo WoRLd"), "hello-world");
});

test("leading spaces", () => {
  assert.is(slugify("   hello world"), "hello-world");
});

test("trailing spaces", () => {
  assert.is(slugify("hello world   "), "hello-world");
});

test("leading and trailing spaces", () => {
  assert.is(slugify("   hello world   "), "hello-world");
});

test("multiple spaces", () => {
  assert.is(slugify("hello    world"), "hello-world");
});

test("multiple hyphens", () => {
  assert.is(slugify("hello----world"), "hello-world");
});

test("leading hyphens", () => {
  assert.is(slugify("---hello-world"), "hello-world");
});

test("trailing hyphens", () => {
  assert.is(slugify("hello-world---"), "hello-world");
});

test("leading underscores", () => {
  assert.is(slugify("___hello_world"), "hello_world");
});

test("trailing underscores", () => {
  assert.is(slugify("hello_world___"), "hello_world");
});

test("special characters", () => {
  assert.is(slugify("hello@world!test#123"), "hello-world-test-123");
});

test("punctuation", () => {
  assert.is(slugify("hello, world. test?"), "hello-world-test");
});

test("brackets and parens", () => {
  assert.is(slugify("hello (world) [test]"), "hello-world-test");
});

test("emojis", () => {
  assert.is(slugify("hello 👋 world 🌍"), "hello-world");
});

test("mixed emojis and text", () => {
  assert.is(slugify("🎉 Party Time 🎊"), "party-time");
});

test("numbers", () => {
  assert.is(slugify("test 123 456"), "test-123-456");
});

test("numbers with letters", () => {
  assert.is(slugify("test123abc456"), "test123abc456");
});

test("accented characters", () => {
  assert.is(slugify("café résumé"), "café-résumé");
});

test("german umlauts", () => {
  assert.is(slugify("Über Größe"), "über-größe");
});

test("full-width characters", () => {
  assert.is(slugify("ｈｅｌｌｏ　ｗｏｒｌｄ"), "hello-world");
});

test("mixed full-width and half-width", () => {
  assert.is(slugify("hello ｗｏｒｌｄ"), "hello-world");
});

test("tabs", () => {
  assert.is(slugify("hello\tworld"), "hello-world");
});

test("newlines", () => {
  assert.is(slugify("hello\nworld"), "hello-world");
});

test("carriage returns", () => {
  assert.is(slugify("hello\rworld"), "hello-world");
});

test("mixed whitespace", () => {
  assert.is(slugify("hello \t\n\r world"), "hello-world");
});

test("zero width space", () => {
  assert.is(slugify("hello\u200Bworld"), "hello-world");
});

test("zero width non-joiner", () => {
  assert.is(slugify("hello\u200Cworld"), "hello-world");
});

test("zero width joiner", () => {
  assert.is(slugify("hello\u200Dworld"), "hello-world");
});

test("soft hyphen", () => {
  assert.is(slugify("hello\u00ADworld"), "hello-world");
});

test("non-breaking space", () => {
  assert.is(slugify("hello\u00A0world"), "hello-world");
});

test("empty string", () => {
  assert.is(slugify(""), "");
});

test("only spaces", () => {
  assert.is(slugify("     "), "");
});

test("only hyphens", () => {
  assert.is(slugify("-----"), "");
});

test("only special characters", () => {
  assert.is(slugify("!@#$%^&*()"), "");
});

test("only emojis", () => {
  assert.is(slugify("👋🌍🎉"), "");
});

test("url with protocol", () => {
  assert.is(slugify("https://example.com"), "https-example-com");
});

test("email address", () => {
  assert.is(slugify("user@example.com"), "user-example-com");
});

test("path-like string", () => {
  assert.is(slugify("path/to/file.txt"), "path-to-file-txt");
});

test("mixed scripts", () => {
  assert.is(slugify("English 日本語 Русский"), "english-日本語-русский");
});

test("right-to-left scripts", () => {
  assert.is(slugify("العربية עברית"), "العربية-עברית");
});

test("chinese characters", () => {
  assert.is(slugify("中文测试"), "中文测试");
});

test("chinese with spaces", () => {
  assert.is(slugify("中文 测试"), "中文-测试");
});

test("thai", () => {
  assert.is(slugify("ทดสอบ ภาษาไทย"), "ทดสอบ-ภาษาไทย");
});

test("quotes and apostrophes", () => {
  assert.is(slugify('it\'s a "test"'), "it-s-a-test");
});

test("slashes and backslashes", () => {
  assert.is(slugify("test/slash\\backslash"), "test-slash-backslash");
});

test("currency symbols", () => {
  assert.is(slugify("$100 €50 ¥1000"), "100-50-1000");
});

test("math symbols", () => {
  assert.is(slugify("x + y = z"), "x-y-z");
});

test("already valid slug", () => {
  assert.is(slugify("hello-world"), "hello-world");
});

test("underscores preserved", () => {
  assert.is(slugify("hello_world_test"), "hello_world_test");
});

test("mixed hyphens and underscores", () => {
  assert.is(slugify("hello-world_test"), "hello-world_test");
});

test("control characters", () => {
  assert.is(slugify("hello\x00\x01\x02world"), "hello-world");
});

test("ligatures", () => {
  assert.is(slugify("ﬁle ﬂag"), "file-flag");
});

test("superscripts and subscripts", () => {
  assert.is(slugify("x² + y₃"), "x2-y3");
});

test("fractions", () => {
  assert.is(slugify("½ + ¼"), "1-2-1-4");
});

test("combining diacritics", () => {
  assert.is(slugify("e\u0301cole"), "école");
});

test("arabic numerals in arabic", () => {
  assert.is(slugify("مثال ١٢٣"), "مثال-١٢٣");
});

test("devanagari numerals", () => {
  assert.is(slugify("परीक्षण १२३"), "परीक्षण-१२३");
});

test.run();
