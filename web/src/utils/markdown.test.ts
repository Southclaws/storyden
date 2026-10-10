import assert from "node:assert/strict";
import { test } from "vitest";

import { markdownURLTransform } from "./markdown";

test("markdownURLTransform preserves SDR URLs", () => {
  assert.strictEqual(
    markdownURLTransform("sdr:node/d8818ueot5pfij6bvm90"),
    "sdr:node/d8818ueot5pfij6bvm90",
  );
});

test("markdownURLTransform rejects unsafe protocols", () => {
  assert.strictEqual(markdownURLTransform("javascript:alert(1)"), "");
});
