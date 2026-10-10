import assert from "node:assert/strict";
import { test } from "vitest";

import { slugify, slugifyDraft } from "./slugify";

test("cyrillic", () => {
  assert.strictEqual(slugify("Документация"), "документация");
});

test("japanese with spaces", () => {
  assert.strictEqual(slugify("日本語 テスト"), "日本語-テスト");
});

test("greek", () => {
  assert.strictEqual(slugify("Παράδειγμα"), "παράδειγμα");
});

test("hindi", () => {
  assert.strictEqual(slugify("परीक्षण दस्तावेज़"), "परीक्षण-दस्तावेज़");
});

test("korean with full-width space", () => {
  assert.strictEqual(slugify("문서　테스트"), "문서-테스트");
});

test("hebrew with hyphen", () => {
  assert.strictEqual(slugify("תיעוד-מערכת"), "תיעוד-מערכת");
});

test("persian with hyphen", () => {
  assert.strictEqual(slugify("مثالِ-سادِه"), "مثالِ-سادِه");
});

test("basic english", () => {
  assert.strictEqual(slugify("Hello World"), "hello-world");
});

test("uppercase to lowercase", () => {
  assert.strictEqual(slugify("HELLO WORLD"), "hello-world");
});

test("mixed case", () => {
  assert.strictEqual(slugify("HeLLo WoRLd"), "hello-world");
});

test("leading spaces", () => {
  assert.strictEqual(slugify("   hello world"), "hello-world");
});

test("trailing spaces", () => {
  assert.strictEqual(slugify("hello world   "), "hello-world");
});

test("leading and trailing spaces", () => {
  assert.strictEqual(slugify("   hello world   "), "hello-world");
});

test("multiple spaces", () => {
  assert.strictEqual(slugify("hello    world"), "hello-world");
});

test("multiple hyphens", () => {
  assert.strictEqual(slugify("hello----world"), "hello-world");
});

test("leading hyphens", () => {
  assert.strictEqual(slugify("---hello-world"), "hello-world");
});

test("trailing hyphens", () => {
  assert.strictEqual(slugify("hello-world---"), "hello-world");
});

test("leading underscores", () => {
  assert.strictEqual(slugify("___hello_world"), "hello_world");
});

test("trailing underscores", () => {
  assert.strictEqual(slugify("hello_world___"), "hello_world");
});

test("special characters", () => {
  assert.strictEqual(slugify("hello@world!test#123"), "hello-world-test-123");
});

test("punctuation", () => {
  assert.strictEqual(slugify("hello, world. test?"), "hello-world-test");
});

test("brackets and parens", () => {
  assert.strictEqual(slugify("hello (world) [test]"), "hello-world-test");
});

test("emojis", () => {
  assert.strictEqual(slugify("hello 👋 world 🌍"), "hello-world");
});

test("mixed emojis and text", () => {
  assert.strictEqual(slugify("🎉 Party Time 🎊"), "party-time");
});

test("numbers", () => {
  assert.strictEqual(slugify("test 123 456"), "test-123-456");
});

test("numbers with letters", () => {
  assert.strictEqual(slugify("test123abc456"), "test123abc456");
});

test("accented characters", () => {
  assert.strictEqual(slugify("café résumé"), "café-résumé");
});

test("german umlauts", () => {
  assert.strictEqual(slugify("Über Größe"), "über-größe");
});

test("full-width characters", () => {
  assert.strictEqual(slugify("ｈｅｌｌｏ　ｗｏｒｌｄ"), "hello-world");
});

test("mixed full-width and half-width", () => {
  assert.strictEqual(slugify("hello ｗｏｒｌｄ"), "hello-world");
});

test("tabs", () => {
  assert.strictEqual(slugify("hello\tworld"), "hello-world");
});

test("newlines", () => {
  assert.strictEqual(slugify("hello\nworld"), "hello-world");
});

test("carriage returns", () => {
  assert.strictEqual(slugify("hello\rworld"), "hello-world");
});

test("mixed whitespace", () => {
  assert.strictEqual(slugify("hello \t\n\r world"), "hello-world");
});

test("zero width space", () => {
  assert.strictEqual(slugify("hello\u200Bworld"), "hello-world");
});

test("zero width non-joiner", () => {
  assert.strictEqual(slugify("hello\u200Cworld"), "hello-world");
});

test("zero width joiner", () => {
  assert.strictEqual(slugify("hello\u200Dworld"), "hello-world");
});

test("soft hyphen", () => {
  assert.strictEqual(slugify("hello\u00ADworld"), "hello-world");
});

test("non-breaking space", () => {
  assert.strictEqual(slugify("hello\u00A0world"), "hello-world");
});

test("empty string", () => {
  assert.strictEqual(slugify(""), "");
});

test("only spaces", () => {
  assert.strictEqual(slugify("     "), "");
});

test("only hyphens", () => {
  assert.strictEqual(slugify("-----"), "");
});

test("only special characters", () => {
  assert.strictEqual(slugify("!@#$%^&*()"), "");
});

test("only emojis", () => {
  assert.strictEqual(slugify("👋🌍🎉"), "");
});

test("url with protocol", () => {
  assert.strictEqual(slugify("https://example.com"), "https-example-com");
});

test("email address", () => {
  assert.strictEqual(slugify("user@example.com"), "user-example-com");
});

test("path-like string", () => {
  assert.strictEqual(slugify("path/to/file.txt"), "path-to-file-txt");
});

test("mixed scripts", () => {
  assert.strictEqual(
    slugify("English 日本語 Русский"),
    "english-日本語-русский",
  );
});

test("right-to-left scripts", () => {
  assert.strictEqual(slugify("العربية עברית"), "العربية-עברית");
});

test("chinese characters", () => {
  assert.strictEqual(slugify("中文测试"), "中文测试");
});

test("chinese with spaces", () => {
  assert.strictEqual(slugify("中文 测试"), "中文-测试");
});

test("thai", () => {
  assert.strictEqual(slugify("ทดสอบ ภาษาไทย"), "ทดสอบ-ภาษาไทย");
});

test("quotes and apostrophes", () => {
  assert.strictEqual(slugify('it\'s a "test"'), "it-s-a-test");
});

test("slashes and backslashes", () => {
  assert.strictEqual(slugify("test/slash\\backslash"), "test-slash-backslash");
});

test("currency symbols", () => {
  assert.strictEqual(slugify("$100 €50 ¥1000"), "100-50-1000");
});

test("math symbols", () => {
  assert.strictEqual(slugify("x + y = z"), "x-y-z");
});

test("already valid slug", () => {
  assert.strictEqual(slugify("hello-world"), "hello-world");
});

test("underscores preserved", () => {
  assert.strictEqual(slugify("hello_world_test"), "hello_world_test");
});

test("mixed hyphens and underscores", () => {
  assert.strictEqual(slugify("hello-world_test"), "hello-world_test");
});

test("control characters", () => {
  assert.strictEqual(slugify("hello\x00\x01\x02world"), "hello-world");
});

test("ligatures", () => {
  assert.strictEqual(slugify("ﬁle ﬂag"), "file-flag");
});

test("superscripts and subscripts", () => {
  assert.strictEqual(slugify("x² + y₃"), "x2-y3");
});

test("fractions", () => {
  assert.strictEqual(slugify("½ + ¼"), "1-2-1-4");
});

test("combining diacritics", () => {
  assert.strictEqual(slugify("e\u0301cole"), "école");
});

test("arabic numerals in arabic", () => {
  assert.strictEqual(slugify("مثال ١٢٣"), "مثال-١٢٣");
});

test("devanagari numerals", () => {
  assert.strictEqual(slugify("परीक्षण १२३"), "परीक्षण-१२३");
});

test("draft slug is normalised while typing", () => {
  assert.strictEqual(
    slugifyDraft("General!!!!!dnwHDAUIDHAODH3r3u18 2"),
    "general-dnwhdauidhaodh3r3u18-2",
  );
});

test("draft slug preserves one pending word separator", () => {
  assert.strictEqual(slugifyDraft("日本語   "), "日本語-");
  assert.strictEqual(slugifyDraft("hello___"), "hello_");
});

test("draft slug removes leading separators", () => {
  assert.strictEqual(slugifyDraft("---hello"), "hello");
  assert.strictEqual(slugifyDraft("---"), "");
});
