import assert from "node:assert/strict";
import { test } from "vitest";

import { isSlugReady, processMarkInput } from "./mark";

test("processMarkInput normalizes spaces, disallowed chars, and casing", () => {
  assert.strictEqual(processMarkInput("  My /Page?# Name  "), "-my-page-name-");
});

test("processMarkInput collapses repeated hyphens", () => {
  assert.strictEqual(processMarkInput("hello---world"), "hello-world");
});

test("processMarkInput collapses double hyphens created by removed characters", () => {
  assert.strictEqual(processMarkInput("a / b"), "a-b");
  assert.strictEqual(processMarkInput("a ? b"), "a-b");
});

test("processMarkInput preserves trailing hyphen while typing", () => {
  assert.strictEqual(processMarkInput("My Page "), "my-page-");
});

test("isSlugReady accepts valid slugs", () => {
  assert.ok(isSlugReady("my-page_1"));
});

test("isSlugReady rejects trailing hyphen", () => {
  assert.ok(!isSlugReady("my-page-"));
});

test("isSlugReady rejects whitespace and reserved URL chars", () => {
  assert.ok(!isSlugReady("my page"));
  assert.ok(!isSlugReady("my/page"));
  assert.ok(!isSlugReady("my?page"));
  assert.ok(!isSlugReady("my#page"));
  assert.ok(!isSlugReady("my%page"));
});
