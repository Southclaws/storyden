import { test } from "uvu";
import * as assert from "uvu/assert";

import { capitalise, humanise, pluralise, truncateText } from "./text";

test("capitalises the first character", () => {
  assert.is(capitalise("active"), "Active");
});

test("pluralises simple nouns", () => {
  assert.is(pluralise(1, "action"), "action");
  assert.is(pluralise(2, "action"), "actions");
  assert.is(pluralise(2, "reply", "replies"), "replies");
});

test("humanises underscore-separated values", () => {
  assert.is(humanise("permission_denied"), "Permission denied");
});

test("trims and truncates text", () => {
  assert.is(truncateText("  short text  ", 20), "short text");
  assert.is(truncateText("  a longer piece of text  ", 8), "a longer…");
  assert.is(truncateText(undefined), undefined);
});

test.run();
