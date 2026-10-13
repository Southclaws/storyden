import assert from "node:assert/strict";
import { test } from "vitest";

import {
  TrailFormSchema,
  initialTrailFormValues,
  trailFormPayload,
  trailFormSchedule,
} from "./trailForm";

test("trail forms normalize text and report nested action errors", () => {
  const input = initialTrailFormValues(undefined, "UTC", "2026-08-24");
  input.name = "  Daily review  ";
  input.actions = [
    {
      type: "robot_run",
      robot_ref: "robot_1",
      instruction: "  Review reports  ",
    },
  ];
  const parsed = TrailFormSchema.parse(input);
  assert.strictEqual(parsed.name, "Daily review");
  assert.strictEqual(parsed.actions[0]?.instruction, "Review reports");

  input.actions[0]!.instruction = "  ";
  const invalid = TrailFormSchema.safeParse(input);
  assert.ok(!invalid.success);
  assert.deepStrictEqual(
    invalid.error.issues.map(({ path, message }) => ({ path, message })),
    [
      {
        path: ["actions", 0, "instruction"],
        message: "Instruction is required",
      },
    ],
  );
});

test("trail intervals must be positive integers", () => {
  const input = initialTrailFormValues(undefined, "UTC", "2026-08-24");
  input.name = "Daily review";
  input.actions = [
    { type: "robot_run", robot_ref: "robot_1", instruction: "Review reports" },
  ];
  for (const interval of [0, -1, 1.5]) {
    const parsed = TrailFormSchema.safeParse({ ...input, interval });
    assert.ok(!parsed.success);
    assert.deepStrictEqual(
      parsed.error.issues.map((issue) => issue.path),
      [["interval"]],
    );
  }
});

test("creates a weekly schedule from explicit form fields", () => {
  const values = initialTrailFormValues(
    undefined,
    "Europe/London",
    "2026-08-24",
  );
  values.interval = 2;
  values.selectedDays = ["tuesday", "thursday"];
  values.localTime = "10:30";

  assert.deepStrictEqual(trailFormSchedule(values), {
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

  assert.deepStrictEqual(trailFormSchedule(values), {
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

  assert.deepStrictEqual(Object.keys(payload).sort(), [
    "actions",
    "description",
    "name",
    "status",
    "trigger",
  ]);
  assert.ok(!("account_id" in payload));
  assert.ok(!("created_by" in payload));
});

test("preserves and submits an event trigger", () => {
  const values = initialTrailFormValues(
    {
      name: "Respond to new reports",
      description: "",
      status: "active",
      trigger: {
        type: "event",
        events: ["EventReportCreated", "EventReportUpdated"],
      },
      actions: [
        {
          type: "robot_run",
          robot_ref: "robot_1",
          instruction: "Review the new report.",
        },
      ],
    },
    "UTC",
    "2026-08-24",
  );

  assert.strictEqual(values.triggerType, "event");
  assert.deepStrictEqual(values.events, [
    "EventReportCreated",
    "EventReportUpdated",
  ]);
  assert.deepStrictEqual(trailFormPayload(values, "active").trigger, {
    type: "event",
    events: ["EventReportCreated", "EventReportUpdated"],
  });
});
