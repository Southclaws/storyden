import { z } from "zod";

import {
  RecurrenceSchedule,
  RecurrenceWeekday,
  TrailMutableProps,
} from "@/api/openapi-schema";

export const scheduleKinds = [
  "once",
  "interval",
  "weekly",
  "monthly",
  "yearly",
] as const;

export const weekdays = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
] as const;

export const TrailFormSchema = z
  .object({
    name: z.string().trim().min(1, "Name is required").max(120),
    description: z.string(),
    scheduleKind: z.enum(scheduleKinds),
    timezone: z.string().trim().min(1, "Timezone is required"),
    startDate: z.string().min(1, "Starting date is required"),
    localTime: z.string().regex(/^\d{2}:\d{2}$/, "Local time is required"),
    interval: z.number().int().min(1, "Interval must be at least 1"),
    intervalUnit: z.enum(["hours", "days"]),
    selectedDays: z.array(z.enum(weekdays)),
    month: z.string(),
    monthDay: z.string(),
    actions: z
      .array(
        z.object({
          type: z.literal("robot_run"),
          robot_ref: z.string().min(1, "Choose a Robot"),
          instruction: z.string().trim().min(1, "Instruction is required"),
        }),
      )
      .min(1, "Add at least one Robot action"),
  })
  .superRefine((values, context) => {
    if (values.scheduleKind === "weekly" && values.selectedDays.length === 0) {
      context.addIssue({
        code: "custom",
        message: "Choose at least one weekday",
        path: ["selectedDays"],
      });
    }
  });

export type TrailFormValues = z.infer<typeof TrailFormSchema>;

export function initialTrailFormValues(
  initialValue: TrailMutableProps | undefined,
  timezone: string,
  startDate: string,
): TrailFormValues {
  const schedule =
    initialValue?.trigger.type === "schedule"
      ? initialValue.trigger.schedule
      : undefined;
  const initialActions = initialValue?.actions ?? [];

  return {
    name: initialValue?.name ?? "",
    description: initialValue?.description ?? "",
    scheduleKind: scheduleKind(schedule),
    timezone: schedule?.timezone ?? timezone,
    startDate: schedule?.start.slice(0, 10) ?? startDate,
    localTime: schedule?.start.slice(11, 16) ?? "09:00",
    interval: schedule?.rule.interval ?? 1,
    intervalUnit: schedule?.rule.frequency === "hourly" ? "hours" : "days",
    selectedDays: schedule?.rule.by_weekday ?? [
      weekdayForDate(schedule?.start ?? startDate),
    ],
    month: String(
      schedule?.rule.by_month?.[0] ?? Number(startDate.slice(5, 7)),
    ),
    monthDay: schedule?.rule.by_month_day?.includes(-1)
      ? "last"
      : String(positiveMonthDay(schedule) ?? Number(startDate.slice(8, 10))),
    actions:
      initialActions.length > 0
        ? initialActions
        : [{ type: "robot_run", robot_ref: "", instruction: "" }],
  };
}

export function trailFormSchedule(values: TrailFormValues): RecurrenceSchedule {
  const start = `${values.startDate.slice(0, 10)}T${values.localTime}:00`;
  const monthDay = values.monthDay === "last" ? -1 : Number(values.monthDay);

  switch (values.scheduleKind) {
    case "once":
      return {
        start,
        timezone: values.timezone,
        rule: { frequency: "daily", interval: 1, count: 1 },
      };

    case "interval":
      return {
        start,
        timezone: values.timezone,
        rule: {
          frequency: values.intervalUnit === "hours" ? "hourly" : "daily",
          interval: values.interval,
        },
      };

    case "weekly":
      return {
        start,
        timezone: values.timezone,
        rule: {
          frequency: "weekly",
          interval: values.interval,
          by_weekday: values.selectedDays as RecurrenceWeekday[],
        },
      };

    case "monthly":
      return {
        start,
        timezone: values.timezone,
        rule: {
          frequency: "monthly",
          interval: values.interval,
          by_month_day: [monthDay],
        },
      };

    case "yearly":
      return {
        start,
        timezone: values.timezone,
        rule: {
          frequency: "yearly",
          interval: values.interval,
          by_month: [Number(values.month)],
          by_month_day: [monthDay],
        },
      };
  }
}

export function trailFormPayload(
  values: TrailFormValues,
  status: TrailMutableProps["status"],
): TrailMutableProps {
  return {
    name: values.name,
    description: values.description,
    status,
    trigger: {
      type: "schedule",
      schedule: trailFormSchedule(values),
    },
    actions: values.actions,
  };
}

function scheduleKind(
  schedule?: RecurrenceSchedule,
): TrailFormValues["scheduleKind"] {
  if (!schedule) return "weekly";
  if (schedule.rule.count === 1) return "once";

  switch (schedule.rule.frequency) {
    case "hourly":
    case "daily":
      return "interval";
    case "weekly":
    case "monthly":
    case "yearly":
      return schedule.rule.frequency;
  }
}

function positiveMonthDay(schedule?: RecurrenceSchedule): number | undefined {
  return schedule?.rule.by_month_day?.find((value) => value > 0);
}

function weekdayForDate(value: string): (typeof weekdays)[number] {
  const date = value.slice(0, 10);
  const day = new Date(`${date}T00:00:00Z`).getUTCDay();

  return weekdays[(day + 6) % 7]!;
}
