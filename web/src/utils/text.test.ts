import assert from "node:assert/strict";
import { test } from "vitest";

import { capitalise, humanise, pluralise, truncateText } from "./text";

test("capitalises the first character", () => {
  assert.strictEqual(capitalise("active"), "Active");
});

test("pluralises simple nouns", () => {
  assert.strictEqual(pluralise(1, "action"), "action");
  assert.strictEqual(pluralise(2, "action"), "actions");
  assert.strictEqual(pluralise(2, "reply", "replies"), "replies");
});

test("humanises underscore-separated values", () => {
  assert.strictEqual(humanise("permission_denied"), "Permission denied");
});

test("trims and truncates text", () => {
  assert.strictEqual(truncateText("  short text  ", 20), "short text");
  assert.strictEqual(
    truncateText("  a longer piece of text  ", 8),
    "a longer…",
  );
  assert.strictEqual(truncateText(undefined), undefined);
});
