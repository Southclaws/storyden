import { FormatDistanceToNowOptions, formatDistanceToNow } from "date-fns";

const formatDistanceLocale = {
  lessThanXSeconds: "{{count}}s",
  xSeconds: "{{count}}s",
  halfAMinute: "30s",
  lessThanXMinutes: "{{count}}min",
  xMinutes: "{{count}}min",
  aboutXHours: "{{count}}h",
  xHours: "{{count}}h",
  xDays: "{{count}}d",
  aboutXWeeks: "{{count}}w",
  xWeeks: "{{count}}w",
  aboutXMonths: "{{count}}mo",
  xMonths: "{{count}}mo",
  aboutXYears: "{{count}}y",
  xYears: "{{count}}y",
  overXYears: "{{count}}y",
  almostXYears: "{{count}}y",
};

export const formatDistance = (
  token: keyof typeof formatDistanceLocale,
  count: number,
) => {
  const result = formatDistanceLocale[token].replace(
    "{{count}}",
    count.toString(),
  );

  return result;
};

export const formatDistanceDefaults: FormatDistanceToNowOptions = {
  locale: { formatDistance },
};

export function timestamp(date: string | number | Date, short = true) {
  try {
    return formatDistanceToNow(
      date,
      short ? formatDistanceDefaults : { addSuffix: true },
    );
  } catch (e: unknown) {
    throw new Error(`Failed to format date: ${date}: error: ${e}`);
  }
}
