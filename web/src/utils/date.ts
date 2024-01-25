import { FormatDistanceToNowOptions } from "date-fns";

const formatDistanceLocale = {
  lessThanXSeconds: "{{count}}s",
  xSeconds: "{{count}}s",
  halfAMinute: "30s",
  lessThanXMinutes: "{{count}}m",
  xMinutes: "{{count}}m",
  aboutXHours: "{{count}}h",
  xHours: "{{count}}h",
  xDays: "{{count}}d",
  aboutXWeeks: "{{count}}w",
  xWeeks: "{{count}}w",
  aboutXMonths: "{{count}}m",
  xMonths: "{{count}}m",
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
