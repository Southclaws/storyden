import { test } from "uvu";
import * as assert from "uvu/assert";

import {
  initialTrailFormValues,
  trailFormPayload,
  trailFormSchedule,
} from "./trailForm";

test("creates a weekly schedule from explicit form fields", () => {
  const values = initialTrailFormValues(
    undefined,
    "Europe/London",
    "2026-08-24",
  );
  values.interval = 2;
  values.selectedDays = ["tuesday", "thursday"];
  values.localTime = "10:30";

  assert.equal(trailFormSchedule(values), {
    start: "2026-08-24T10:30:00",
    timezone: "Europe/London",
    rule: {
      frequency: "weekly",
      interval: 2,
      by_weekday: ["tuesday", "thursday"],
    },
  });
});

test("normalizes a monthly last-day schedule", () => {
  const values = initialTrailFormValues(undefined, "UTC", "2026-08-24");
  values.scheduleKind = "monthly";
  values.monthDay = "last";

  assert.equal(trailFormSchedule(values), {
    start: "2026-08-24T09:00:00",
    timezone: "UTC",
    rule: {
      frequency: "monthly",
      interval: 1,
      by_month_day: [-1],
    },
  });
});

test("builds an initial payload without a creator field", () => {
  const values = initialTrailFormValues(undefined, "UTC", "2026-08-24");
  values.name = "Scheduled moderation check";
  values.actions = [
    {
      type: "robot_run",
      robot_ref: "robot_1",
      instruction: "Review new reports and finish with a structured result.",
    },
  ];

  const payload = trailFormPayload(values, "active");

  assert.equal(Object.keys(payload).sort(), [
    "actions",
    "description",
    "name",
    "status",
    "trigger",
  ]);
  assert.not.ok("account_id" in payload);
  assert.not.ok("created_by" in payload);
});

test.run();
