import { suite } from "uvu";
import * as assert from "uvu/assert";

import { parseThemeEditingFlag } from "./theme-editing";

const test = suite("theme editing flag");

test("only the explicit stored true value enables editing", () => {
  assert.is(parseThemeEditingFlag("true"), true);
  assert.is(parseThemeEditingFlag("false"), false);
  assert.is(parseThemeEditingFlag("1"), false);
  assert.is(parseThemeEditingFlag(null), false);
});

test.run();
