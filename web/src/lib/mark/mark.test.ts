import { test } from "uvu";
import * as assert from "uvu/assert";

import { isSlugReady, processMarkInput } from "./mark";

test("processMarkInput normalizes spaces, disallowed chars, and casing", () => {
  assert.is(
    processMarkInput("  My /Page?# Name  "),
    "-my-page-name-",
  );
});

test("processMarkInput collapses repeated hyphens", () => {
  assert.is(processMarkInput("hello---world"), "hello-world");
});

test("processMarkInput collapses double hyphens created by removed characters", () => {
  assert.is(processMarkInput("a / b"), "a-b");
  assert.is(processMarkInput("a ? b"), "a-b");
});

test("processMarkInput preserves trailing hyphen while typing", () => {
  assert.is(processMarkInput("My Page "), "my-page-");
});

test("isSlugReady accepts valid slugs", () => {
  assert.ok(isSlugReady("my-page_1"));
});

test("isSlugReady rejects trailing hyphen", () => {
  assert.not.ok(isSlugReady("my-page-"));
});

test("isSlugReady rejects whitespace and reserved URL chars", () => {
  assert.not.ok(isSlugReady("my page"));
  assert.not.ok(isSlugReady("my/page"));
  assert.not.ok(isSlugReady("my?page"));
  assert.not.ok(isSlugReady("my#page"));
  assert.not.ok(isSlugReady("my%page"));
});

test.run();
