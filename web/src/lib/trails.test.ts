import assert from "node:assert/strict";
import { test } from "vitest";

import { describeTrailSchedule } from "./trails";

test("describes monthly last-day schedules", () => {
  assert.strictEqual(
    describeTrailSchedule({
      start: "2026-08-31T09:30:00",
      timezone: "Europe/London",
      rule: {
        frequency: "monthly",
        interval: 1,
        by_month_day: [-1],
      },
    }),
    "Every month on the last day at 09:30",
  );
});

test("describes selected weekdays", () => {
  assert.strictEqual(
    describeTrailSchedule({
      start: "2026-08-17T18:00:00",
      timezone: "UTC",
      rule: {
        frequency: "weekly",
        interval: 2,
        by_weekday: ["monday", "friday"],
      },
    }),
    "Every 2 weeks on Monday, Friday at 18:00",
  );
});

test("describes normalized one-time schedules", () => {
  assert.strictEqual(
    describeTrailSchedule({
      start: "2026-12-24T18:00:00",
      timezone: "Europe/London",
      rule: { frequency: "daily", interval: 1, count: 1 },
    }),
    "Once on 2026-12-24 at 18:00",
  );
});
