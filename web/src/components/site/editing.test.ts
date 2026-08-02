import { test } from "uvu";
import * as assert from "uvu/assert";

import { EditingSchema } from "./editing";

test("site edit mode parses the canonical query value", () => {
  assert.is(EditingSchema.parse("site"), "site");
});

test("legacy feed edit links enter site edit mode", () => {
  assert.is(EditingSchema.parse("feed"), "site");
});

test("empty edit values remain disabled", () => {
  assert.is(EditingSchema.parse(""), undefined);
});

test.run();
