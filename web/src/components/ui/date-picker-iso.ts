import { type DateValue, parseDate } from "@internationalized/date";
import { formatInTimeZone } from "date-fns-tz";

export function formatISODate(value: DateValue) {
  return formatInTimeZone(value.toDate("UTC"), "UTC", "yyyy-MM-dd");
}

export function parseISODate(value: string) {
  try {
    return parseDate(value);
  } catch {
    return undefined;
  }
}
