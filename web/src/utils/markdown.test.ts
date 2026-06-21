import { test } from "uvu";
import * as assert from "uvu/assert";

import { markdownURLTransform } from "./markdown";

test("markdownURLTransform preserves SDR URLs", () => {
  assert.is(
    markdownURLTransform("sdr:node/d8818ueot5pfij6bvm90"),
    "sdr:node/d8818ueot5pfij6bvm90",
  );
});

test("markdownURLTransform rejects unsafe protocols", () => {
  assert.is(markdownURLTransform("javascript:alert(1)"), "");
});

test.run();
