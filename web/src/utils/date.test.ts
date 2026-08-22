import { test } from "uvu";
import * as assert from "uvu/assert";

import { formatSeconds } from "./date";

test("formats a number of seconds as a duration", () => {
  assert.is(formatSeconds(3661), "1 hour 1 minute 1 second");
});

test.run();
