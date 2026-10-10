import assert from "node:assert/strict";
import { describe, test } from "vitest";

import { parseThemeEditingFlag } from "./theme-editing";

describe("theme editing flag", () => {
  test("only the explicit stored true value enables editing", () => {
    assert.strictEqual(parseThemeEditingFlag("true"), true);
    assert.strictEqual(parseThemeEditingFlag("false"), false);
    assert.strictEqual(parseThemeEditingFlag("1"), false);
    assert.strictEqual(parseThemeEditingFlag(null), false);
  });
});
