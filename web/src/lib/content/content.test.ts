import assert from "node:assert/strict";
import test from "node:test";

import { isContentEmpty } from "./content";

test("isContentEmpty recognises empty rich-text editor output", () => {
  assert.equal(isContentEmpty("<p></p>"), true);
  assert.equal(isContentEmpty("<p>Community update</p>"), false);
});
