import { RecurrenceSchedule } from "@/api/openapi-schema";
import { capitalise, pluralise } from "@/utils/text";

const occurrenceFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: "medium",
  timeStyle: "short",
});
const monthFormatter = new Intl.DateTimeFormat(undefined, { month: "long" });
const weekdayFormatter = new Intl.DateTimeFormat("en-US", {
  weekday: "long",
  timeZone: "UTC",
});

export function describeTrailSchedule(schedule: RecurrenceSchedule): string {
  const time = schedule.start.slice(11, 16);
  if (schedule.rule.count === 1) {
    return `Once on ${schedule.start.slice(0, 10)} at ${time}`;
  }

  const interval = schedule.rule.interval;
  switch (schedule.rule.frequency) {
    case "hourly":
      return every(interval, "hour", "hours");
    case "daily":
      return every(interval, "day", "days");
    case "weekly": {
      const days = schedule.rule.by_weekday ?? [
        weekdayForStart(schedule.start),
      ];
      return `${every(interval, "week", "weeks")} on ${days.map(capitalise).join(", ")} at ${time}`;
    }
    case "monthly":
      return `${every(interval, "month", "months")} on ${describeMonthDays(schedule.rule.by_month_day, schedule.start)} at ${time}`;
    case "yearly": {
      const months = schedule.rule.by_month ?? [
        Number(schedule.start.slice(5, 7)),
      ];
      return `${every(interval, "year", "years")} in ${months.map(monthName).join(", ")} on ${describeMonthDays(schedule.rule.by_month_day, schedule.start)} at ${time}`;
    }
  }
}

export function formatOccurrence(value?: string): string {
  if (!value) return "—";
  return occurrenceFormatter.format(new Date(value));
}

function every(interval: number, singular: string, plural: string): string {
  return interval === 1
    ? `Every ${singular}`
    : `Every ${interval} ${pluralise(interval, singular, plural)}`;
}

function describeMonthDays(
  values: number[] | undefined,
  start: string,
): string {
  const days = values ?? [Number(start.slice(8, 10))];
  return days
    .map((value) => {
      if (value === -1) return "the last day";
      if (value < 0) return `the ${ordinal(Math.abs(value))} day from the end`;
      return `day ${value}`;
    })
    .join(", ");
}

function ordinal(value: number): string {
  const remainder = value % 100;
  if (remainder >= 11 && remainder <= 13) return `${value}th`;
  switch (value % 10) {
    case 1:
      return `${value}st`;
    case 2:
      return `${value}nd`;
    case 3:
      return `${value}rd`;
    default:
      return `${value}th`;
  }
}

function monthName(value: number): string {
  return monthFormatter.format(new Date(2000, value - 1, 1));
}

function weekdayForStart(value: string): string {
  return weekdayFormatter
    .format(new Date(`${value.slice(0, 10)}T00:00:00Z`))
    .toLowerCase();
}
