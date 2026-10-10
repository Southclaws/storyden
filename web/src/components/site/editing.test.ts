import assert from "node:assert/strict";
import { test } from "vitest";

import { EditingSchema } from "./editing";

test("site edit mode parses the canonical query value", () => {
  assert.strictEqual(EditingSchema.parse("site"), "site");
});

test("legacy feed edit links enter site edit mode", () => {
  assert.strictEqual(EditingSchema.parse("feed"), "site");
});

test("empty edit values remain disabled", () => {
  assert.strictEqual(EditingSchema.parse(""), undefined);
});
