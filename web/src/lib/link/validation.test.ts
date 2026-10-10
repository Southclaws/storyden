import assert from "node:assert/strict";
import { test } from "vitest";

import { isValidLinkLike, normalizeLink } from "./validation";

test("isValidLinkLike accepts plain domains", () => {
  assert.ok(isValidLinkLike("example.com"));
  assert.ok(isValidLinkLike("sub.example.com"));
});

test("isValidLinkLike accepts http/https URLs", () => {
  assert.ok(isValidLinkLike("https://example.com"));
  assert.ok(isValidLinkLike("http://example.com/path"));
});

test("isValidLinkLike rejects unsafe protocols", () => {
  assert.ok(!isValidLinkLike("javascript:alert(1)"));
  assert.ok(!isValidLinkLike("data:text/plain,hello"));
  assert.ok(!isValidLinkLike("ftp://example.com"));
});

test("isValidLinkLike rejects whitespace and non-domain text", () => {
  assert.ok(!isValidLinkLike("hello world"));
  assert.ok(!isValidLinkLike("example"));
  assert.ok(!isValidLinkLike("   "));
});

test("normalizeLink normalizes domains to https", () => {
  assert.strictEqual(normalizeLink("example.com"), "https://example.com/");
  assert.strictEqual(
    normalizeLink("sub.example.com"),
    "https://sub.example.com/",
  );
});

test("normalizeLink preserves safe absolute URLs", () => {
  assert.strictEqual(
    normalizeLink("https://example.com/path?q=1"),
    "https://example.com/path?q=1",
  );
  assert.strictEqual(
    normalizeLink("http://example.com"),
    "http://example.com/",
  );
});

test("normalizeLink rejects unsafe and invalid values", () => {
  assert.strictEqual(normalizeLink(undefined), undefined);
  assert.strictEqual(normalizeLink(""), undefined);
  assert.strictEqual(normalizeLink("   "), undefined);
  assert.strictEqual(normalizeLink("javascript:alert(1)"), undefined);
  assert.strictEqual(normalizeLink("hello world"), undefined);
});
