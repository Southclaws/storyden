import assert from "node:assert/strict";
import { test } from "vitest";

import { formatSeconds } from "./date";

test("formats a number of seconds as a duration", () => {
  assert.strictEqual(formatSeconds(3661), "1 hour 1 minute 1 second");
});
