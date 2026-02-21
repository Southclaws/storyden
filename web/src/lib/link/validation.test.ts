import { test } from "uvu";
import * as assert from "uvu/assert";

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
  assert.not.ok(isValidLinkLike("javascript:alert(1)"));
  assert.not.ok(isValidLinkLike("data:text/plain,hello"));
  assert.not.ok(isValidLinkLike("ftp://example.com"));
});

test("isValidLinkLike rejects whitespace and non-domain text", () => {
  assert.not.ok(isValidLinkLike("hello world"));
  assert.not.ok(isValidLinkLike("example"));
  assert.not.ok(isValidLinkLike("   "));
});

test("normalizeLink normalizes domains to https", () => {
  assert.is(normalizeLink("example.com"), "https://example.com/");
  assert.is(normalizeLink("sub.example.com"), "https://sub.example.com/");
});

test("normalizeLink preserves safe absolute URLs", () => {
  assert.is(
    normalizeLink("https://example.com/path?q=1"),
    "https://example.com/path?q=1",
  );
  assert.is(normalizeLink("http://example.com"), "http://example.com/");
});

test("normalizeLink rejects unsafe and invalid values", () => {
  assert.is(normalizeLink(undefined), undefined);
  assert.is(normalizeLink(""), undefined);
  assert.is(normalizeLink("   "), undefined);
  assert.is(normalizeLink("javascript:alert(1)"), undefined);
  assert.is(normalizeLink("hello world"), undefined);
});

test.run();
